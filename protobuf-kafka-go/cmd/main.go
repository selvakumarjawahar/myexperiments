package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	exec "github.com/selvakumarjawahar/myexperiments/protobuf-kafka-go/gen/callexecuter"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"errors"
	//protodynamic "google.golang.org/protobuf/types/dynamicpb"
	"os"
	"os/signal"
	"syscall"
)

func deserializeCallexecuter(msg []byte ) {
	metrics := exec.Metrics{}
	if err := proto.Unmarshal(msg,&metrics); err != nil {
		fmt.Println("Error in unmarshalling protobuf %s",err)
		return
	}
	fmt.Printf("Timestamp = %d \n", metrics.Timestamp.GetSeconds())
	fmt.Printf("Stream ID = %d \n", metrics.StreamId)
	for key,val := range metrics.Values {
		switch val.GetType().(type) {
		case *exec.MetricValue_FloatVal:
			fmt.Printf("%s = %f \n", key, val.GetFloatVal())
		case *exec.MetricValue_IntVal:
			fmt.Printf("%s = %d \n", key, val.GetIntVal())
		}
	}

}

//message {
//metrics
//}
//connect_time_avg:float
//first_byte_time_avg:float
//response_time_min:float
//response_time_avg:float
//response_time_max : float
//size_avg : float
//speed_avg : float
//es_timeout : int
//es_response : int
//es : int

type MetricFields map[string]interface{}

func buildProto(fields MetricFields) (*desc.MessageDescriptor,error) {

	msg_builder := builder.NewMessage("Metrics")

	for key,value := range fields {
		switch value.(type) {
		case int:
			msg_builder.AddField(builder.NewField(key,builder.FieldTypeInt64()))
		default:
			msg_builder.AddField(builder.NewField(key,builder.FieldTypeFloat()))
		}
	}
	return msg_builder.Build()
}

func dynamicDeserializationFromDesc(msg []byte, messageBuilder* desc.MessageDescriptor) bool {

	data := dynamic.NewMessage(messageBuilder)
	data.Unmarshal(msg)
	fields := data.GetKnownFields()
	for _,field := range fields {
		fmt.Printf("The Field name = %s Field value = %v \n",
			field.GetName(),
			data.GetFieldByName(field.GetName()))
	}
	return true
}

func dynamicDeserialization(msg []byte, protofile string, msgName string) (*dynamic.Message, error) {
	msg_bytes := pref.RawFields(msg)
	if !msg_bytes.IsValid() {
		fmt.Printf("Error in the Message, not wireformat")
		return nil, errors.New("Error in the Message, not wireformat")
	}
	parser := protoparse.Parser{}
	fd,err := parser.ParseFiles(protofile)
	if err != nil {
		fmt.Printf("error in parsing protofile %s\n",err)
		return nil, errors.New("error in parsing protofile")
	}
	msgfd := fd[0].FindMessage(msgName)
	if msgfd == nil {
		fmt.Printf("Message Not found \n")
		return nil, errors.New("Message not found")
	}
	data := dynamic.NewMessage(msgfd)

	data.Unmarshal(msg)

	return data, nil
}

func parseHTTPMetrics(message *dynamic.Message) {

	var intVal int
	var floatVal float64
	httpFields :=  MetricFields{
		"connect_time_avg" : floatVal,
		"first_byte_time_avg" : floatVal,
		"response_time_min" : floatVal,
		"response_time_avg" : floatVal,
		"response_time_max" : floatVal,
		"size_avg" : floatVal,
		"speed_avg" : floatVal,
		"es_timeout" : intVal,
		"es_response" : intVal,
		"es" : intVal }
	fmt.Print("running build proto \n")

	md, err := buildProto(httpFields)
	if  err != nil {
		fmt.Printf("Error in Generating message descriptior")
		return
	}
	fmt.Print("build proto successful")

	if message.HasFieldName("metrics") {
		values := message.GetFieldByName("metrics").(*dynamic.Message).GetFieldByName("values")
		dynamicDeserializationFromDesc(values.([]byte),md)
	} else {
		fmt.Print("Field not found\n")
	}
}

func main() {

	if len(os.Args) < 6 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <group> <topic> <protofile> <MessageName> \n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	group := os.Args[2]
	topic := os.Args[3]
	protofile := os.Args[4]
	messgName := os.Args[5]

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"broker.address.family": "v4",
		"group.id":              group,
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest"})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.Subscribe(topic, nil)

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				data,err := dynamicDeserialization(e.Value, protofile, messgName)
				if err != nil {
					fmt.Printf("Error in serialization %s\n",err.Error())
				} else {
					parseHTTPMetrics(data)
				}
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
	fmt.Printf("Closing consumer\n")
	c.Close()
}

