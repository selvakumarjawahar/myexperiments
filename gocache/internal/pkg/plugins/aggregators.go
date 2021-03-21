package plugins

type Aggregator func(currentMetric, incomingMetric interface{}) interface{}

func Min(currentMetric, incomingMetric interface{}) interface{} {
	var min interface{}
	switch currentMetric.(type) {
	case int64:
		if currentMetric.(int64) < incomingMetric.(int64) {
			min = currentMetric.(int64)
		} else {
			min = incomingMetric.(int64)
		}
	case float32:
		if currentMetric.(float32) < incomingMetric.(float32) {
			min = currentMetric.(float32)
		} else {
			min = incomingMetric.(float32)
		}
	default:
	}
	return min
}

func Max(currentMetric, incomingMetric interface{}) interface{} {
	var max interface{}
	switch currentMetric.(type) {
	case int64:
		if currentMetric.(int64) > incomingMetric.(int64) {
			max = currentMetric.(int64)
		} else {
			max = incomingMetric.(int64)
		}
	case float32:
		if currentMetric.(float32) > incomingMetric.(float32) {
			max = currentMetric.(float32)
		} else {
			max = incomingMetric.(float32)
		}
	default:
	}
	return max
}

func Sum(currentMetric, incomingMetric interface{}) interface{} {
	var sum interface{}
	switch currentMetric.(type) {
	case int64:
		sum = currentMetric.(int64) + incomingMetric.(int64)
	case float32:
		sum = currentMetric.(float32) + incomingMetric.(float32)
	default:
	}
	return sum
}

func Avg(currentMetric, incomingMetric interface{}) interface{} {
	var avg interface{}
	switch currentMetric.(type) {
	case int64:
		avg = (currentMetric.(int64) + incomingMetric.(int64)) / 2
	case float32:
		avg = (currentMetric.(float32) + incomingMetric.(float32)) / 2
	default:
	}
	return avg
}

type AggregatorMap struct {
	metricAggregators map[string]Aggregator
}

func NewAggregatorMap() *AggregatorMap {
	aMap := new(AggregatorMap)
	aMap.metricAggregators = make(map[string]Aggregator)
	return aMap
}

func (aMap *AggregatorMap) CopyAggregatorMap()*AggregatorMap{
	copyMap := new(AggregatorMap)
	for key,value :=  range aMap.metricAggregators {
		copyMap.metricAggregators[key] = value
	}
	return copyMap
}

func (aMap *AggregatorMap) GetAggregator(metricName string) Aggregator {
	return aMap.metricAggregators[metricName]
}

func (aMap *AggregatorMap) SetAggregator(metricName string, aggregator Aggregator) {
	aMap.metricAggregators[metricName] = aggregator
}

func (aMap *AggregatorMap) SetAggregators(metricMap map[string]Aggregator) {
	for key, _ := range metricMap {
		aMap.metricAggregators[key] = metricMap[key]
	}
}
