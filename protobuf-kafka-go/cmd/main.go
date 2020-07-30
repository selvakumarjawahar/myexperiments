package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	protokafka "github.com/selvakumarjawahar/myexperiments/protobuf-kafka-go/gen"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	//protodynamic "google.golang.org/protobuf/types/dynamicpb"
	"os"
	"os/signal"
	"syscall"
)

func deserializeProtobuf(msg []byte ) {
	metrics := protokafka.Metrics{}
	if err := proto.Unmarshal(msg,&metrics); err != nil {
		fmt.Println("Error in unmarshalling protobuf %s",err)
		return
	}
	fmt.Printf("Timestamp = %d \n", metrics.Timestamp.GetSeconds())
	fmt.Printf("Stream ID = %d \n", metrics.StreamId)
	for key,val := range metrics.Values {
		switch val.GetType().(type) {
		case *protokafka.MetricValue_FloatVal:
			fmt.Printf("%s = %f \n", key, val.GetFloatVal())
		case *protokafka.MetricValue_IntVal:
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

	var msg_builder builder.MessageBuilder
	for key,value := range fields {
		switch value.(type) {
		case int:
			fld_builder := builder.NewField(key,builder.FieldTypeInt64())
			msg_builder.AddField(fld_builder)
		default:
			fld_builder := builder.NewField(key,builder.FieldTypeFloat())
			msg_builder.AddField(fld_builder)
		}
	}
	return msg_builder.Build()
}

func dynamicDeserializationFromDesc(msg []byte, messageBuilder* desc.MessageDescriptor) bool {

	data := dynamic.NewMessage(messageBuilder)
	data.Unmarshal(msg)
	fields := data.GetKnownFields()
	for _,field := range fields {
		fmt.Printf("The Field name = %s Field value = %s",field.String(),data.GetFieldByName(field.String()) )
	}
	return true
}

func dynamicDeserialization(msg []byte) bool {
	msg_bytes := pref.RawFields(msg)
	if !msg_bytes.IsValid() {
		fmt.Printf("Error in the Message, not wireformat")
		return false
	}
	parser := protoparse.Parser{}
	fd,err := parser.ParseFiles("/home/selva/Projects/personal/myexperiments/protobuf-kafka-go/api/proto/callexecuter_metrics.proto")
	if err != nil {
		fmt.Printf("error in parsing protofile %s\n",err)
		return false
	}
	msgfd := fd[0].FindMessage("netrounds.callexecuter.Metrics")
	if msgfd == nil {
		fmt.Printf("Message Not found \n")
		return false
	}
	data := dynamic.NewMessage(msgfd)
	data.Unmarshal(msg)
	fmt.Printf("Timestamp = %d \n", data.GetFieldByName("timestamp"))
	fmt.Printf("Stream ID = %d \n", data.GetFieldByName("stream_id"))
	data.ForEachMapFieldEntryByName("values",func(key interface{}, val interface{}) bool {
			fmt.Printf("%s \n", key)
    		return true
	} )
	return true
}

func main() {

	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <group> <topics..>\n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	group := os.Args[2]
	topics := os.Args[3:]
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

	err = c.SubscribeTopics(topics, nil)

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
					//deserializeProtobuf(e.Value)
					dynamicDeserialization(e.Value)
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
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

