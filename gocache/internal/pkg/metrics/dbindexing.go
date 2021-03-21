package metrics

/*
func DBIndexer(db *sql.DB, dbPipe <-chan RollupState) {
	for data := range dbPipe {
		fmt.Printf("%s, %s, %d\n",
			data.granularity.windowName,
			data.msg.PluginName,
			data.msg.Metrics.StreamId,
			//data.msg.Metrics.Timestamp.Seconds,
			//data.msg.Metrics.Values
			)
	}
}

*/

import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/selvakumarjawahar/myexperiments/gocache/cmd/bulkinsert"
	metricproto "github.com/selvakumarjawahar/myexperiments/gocache/gen"
	"github.com/selvakumarjawahar/myexperiments/gocache/internal/pkg/plugins"
	"time"
)

//The below map enables to config the bulk inserter with different
//batch sizes and batchage for different granularity
//Currently all these values are same for all granularity but
//This can be changed
var dbConfigMap = map[string]struct {
	maxBatchSize int //Message Count
	maxBatchAge  int //seconds
}{
	plugins.Minute1:  {100000, 30},
	plugins.Minute5:  {100000, 30},
	plugins.Minute30: {100000, 30},
	plugins.Minute60: {100000, 30},
}

//This function is another goroutine reads rolled up metrics and indexes into appropriate granularity tables
func DBIndexer(db *sql.DB, dbPipe <-chan RollupState) {

	dbInsertMap := make(map[string]*bulkinsert.BulkInsert)
	for data := range dbPipe {
		tableName := fmt.Sprintf("monitor_metrics_%s_%s", data.granularity.windowName, data.msg.PluginName)
		columnNames, err := plugins.GetColumnsForTable(tableName, db)
		if len(columnNames) == 0 || err != nil {
			log.Error().Err(err).Msgf("Table %s does not exist ", tableName)
			continue
		}

		if _, ok := dbInsertMap[tableName]; !ok { //New DB Bulk inserter
			dbInsertMap[tableName] = bulkinsert.NewBulkInsert(
				db,
				tableName,
				dbConfigMap[data.granularity.windowName].maxBatchSize,
				dbConfigMap[data.granularity.windowName].maxBatchAge,
				columnNames)
		}
		metricNames := []string{
			"ConnectTimeAvg",
			"Es",
			"EsResponse",
			"EsTimeout",
			"FirstByteTimeAvg",
			"ResponseTimeAvg",
			"ResponseTimeMax",
			"ResponseTimeMin",
			"SizeAvg",
			"SpeedAvg",
		}
		columnValues := metricsToString(metricNames,data.msg)
		if err := dbInsertMap[tableName].Insert(columnValues); err != nil {
			log.Error().Err(err).Msgf("Error Inserting column values to table %s ", tableName)
		}
	}
}

//This is a helper function , which converts the metrics to string, which can be then used in db insert query
func metricsToString(columnNames []string, msg *MetricMsg) []string {

	columnValues := []string{fmt.Sprintf(
		"'%s'",
		time.Unix(msg.Metrics.Timestamp.Seconds, 0).UTC().Format("2006-01-02 15:04:05")),
		fmt.Sprintf("%d", msg.Metrics.StreamId)}



	for _, metricName := range columnNames {
		if metricName == "time" || metricName == "measurement_id" {
			continue
		}

		value := msg.Metrics.Values[metricName]
		// add \nil for all missing values
		if value == nil {
			columnValues = append(columnValues, "\nil")
			log.Debug().Msgf("setting metric %s value nil",metricName)
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
	return columnValues
}

