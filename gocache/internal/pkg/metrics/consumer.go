package metrics

/*

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Netrounds/ncc3-services/ncc3-metrics-service/cmd/bulkinsert"
	"github.com/Netrounds/ncc3-services/ncc3-metrics-service/cmd/metadata"
	metricproto "github.com/Netrounds/ncc3-services/ncc3-metrics-service/gen/protobuf/callexecuter"
	pluginproto "github.com/Netrounds/ncc3-services/ncc3-metrics-service/gen/protobuf/plugin"
	"github.com/Netrounds/ncc3-services/ncc3-metrics-service/internal/pkg/plugins"
	"github.com/golang/protobuf/proto"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
)

// Consumer struct holds the variables
type Consumer struct {
	maxBatchAge   int
	maxBatchSize  int
	pluginHandler *plugins.PluginHandler
	kafkaURL      string
	db            *sql.DB
	maxTimestamp  map[int32]int64
	mapBi         map[string]*bulkinsert.BulkInsert
	pluginChannel chan *pluginproto.Plugin
	mapPlugins    map[string]*pluginproto.Plugin
	aggregator    *plugins.AggregatorMap //localcopy
}
type MetricMsg struct {
	PluginName string
	Metrics    *metricproto.Metrics
}

// NewConsumer creates and initializes metrics consumer
func NewConsumer(kafkaURL string, maxBatchSize int, maxBatchAge int, pluginHandler *plugins.PluginHandler, db *sql.DB, pluginChannel chan *pluginproto.Plugin) *Consumer {
	consumer := Consumer{
		maxBatchAge:   maxBatchAge,
		maxBatchSize:  maxBatchSize,
		pluginHandler: pluginHandler,
		kafkaURL:      kafkaURL,
		db:            db,
		maxTimestamp:  make(map[int32]int64),
		mapBi:         make(map[string]*bulkinsert.BulkInsert),
		pluginChannel: pluginChannel,
		mapPlugins:    make(map[string]*pluginproto.Plugin),
	}
	return &consumer
}

// queryLatestTimestamps queries latest timestamps for all streams in the table
func (consumer *Consumer) queryLatestTimestamps(tableName string) {
	// check if table exists in the database
	query := fmt.Sprintf("SELECT max(time), measurement_id FROM %s GROUP BY measurement_id", tableName)
	log.Debug().Msgf("%s", query)
	rows, err := consumer.db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return
	}

	defer rows.Close()
	var maxTime time.Time
	var measurementID int32
	for rows.Next() {
		rows.Scan(&maxTime, &measurementID)
		log.Debug().Msgf("Read max time for measurement %d to %s", measurementID, maxTime.Format("2006-01-02 15:04:05"))
		consumer.maxTimestamp[measurementID] = maxTime.Unix()
	}
}
func (consumer *Consumer) waitForPluginsInitialized() {
	log.Info().Msg("Waiting for plugins to be intialized")
	for {
		var plugin *pluginproto.Plugin
		select {
		case plugin = <-consumer.pluginChannel:
			if plugin == nil {
				return
			}

		case <-time.After(60 * time.Second):
			return

		}
		// add plugin to map
		consumer.mapPlugins[plugin.Name] = plugin
		log.Debug().Msgf("Got plugin %s", plugin.Name)
	}
	consumer.aggregator = consumer.pluginHandler.AggregatorMap.CopyAggregatorMap()
}

// ConsumeMetrics consumes metrics from the topics
func (consumer *Consumer) ConsumeMetrics(topics []string, producer *kafka.Producer, metadataCache *metadata.Cache) error {

	consumer.waitForPluginsInitialized()

	// Create new consumer
	log.Info().Strs("topics", topics).Msg("creating consumer")

	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               consumer.kafkaURL,
		"go.application.rebalance.enable": true,
		"group.id":                        "metrics-service",
		"auto.offset.reset":               "earliest",
	})
	if err != nil {
		log.Panic().Err(err).Msg("failed to create Kafka consumer")
	}
	defer kafkaConsumer.Close()
	kafkaConsumer.SubscribeTopics(topics, nil)

	// reset message counters
	var msgCount = 0
	lastReport := time.Now()
	var countRejected = 0

	run := true
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	rollupState := NewRollUpStateManager(consumer.db,consumer.aggregator)

	rollupState.InitRollUpStates()


	for run == true {
		select {
		case sig := <-sigchan:
			log.Info().Str("signal", sig.String()).Msg("caught terminating signal")
			run = false
		case plugin := <-consumer.pluginChannel:
			// add plugin to map
			if plugin != nil {
				consumer.mapPlugins[plugin.Name] = plugin
				log.Debug().Msgf("Got new plugin %s", plugin.Name)
				//update aggregator local map
				consumer.aggregator = consumer.pluginHandler.AggregatorMap.CopyAggregatorMap()
			}
		default:
			consumerEvent := kafkaConsumer.Poll(100)
			if consumerEvent == nil {
				continue
			}
			switch event := consumerEvent.(type) {
			case kafka.AssignedPartitions:
				kafkaConsumer.Assign(event.Partitions)
				log.Debug().Msgf("assigned partitions %v", event)

			case kafka.RevokedPartitions:
				log.Debug().Msgf("revoked partitions %v", event)
				kafkaConsumer.Unassign()

			case *kafka.Message:
				// deserialize protobuf message
				metrics := &metricproto.Metrics{}
				err = proto.Unmarshal(event.Value, metrics)
				if err != nil {
					log.Error().Err(err).Msgf("Failed to unmarshal data from metrics topic %v", topics)
					continue
				}
				// todo: write metrics counter to prometheus as part of https://netrounds.atlassian.net/browse/NCC2019-1827
				// prom_metrics_received.Inc()

				pluginName := metadataCache.GetPluginName(metrics.StreamId)
				if len(pluginName) == 0 {
					log.Debug().Msgf("Cannot find plugin name for stream %v", metrics.StreamId)
					continue
				}
				plugin := consumer.mapPlugins[pluginName]
				if plugin == nil {
					log.Debug().Msgf("Plugin not found %s", pluginName)
					continue
				}

				// check for out of order timestamps
				if consumer.maxTimestamp[metrics.StreamId] != 0 && metrics.Timestamp.Seconds < consumer.maxTimestamp[metrics.StreamId] {
					log.Debug().Msgf("Rejected message from measurement id %d, timestamp %s (max %s)", metrics.StreamId,
						time.Unix(metrics.Timestamp.Seconds, 0).Format("2006-01-02 15:04:05"),
						time.Unix(consumer.maxTimestamp[metrics.StreamId], 0).Format("2006-01-02 15:04:05"))
					msgCount++
					countRejected++
					continue
				}
				consumer.maxTimestamp[metrics.StreamId] = metrics.Timestamp.Seconds

				// process metrics in same goroutine because user internal data
				consumer.processMetrics(&MetricMsg{PluginName: pluginName, Metrics: metrics})

				//send metrics for rollups
				rollupState.OnMessage(&MetricMsg{PluginName: pluginName, Metrics: metrics})

				//print stats
				elapsed := time.Since(lastReport)
				if elapsed.Seconds() > 1.0 {
					lag := time.Since(time.Unix(int64(metrics.Timestamp.Seconds), 0))
					log.Debug().Msgf("%v: consumed %d messages, %.0f per second, lag %.0f s, rejected %d messages", topics, msgCount,
						float64(msgCount)/elapsed.Seconds(), lag.Seconds(), countRejected)
					msgCount = 0
					lastReport = time.Now()
					countRejected = 0

				}
				msgCount++

			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				log.Error().Err(event)
			default:
				log.Debug().Msgf("Ignored %v", event)
			}
		}
	}

	log.Info().Msg("closing consumer")
	return nil
}

// processMetrics processes the metrics
func (consumer *Consumer) processMetrics(msg *MetricMsg) {
	plugin := consumer.mapPlugins[msg.PluginName]
	if plugin == nil {
		log.Debug().Msgf("Plugin not found %s", msg.PluginName)
		return
	}
	columnNames := []string{"time", "measurement_id"}
	table := "monitor_metrics_" + msg.PluginName
	for _, metric := range plugin.Metrics {
		columnNames = append(columnNames, metric.Name)
	}
	// Create bulk insert based to plugin spec
	if consumer.mapBi[table] == nil {
		consumer.mapBi[table] = bulkinsert.NewBulkInsert(consumer.db, table, consumer.maxBatchSize, consumer.maxBatchAge, columnNames)
		// this is first time plugin is seen, initialize timestamps for it
		consumer.queryLatestTimestamps(table)
	}

	columnValues := []string{fmt.Sprintf("'%s'", time.Unix(msg.Metrics.Timestamp.Seconds, 0).UTC().Format("2006-01-02 15:04:05")),
		fmt.Sprintf("%d", msg.Metrics.StreamId)}

	for _, metric := range plugin.Metrics {
		value := msg.Metrics.Values[metric.Name]

		// add \nil for all missing values
		if value == nil {
			columnValues = append(columnValues, "\nil")
			continue
		}
		switch value.Type.(type) {
		case *metricproto.MetricValue_IntVal:
			if value.GetIntVal() == -1 {
				columnValues = append(columnValues, "\nil")
			} else {
				columnValues = append(columnValues, fmt.Sprintf("%d", value.GetIntVal()))
			}
		case *metricproto.MetricValue_FloatVal:
			columnValues = append(columnValues, fmt.Sprintf("%f", value.GetFloatVal()))
		}
	}

	err := consumer.mapBi[table].Insert(columnValues)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to Insert from plugin %s", msg.PluginName)
	}

}

*/
