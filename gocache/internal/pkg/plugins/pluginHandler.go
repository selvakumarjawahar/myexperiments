package plugins

import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
)

/*
import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"sort"
)
*/
//Constants defining granularity
const (
	Minute1 = "minute1"
	Minute5 = "minute5"
	Minute30 = "minute30"
	Minute60 = "minute60"
)

func GetColumnsForTable(tableName string, db *sql.DB) ([]string, error) {
	var existingColumns []string
	// check if table exists in the database
	query := fmt.Sprintf(
		"SELECT column_name FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '%s'",
		tableName)
	log.Debug().Msgf("%s", query)
	rows, err := db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return existingColumns, err
	}
	defer rows.Close()
	var columnName string
	for rows.Next() {
		rows.Scan(&columnName)
		existingColumns = append(existingColumns, columnName)
	}
	return existingColumns, nil
}

/*
import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	//pluginproto "github.com/Netrounds/ncc3-services/ncc3-metrics-service/gen/protobuf/plugin"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
)

//Constants defining granularity
//const (
//	Minute1 = "minute1"
//	Minute5 = "minute5"
//	Minute30 = "minute30"
//	Minute60 = "minute60"
//)

/*
// PluginHandler is object for handling plugins
type PluginHandler struct {
	pluginColumns map[string]bool
	plugins       map[string]*pluginproto.Plugin
	channel       chan *pluginproto.Plugin
	AggregatorMap *AggregatorMap
}

// NewPluginHandler creates a new handler for plugins
func NewPluginHandler(channel chan *pluginproto.Plugin) *PluginHandler {

	return &PluginHandler{
		pluginColumns: make(map[string]bool),
		plugins:       make(map[string]*pluginproto.Plugin),
		channel:       channel,
		AggregatorMap: NewAggregatorMap(),
	}
}

// HasMetric checks if metric with the key exists in plugin definition
func (handler *PluginHandler) HasMetric(pluginName string, metricName string) bool {
	key := pluginName + "." + metricName
	return handler.pluginColumns[key]
}

func resetOffsets(consumer *kafka.Consumer, topics[] string) error{
	var err error
	for {
		for _, topic := range topics {
			_, err = consumer.CommitOffsets([]kafka.TopicPartition{{
				Topic:     &topic,
				Partition: int32(0),
				Offset:    kafka.Offset(0),
			}})
			if err != nil {
				log.Debug().Err(err).Msgf("failed to reset offset for topic %s", topic)
				break
			}
		}
		if err == nil {
			return nil
		}
	}
	
}
// ConsumePlugins consumes plugins topic from kafka and creates metrics database tables according
func (handler *PluginHandler) ConsumePlugins(db *sql.DB, topics []string, kafkaURL string) {
	// Create new consumer
	log.Info().Strs("topics", topics).Msg("creating consumer")

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaURL,
		"group.id":          "metrics-service",
		// Enable generation of PartitionEOF when the
		// end of a partition is reached.
		"enable.partition.eof": true,
		// auto commit offsets
		"enable.auto.commit": true,
		//start from beginning
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Panic().Err(err).Msg("failed to create Kafka consumer for plugins topic")
	}
	defer consumer.Close()

	
	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Panic().Err(err).Msg("failed to subscribe to metadata topic")
	}

	resetOffsets(consumer, topics)


	var msgCount = 0
	lastReport := time.Now()
	run := true
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for run == true {
		select {
		case sig := <-sigchan:
			log.Info().Str("signal", sig.String()).Msg("caught terminating signal")
			run = false
		default:
			consumerEvent := consumer.Poll(100)
			if consumerEvent == nil {
				continue
			}

			log.Debug().Msgf("Event: %v", consumerEvent)
			switch event := consumerEvent.(type) {
			case kafka.PartitionEOF:
				log.Debug().Msgf("%% Reached %v", event)
				handler.channel <- nil

			case *kafka.Message:
				plugin := &pluginproto.Plugin{}
				syncSchema := false
				if err = proto.Unmarshal(event.Value, plugin); err != nil {
					log.Error().Err(err).Msg("Failed to unmarshal data from plugins topic")
					continue
				}
				log.Debug().Msgf("Received plugin %s", plugin.Name)
				handler.plugins[plugin.Name] = plugin
				for _, metric := range plugin.Metrics {
					field := plugin.Name + "." + metric.Name
					// If metric does not exists in cache then syncPluginWithDB will be called.
					// Tt means it will call sync function when cache is empty, but there is a check
					// inside syncPluginWithDB as well, which will not execute create/alter statement
					if !handler.HasMetric(plugin.Name, metric.Name) {
						syncSchema = true
						// Fill in the map of metric field and associated aggregator function
						switch metric.GetAggregation() {
						case pluginproto.PluginMetric_MIN:
							handler.AggregatorMap.SetAggregator(field,Min)
						case pluginproto.PluginMetric_MAX:
							handler.AggregatorMap.SetAggregator(field,Max)
						case pluginproto.PluginMetric_SUM:
							handler.AggregatorMap.SetAggregator(field,Sum)
						case pluginproto.PluginMetric_AVG:
							handler.AggregatorMap.SetAggregator(field,Avg)
						default:
							handler.AggregatorMap.SetAggregator(field,Min)
						}
					}
					handler.pluginColumns[field] = true
				}

				if syncSchema {
					if err = syncPluginWithDB(plugin, db); err != nil {
						log.Error().Err(err).Msg("Failed to add plugins to DB")
						continue
					}
					if err = syncPluginWithGranularityDB(plugin, db); err != nil {
						log.Error().Err(err).Msg("Failed to add plugins to Granularity DB")
						continue
					}
				}

				handler.channel <- plugin

				elapsed := time.Since(lastReport)
				if elapsed.Seconds() > 1.0 {
					log.Debug().Msgf("%v: consumed %d messages, %.0f per second", topics, msgCount, float64(msgCount)/elapsed.Seconds())
					msgCount = 0
					lastReport = time.Now()
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
}


// GetColumnsForTable returns array of columns for tableName or empty array if it does not exist
func GetColumnsForTable(tableName string, db *sql.DB) ([]string, error) {
	var existingColumns []string
	// check if table exists in the database
	query := fmt.Sprintf(
		"SELECT column_name FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '%s'",
		tableName)
	log.Debug().Msgf("%s", query)
	rows, err := db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return existingColumns, err
	}
	defer rows.Close()
	var columnName string
	for rows.Next() {
		rows.Scan(&columnName)
		existingColumns = append(existingColumns, columnName)
	}
	return existingColumns, nil
}

/*
// getColumnDefinition returns db column definition
func getColumnDefinition(name string, metricType pluginproto.PluginMetric_Type) string {
	switch metricType {
	case pluginproto.PluginMetric_INT:
		return fmt.Sprintf("%s INT", name)
	default:
		return fmt.Sprintf("%s FLOAT", name)
	}
}

func columnExists(columns []string, newColumn string) bool {
	sort.Strings(columns)
	i := sort.SearchStrings(columns, newColumn)
	return i < len(columns) && columns[i] == newColumn
}

// syncTable will create or update plugin tables with given name, views and indexes
// The assumption is that the only allowed operation in case of updating is adding new column
func syncTable(plugin *pluginproto.Plugin, db *sql.DB, tableName string) error {

	existingColumns, err := GetColumnsForTable(tableName, db)
	if err != nil {
		return err
	}

	var query string
	var columns []string

	if len(existingColumns) == 0 {
		// table does not exist so it has to be created
		columns = append(columns, "time TIMESTAMPTZ NOT NULL")
		columns = append(columns, "measurement_id INT NOT NULL")
		for _, metric := range plugin.Metrics {
			columns = append(columns, getColumnDefinition(metric.Name, metric.Type))
		}
		query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columns, ","))
	} else {
		// table exists, search for new columns and add them to the query
		for _, metric := range plugin.Metrics {
			if !columnExists(existingColumns, metric.Name) {
				columns = append(columns, fmt.Sprintf("ADD COLUMN %s", getColumnDefinition(metric.Name, metric.Type)))
			}
		}
		// if no columns have been changed then skip the rest of the function
		if len(columns) == 0 {
			return nil
		}
		query = fmt.Sprintf("ALTER TABLE %s %s", tableName, strings.Join(columns, ","))
	}

	// sync table
	log.Debug().Msgf("%s", query)
	_, err = db.Exec(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return err
	}

	// Create views joining metrics to measurements.
	// New columns may only appear at the end of the view - otherwise REPLACE will not help here.
	// Drop will be safer, it will happen only on schema changes so this operation should be safe.
	query = fmt.Sprintf(
		"DROP VIEW IF EXISTS vw_%s;"+
			"CREATE VIEW vw_%s AS SELECT measurements.*, metrics.*"+
			"FROM %s metrics INNER JOIN measurements ON (metrics.measurement_id=measurements.id)",
		tableName,
		tableName,
		tableName)
	log.Debug().Msgf("%s", query)
	_, err = db.Exec(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return err
	}

	// Create index
	query = fmt.Sprintf(
		"CREATE INDEX IF NOT EXISTS %s_measurement_id_idx ON %s (measurement_id)",
		tableName,
		tableName)
	log.Debug().Msgf("%s", query)
	_, err = db.Exec(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return err
	}
	// Create hypertable, it will be automatically updated on altering base table
	query = fmt.Sprintf(
		"SELECT create_hypertable('%s', 'time', if_not_exists => true, chunk_time_interval => 3600000000, migrate_data => true)",
		tableName)
	log.Debug().Msgf("%s", query)
	_, err = db.Exec(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return err
	}
	log.Info().Msgf("Loaded plugin %s", plugin.Name)
	return nil

}

func syncPluginWithDB(plugin *pluginproto.Plugin, db *sql.DB) error {
	tableName := fmt.Sprintf("monitor_metrics_%s", plugin.Name)
	return syncTable(plugin,db,tableName)
}

func syncPluginWithGranularityDB(plugin *pluginproto.Plugin, db *sql.DB) error {

	//Create/Update Granularity DB - 1min
	tableName := fmt.Sprintf("monitor_metrics_%s_%s",Minute1, plugin.Name)
	err := syncTable(plugin,db,tableName)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create %s",tableName)
		return err
	}
	//Create/Update Granularity DB - 5min
	tableName = fmt.Sprintf("monitor_metrics_%s_%s", Minute5, plugin.Name)
	err = syncTable(plugin,db,tableName)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create %s",tableName)
		return err
	}
	//Create/Update Granularity DB - 30min
	tableName = fmt.Sprintf("monitor_metrics_%s_%s", Minute30, plugin.Name)
	err = syncTable(plugin,db,tableName)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create %s",tableName)
		return err
	}
	//Create/Update Granularity DB - 60min
	tableName = fmt.Sprintf("monitor_metrics_%s_%s", Minute60, plugin.Name)
	err = syncTable(plugin,db,tableName)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create %s",tableName)
		return err
	}
	return nil
}

*/