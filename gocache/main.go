package main

import (
	"database/sql"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/protobuf/proto"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	metricproto "github.com/selvakumarjawahar/myexperiments/gocache/gen"
	met "github.com/selvakumarjawahar/myexperiments/gocache/internal/pkg/metrics"
	"github.com/selvakumarjawahar/myexperiments/gocache/internal/pkg/plugins"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <group> <topic> \n",
			os.Args[0])
		os.Exit(1)
	}
	broker := os.Args[1]
	group  := os.Args[2]
	topic  := os.Args[3]

	kConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"broker.address.family": "v4",
		"group.id":              group,
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest"})

	defer kConsumer.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created Consumer %v\n", kConsumer)
	err = kConsumer.Subscribe(topic, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe topic: %s\n", err)
		os.Exit(1)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "netrounds", "postgres")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open DB")
		os.Exit(1)
	}
	log.Info().Msg("DB Opened successfully")
	defer db.Close()

	err = CreateDBTables(db,"monitor_metrics_http")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to Create table monitor_metrics_http ")
		os.Exit(1)
	}
	err = CreateDBTables(db,"monitor_metrics_minute1_http")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to Create table monitor_metrics_minute1_http ")
		os.Exit(1)
	}
	err = CreateDBTables(db,"monitor_metrics_minute5_http")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to Create table monitor_metrics_minute5_http ")
		os.Exit(1)
	}
	err = CreateDBTables(db,"monitor_metrics_minute30_http")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to Create table monitor_metrics_minute30_http ")
		os.Exit(1)
	}
	err = CreateDBTables(db,"monitor_metrics_minute60_http")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to Create table monitor_metrics_minute60_http ")
		os.Exit(1)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	aggregator := plugins.NewAggregatorMap()
	aggregator.SetAggregators(map[string]plugins.Aggregator{
		"http.ConnectTimeAvg" : plugins.Avg,
		"http.Es" : plugins.Sum,
		"http.EsResponse" : plugins.Min,
		"http.EsTimeout" : plugins.Max,
		"http.FirstByteTimeAvg" : plugins.Avg,
		"http.ResponseTimeAvg" : plugins.Avg,
		"http.ResponseTimeMax" : plugins.Max,
		"http.ResponseTimeMin" : plugins.Min,
		"http.SizeAvg" : plugins.Avg,
		"http.SpeedAvg" : plugins.Avg})

	rollUpState := met.NewRollUpStateManager(db,aggregator)

	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := kConsumer.Poll(100)
			if ev == nil {
				continue
			}
			switch event := ev.(type) {
			case *kafka.Message:
				metrics := &metricproto.Metrics{}
				err = proto.Unmarshal(event.Value, metrics)
				if err != nil {
					log.Error().Err(err).Msgf("Failed to unmarshal data from metrics topic %v", topic)
					continue
				}
				pluginMetric := new(met.MetricMsg)
				pluginMetric.Metrics = metrics
				pluginMetric.PluginName = "http"
				rollUpState.OnMessage(pluginMetric)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", event.Code(), event)
				if event.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", event)
			}
		}
	}
}


func CreateDBTables(db *sql.DB,tableName string)error {
	var query string
	query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s " +
		"(time TIMESTAMPTZ NOT NULL, " +
		"measurement_id INT NOT NULL, " +
		"ConnectTimeAvg FLOAT, " +
		"Es INT, " +
		"EsResponse INT, " +
		"EsTimeout INT, " +
		"FirstByteTimeAvg FLOAT, " +
		"ResponseTimeAvg FLOAT, " +
		"ResponseTimeMax FLOAT, " +
		"ResponseTimeMin FLOAT, " +
		"SizeAvg FLOAT, " +
		"SpeedAvg FLOAT)", tableName)
	log.Debug().Msgf("%s", query)
	_, err := db.Exec(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return err
	}

	// Create views joining metrics to measurements.
	// New columns may only appear at the end of the view - otherwise REPLACE will not help here.
	// Drop will be safer, it will happen only on schema changes so this operation should be safe.
	/*
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
	*/


	// Create index
	//query = fmt.Sprintf(
	//	"CREATE INDEX IF NOT EXISTS %s_measurement_id_idx ON %s (measurement_id)",
	//	tableName,
	//	tableName)
	//log.Debug().Msgf("%s", query)
	//_, err = db.Exec(query)
	//if err != nil {
	//	log.Error().Err(err).Msg("Failed to execute query")
	//	return err
	//}
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
	log.Info().Msgf("Created Table plugin %s", tableName)
	return nil

}