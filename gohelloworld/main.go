
package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Metrics struct {
	ConnectTimeAvg   *float32 `protobuf:"fixed32,1,opt,name=connect_time_avg,json=connectTimeAvg" json:"connect_time_avg,omitempty"`
	FirstByteTimeAvg *float32 `protobuf:"fixed32,2,opt,name=first_byte_time_avg,json=firstByteTimeAvg" json:"first_byte_time_avg,omitempty"`
	ResponseTimeMin  *float32 `protobuf:"fixed32,3,opt,name=response_time_min,json=responseTimeMin" json:"response_time_min,omitempty"`
	ResponseTimeAvg  *float32 `protobuf:"fixed32,4,opt,name=response_time_avg,json=responseTimeAvg" json:"response_time_avg,omitempty"`
	ResponseTimeMax  *float32 `protobuf:"fixed32,5,opt,name=response_time_max,json=responseTimeMax" json:"response_time_max,omitempty"`
	SizeAvg          *float32 `protobuf:"fixed32,6,opt,name=size_avg,json=sizeAvg" json:"size_avg,omitempty"`
	SpeedAvg         *float32 `protobuf:"fixed32,7,opt,name=speed_avg,json=speedAvg" json:"speed_avg,omitempty"`
	EsTimeout        *int64   `protobuf:"varint,8,opt,name=es_timeout,json=esTimeout" json:"es_timeout,omitempty"`
	EsResponse       *int64   `protobuf:"varint,9,opt,name=es_response,json=esResponse" json:"es_response,omitempty"`
	Es               *int64   `protobuf:"varint,10,opt,name=es" json:"es,omitempty"`
}

func (b Metrics) PrintFields() {
	val := reflect.ValueOf(b)
	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)
		fieldName := t.Name

		if jsonTag := t.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
				fieldName = jsonTag[:commaIdx]
			}
		}

		fmt.Println(fieldName)
	}
}


func main() {
	var test Metrics
	test.PrintFields()

}

