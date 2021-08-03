package logic

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	serialization "github.com/richiMarchi/latency-tester/latency-tester-mqtt/pkg/message/serialization/protobuf"

	"google.golang.org/protobuf/proto"
	"k8s.io/klog/v2"
)

type ServerSubscriber struct {
	client  mqtt.Client
	qos     byte
	output  *os.File
}

func NewServerSubscriber(client mqtt.Client, qos byte, outputFile string) *ServerSubscriber {
	file, err := os.Create(outputFile)
	if err != nil {
		klog.Fatal("Failed to open output file: ", err)
	}
	_, _ = file.WriteString("client-send-timestamp,msg-latency-ms\n")

	return &ServerSubscriber{
		client:  client,
		qos:     qos,
		output:  file,
	}
}

func (s *ServerSubscriber) Subscribe() {
	klog.Infof("Subscribing to topic: %s", RequestTopic)
	token := s.client.Subscribe(RequestTopic, s.qos, s.onMessage)
	token.Wait()
	klog.Infof("Subscribed to topic: %s", RequestTopic)
}

func (s *ServerSubscriber) onMessage(client mqtt.Client, msg mqtt.Message) {
	request := &serialization.Message{}
	err := proto.Unmarshal(msg.Payload(), request)
	if err != nil {
		klog.Errorf("Failed to unmarshal message: %v", err)
	}

	// Save the response
	latency := time.Since(request.ClientTimestamp.AsTime())
	latencyMs := float64(latency.Nanoseconds()) / float64(time.Millisecond.Nanoseconds())

	klog.Infof("Received message %d in %.2f ms", request.Id, latencyMs)

	outputStr := fmt.Sprintf("%v,%f\n",
		request.ClientTimestamp.AsTime().UnixNano(),
		latencyMs)
	_, _ = s.output.WriteString(outputStr)
}