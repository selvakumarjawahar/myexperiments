package metrics

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/rs/zerolog/log"
	metricproto "github.com/selvakumarjawahar/myexperiments/gocache/gen"
	"github.com/selvakumarjawahar/myexperiments/gocache/internal/pkg/plugins"
	"math"
	"time"
)

type MetricMsg struct {
	PluginName string
	Metrics    *metricproto.Metrics
}

//granularityWindow struct stores the time window for each granularity defined in the system.
//The timeWindow field saves the time window in seconds. The window name field is a enum which
//has the infix of granularity table name
type granularityWindow struct {
	timeWindow int //seconds
	windowName string
}

//rollupState struct stores the rolled up metric for a given granularity window
type RollupState struct {
	granularity granularityWindow
	msg *MetricMsg
}

//measurement struct stores list of rollup states for different granularity window
//It also saves the latest timestamp which was saved in rollups
type measurement struct {
	maxTime int64
	rollUps []RollupState
}

const LowestGranularity = 3600

//RollUpStateManager has all the state value of rollups along with methods to operate on
//the states
type RollUpStateManager struct {
	measurements map[int32]*measurement
	funcMap *plugins.AggregatorMap
	dbChannel chan RollupState
	db *sql.DB
}

func NewRollUpStateManager(db *sql.DB, aggregatorMap *plugins.AggregatorMap)*RollUpStateManager{
	rollupState := new(RollUpStateManager)
	rollupState.measurements = make(map[int32]*measurement)
	rollupState.db = db
	rollupState.dbChannel = make(chan RollupState)
	rollupState.funcMap = aggregatorMap
	go DBIndexer(rollupState.db,rollupState.dbChannel)
	return rollupState
}

func (rsm *RollUpStateManager ) InitRollUpStates()error{
	availableMeasurements,err := rsm.getMeasurementIDs()
	if err != nil {
		log.Error().Err(err).Msg("error in fetching metric Ids")
		return err
	}
	for key,value :=  range availableMeasurements {
		rawTableName := fmt.Sprintf("monitor_metrics_%s", value)
		//columnNames,err := plugins.GetColumnsForTable(rawTableName,rsm.db)
		columnNames := []string{
			"es",
			"timeout",
		}
		if err != nil {
			log.Error().Err(err).Msg("error in fetching column name")
			continue
		}
		dbrecord := NewDBRecord(value,columnNames,rsm.db)
		maxTimestamp,err := rsm.getMaxTimeStamp(rawTableName,key)
		if err != nil {
			log.Error().Err(err).Msgf("error in fetching max timestamp for measurement %d",key)
			continue
		}
		earliestTimestamp := (maxTimestamp / LowestGranularity) * LowestGranularity
		messages,err := dbrecord.GetMetricsFromTimeRange(earliestTimestamp,key)
		if err != nil {
			log.Error().Err(err).Msgf("error in fetching db records")
			continue
		}
		for _,message := range messages {
			rsm.OnMessage(message)
		}
	}
	return nil
}

func (rsm *RollUpStateManager ) OnMessage(incomingMsg *MetricMsg) {
	if _,ok := rsm.measurements[incomingMsg.Metrics.StreamId]; !ok {
		rsm.addNewMeasurement(incomingMsg)
	} else {
		rsm.updateRollUps(incomingMsg)
	}
}

//NewMetricMessage Does a deep copy of metric msg
func NewMetricMessage(msg *MetricMsg)*MetricMsg{
	newMsg := new(MetricMsg)
	newMsg.PluginName = msg.PluginName
	newMsg.Metrics = new(metricproto.Metrics)
	newMsg.Metrics.StreamId = msg.Metrics.StreamId
	newMsg.Metrics.Timestamp = msg.Metrics.Timestamp
	newMsg.Metrics.Values = make(map[string]*metricproto.MetricValue)
	for key,val := range msg.Metrics.Values {
		newMsg.Metrics.Values[key] = new(metricproto.MetricValue)
		*newMsg.Metrics.Values[key] = *val
	}
	return newMsg
}

func (rsm *RollUpStateManager ) addNewMeasurement(incomingMsg *MetricMsg) {
	newMeasurement := new(measurement)
	newMeasurement.rollUps = make([]RollupState,4)
	newMeasurement.rollUps[0] = RollupState{
		granularity: granularityWindow{
			timeWindow: 60,
			windowName: plugins.Minute1,
		},
		msg:         NewMetricMessage(incomingMsg),
	}
	newMeasurement.rollUps[1] = RollupState{
		granularity: granularityWindow{
			timeWindow: 300,
			windowName: plugins.Minute5,
		},
		msg:         NewMetricMessage(incomingMsg),
	}
	newMeasurement.rollUps[2] = RollupState{
		granularity: granularityWindow{
			timeWindow: 1800,
			windowName: plugins.Minute30,
		},
		msg:         NewMetricMessage(incomingMsg),
	}
	newMeasurement.rollUps[3] = RollupState{
		granularity: granularityWindow{
			timeWindow: 3600,
			windowName: plugins.Minute60,
		},
		msg:         NewMetricMessage(incomingMsg),
	}
	newMeasurement.maxTime = incomingMsg.Metrics.Timestamp.Seconds
	rsm.measurements[incomingMsg.Metrics.StreamId] = newMeasurement
}

func (rsm *RollUpStateManager ) resetMeasurement(timeWindowIndex int,incomingMsg *MetricMsg){
	rsm.measurements[incomingMsg.Metrics.StreamId].rollUps[timeWindowIndex].msg = NewMetricMessage(incomingMsg)
}

func (rsm *RollUpStateManager ) computeRollUp(incomingMetric *MetricMsg) {

	for _,rollup := range rsm.measurements[incomingMetric.Metrics.StreamId].rollUps{
		for key, value := range incomingMetric.Metrics.Values {
			field := incomingMetric.PluginName + "." + key
			aggregator := rsm.funcMap.GetAggregator(field)
			switch value.Type.(type) {
			case *metricproto.MetricValue_IntVal:
				protoVal := aggregator(rollup.msg.Metrics.Values[key].GetIntVal(), value.GetIntVal()).(int64)
				rollup.msg.Metrics.Values[key] = &metricproto.MetricValue{
					Type: &metricproto.MetricValue_IntVal{protoVal}}

			case *metricproto.MetricValue_FloatVal:
				protoVal := aggregator(rollup.msg.Metrics.Values[key].GetFloatVal(), value.GetFloatVal()).(float32)
				rollup.msg.Metrics.Values[key] = &metricproto.MetricValue{
					Type: &metricproto.MetricValue_FloatVal{protoVal}}
			}
		}
	}
}

func (rsm *RollUpStateManager ) processGranularity(incomingMsg *MetricMsg) {

	for index,rollup := range rsm.measurements[incomingMsg.Metrics.StreamId].rollUps {
		bucket := math.Remainder(float64(incomingMsg.Metrics.Timestamp.Seconds), float64(rollup.msg.Metrics.Timestamp.Seconds)) //*
		if int(bucket) >= rollup.granularity.timeWindow {
			sendRollup := RollupState{
				granularity: rollup.granularity,
				msg:         NewMetricMessage(rollup.msg),
			}
			rsm.dbChannel <- sendRollup
			rsm.resetMeasurement(index,incomingMsg)
		}
//		if index == 0 {
//			fmt.Printf("bucket = %f origEs = %f rollupEs = %f \n",bucket,
//				incomingMsg.Metrics.Values["SizeAvg"].GetFloatVal(),
//				rsm.measurements[incomingMsg.Metrics.StreamId].rollUps[index].msg.Metrics.Values["SizeAvg"].GetFloatVal())
//		}
	}
}

func (rsm *RollUpStateManager ) updateRollUps(incomingMsg *MetricMsg) {
	measurement := rsm.measurements[incomingMsg.Metrics.StreamId]
	if measurement.maxTime > incomingMsg.Metrics.Timestamp.Seconds {
		log.Info().Msg("Duplicate message skipping")
		return
	}
	//for index,val := range measurement.rollUps{
	//	fmt.Printf("rollup[%d]=%d,%s \n",index,val.granularity.timeWindow,val.granularity.windowName)
	//}
	rsm.processGranularity(incomingMsg)
	rsm.computeRollUp(incomingMsg)
}

func (rsm *RollUpStateManager ) getMeasurementIDs()(map[int32]string,error){
	measurementPluginMap := make(map[int32]string)
	query := "SELECT DISTINCT ON (task_id) task_id, plugin_name FROM measurement ORDER BY task_id, plugin_name"
	log.Debug().Msgf("%s", query)
	rows, err := rsm.db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return nil, err
	}
	defer rows.Close()
	var pluginName string
	var measurementID int32
	for rows.Next() {
		if err := rows.Scan(&measurementID,&pluginName); err != nil {
			log.Error().Err(err).Msg("Error Scanning the row")
			return nil, err
		}
		measurementPluginMap[measurementID] = pluginName
	}
	return measurementPluginMap,nil
}

func (rsm *RollUpStateManager )getMaxTimeStamp(tableName string,measurementID int32) (int64, error) {
	query := fmt.Sprintf("SELECT max(time) FROM %s WHERE measurement_id = %d", tableName, measurementID)
	log.Debug().Msgf("%s", query)
	rows, err := rsm.db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return 0, err
	}
	defer rows.Close()
	var maxTime time.Time
	rows.Scan(&maxTime)
	return maxTime.Unix(), nil
}

type DBRecord struct {
	pluginName string
	ColumnNames []string
	db *sql.DB
	RowData map[int][]interface{}
}

func NewDBRecord(pluginName string, columnNames []string,db *sql.DB) *DBRecord {
	dbRecord := new(DBRecord)
	dbRecord.ColumnNames = make([]string, len(columnNames))
	dbRecord.RowData = make(map[int][]interface{})
	dbRecord.pluginName = pluginName
	dbRecord.db = db
	for index, value := range columnNames {
		dbRecord.ColumnNames[index] = value
	}
	return dbRecord
}

func (dbRecord *DBRecord) GetMetricsFromTimeRange(rangeStart int64, measurementID int32)([]*MetricMsg,error) {

	timeStart := time.Unix(rangeStart, 0).UTC().Format("2006-01-02 15:04:05")
	tableName := fmt.Sprintf("monitor_metrics_%s", dbRecord.pluginName)
	query := fmt.Sprintf("SELECT * FROM %s WHERE time > %s AND measurement_id = %d", tableName, timeStart, measurementID)
	log.Debug().Msgf("%s", query)
	rows, err := dbRecord.db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return  nil,err
	}
	defer rows.Close()
	temp := make([]interface{}, len(dbRecord.ColumnNames))
	count := 0
	for rows.Next() {
		if err := rows.Scan(temp); err != nil {
			log.Error().Err(err).Msg("Error Scanning the row")
			return nil,err
		}
		dbRecord.RowData[count] = temp
		count++
	}
	messages , err := dbRecord.packageMetricsFromRows()
	if err := rows.Scan(temp); err != nil {
		log.Error().Err(err).Msg("Error Scanning the row")
		return nil,err
	}
	return messages,nil

}

func (dbRecord *DBRecord) packageMetricsFromRows()([]*MetricMsg, error) {
	result := make([]*MetricMsg, len(dbRecord.RowData))
	if len(dbRecord.ColumnNames) > 1 && (dbRecord.ColumnNames[0] != "time" && dbRecord.ColumnNames[1] != "measurement_id") {
		log.Error().Msg("Coloumn names not proper")
		return nil, errors.New("column names not proper ")
	}
	for _, value := range dbRecord.RowData {
		tempProto := new(metricproto.Metrics)
		tempMetrics := new(MetricMsg)
		for index, val := range value {
			if index == 0 {
				timestamp := new(timestamp.Timestamp)
				timestamp.Seconds = val.(time.Time).Unix()
				tempProto.Timestamp = timestamp
				continue
			}
			if index == 1 {
				tempProto.StreamId = val.(int32)
				continue
			}
			switch val.(type) {
			case int64:
				tempProto.Values[dbRecord.ColumnNames[index]] = &metricproto.MetricValue{
					Type: &metricproto.MetricValue_IntVal{val.(int64)}}

			case float32:
				tempProto.Values[dbRecord.ColumnNames[index]] = &metricproto.MetricValue{
					Type: &metricproto.MetricValue_FloatVal{val.(float32)}}

			default:
				continue
			}
		}
		tempMetrics.Metrics = tempProto
		tempMetrics.PluginName = dbRecord.pluginName
		result = append(result, tempMetrics)
	}
	return result, nil
}

