package plugins

import (
	"math"
	"testing"
)

func TestMin(t *testing.T) {
	testCases := []struct {
		input1   interface{}
		input2   interface{}
		expected interface{}
	}{
		{float32(45.56), float32(40.23), float32(40.23)},
		{int64(1), int64(5), int64(1)},
	}
	for index, testCase := range testCases {
		output := Min(testCase.input1, testCase.input2)
		if testCase.expected != output {
			t.Errorf(""+
				"For Test Index = %d The output value expected was %v but got %v \n",
				index, testCase.expected, output)
		}
	}
}

func TestMax(t *testing.T) {
	testCases := []struct {
		input1   interface{}
		input2   interface{}
		expected interface{}
	}{
		{float32(45.56), float32(40.23), float32(45.56)},
		{int64(1), int64(5), int64(5)},
	}
	for index, testCase := range testCases {
		output := Max(testCase.input1, testCase.input2)
		if testCase.expected != output {
			t.Errorf(""+
				"For Test Index = %d The output value expected was %v but got %v \n",
				index, testCase.expected, output)
		}
	}
}

func TestAvg(t *testing.T) {
	testCases := []struct {
		input1   interface{}
		input2   interface{}
		expected interface{}
	}{
		{float32(45.56), float32(40.23), float32(43)},
		{int64(1), int64(5), int64(3)},
	}
	for index, testCase := range testCases {
		switch testCase.expected.(type) {
		case int64:
			output := Avg(testCase.input1, testCase.input2)
			if testCase.expected != output {
				t.Errorf(""+
					"For Test Index = %d The output value expected was %v but got %v \n",
					index, testCase.expected, output)
			}
		case float32:
			output := float32(math.Round(float64(Avg(testCase.input1, testCase.input2).(float32))))
			if testCase.expected != output {
				t.Errorf(""+
					"For Test Index = %d The output value expected was %v but got %v \n",
					index, testCase.expected, output)
			}
		default:
			t.Error("Error unknown type")
		}
	}
}

func TestSum(t *testing.T) {
	testCases := []struct {
		input1   interface{}
		input2   interface{}
		expected interface{}
	}{
		{float32(45.56), float32(40.23), float32(86)},
		{int64(1), int64(5), int64(6)},
	}

	for index, testCase := range testCases {
		switch testCase.expected.(type) {
		case int64:
			output := Sum(testCase.input1, testCase.input2)
			if testCase.expected != output {
				t.Errorf(""+
					"For Test Index = %d The output value expected was %v but got %v \n",
					index, testCase.expected, output)
			}
		case float32:
			output := float32(math.Round(float64(Sum(testCase.input1, testCase.input2).(float32))))
			if testCase.expected != output {
				t.Errorf(""+
					"For Test Index = %d The output value expected was %v but got %v \n",
					index, testCase.expected, output)
			}
		default:
			t.Error("Error unknown type")
		}
	}
}

func TestAggregatorMap_SetAggregators(t *testing.T) {
	metricMap := map[string]Aggregator{
		"es":          Sum,
		"min_latency": Min,
		"max_latency": Max,
		"avg_latency": Avg,
	}
	aggregatorMap := NewAggregatorMap()
	aggregatorMap.SetAggregators(metricMap)

	testCases := []struct {
		metricName     string
		metricInput1   interface{}
		metricInput2   interface{}
		expectedOutput interface{}
	}{
		{"es", int64(23), int64(56), int64(79)},
		{"min_latency", int64(23), int64(56), int64(23)},
		{"max_latency", int64(23), int64(56), int64(56)},
		{"avg_latency", int64(23), int64(56), int64(39)},
	}

	for index, testCase := range testCases {
		actualOutput := aggregatorMap.GetAggregator(testCase.metricName)(testCase.metricInput1, testCase.metricInput2)
		if testCase.expectedOutput != actualOutput {
			t.Errorf("The test case index %d failed, expected output %v actual output %v \n",
				index, testCase.expectedOutput, actualOutput)
		}
	}
}

func TestAggregatorMap_SetAggregator(t *testing.T) {

	aggregatorMap := NewAggregatorMap()
	aggregatorMap.SetAggregator("es", Sum)

	testCases := []struct {
		metricName     string
		metricInput1   interface{}
		metricInput2   interface{}
		expectedOutput interface{}
	}{
		{"es", int64(23), int64(56), int64(79)},
	}

	for index, testCase := range testCases {
		actualOutput := aggregatorMap.GetAggregator(testCase.metricName)(testCase.metricInput1, testCase.metricInput2)
		if testCase.expectedOutput != actualOutput {
			t.Errorf("The test case index %d failed, expected output %v actual output %v \n",
				index, testCase.expectedOutput, actualOutput)
		}
	}
}
