package logic

import (
	"crypto/rand"
	"os"
	"time"

	serialization "github.com/richiMarchi/latency-tester/latency-tester-mqtt/pkg/message/serialization/protobuf"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/klog/v2"
)

type ClientRequester struct {
	client           mqtt.Client
	currentMessageID   uint
	messageRate        uint
	experimentDuration time.Duration
	qos                byte
	payload            []byte

	shutdown chan os.Signal
}

func NewClientRequester(client mqtt.Client, messageRate uint, experimentDuration uint, requestSize uint,
						qos byte, shutdown chan os.Signal) *ClientRequester {
	payload := make([]byte, requestSize)
	if _, err := rand.Read(payload); err != nil {
		klog.Fatal("Failed to build payload", err)
	}

	return &ClientRequester{
		client:             client,
		currentMessageID:   0,
		messageRate: 	    messageRate,
		experimentDuration: time.Duration(experimentDuration) * time.Second,
		payload:            payload,
		qos:                qos,
		shutdown:           shutdown,
	}
}

func (r *ClientRequester) WorkOutRate(currentRate uint) time.Duration{
	return time.Duration(uint(1 / float32(currentRate) * 1000000)) * time.Microsecond
}

func (r *ClientRequester) PublishRequests() {
	klog.Infof("Starting to publish requests on %s", RequestTopic)
	currentRate := r.messageRate
	intervalNs := r.WorkOutRate(currentRate)
	experimentStart := time.Now()
	klog.Infof("Start publishing with rate %v", intervalNs)
	klog.Infof("Experiment will last %v", r.experimentDuration)
	for r.experimentDuration == 0 || time.Since(experimentStart) <  r.experimentDuration{
		start := time.Now()
		r.publishRequest(&start)

		waitTime := intervalNs - time.Since(start)
		if waitTime < 0 {
			klog.Warningf("Missed deadline when sending message %v %v interval %v since start", r.currentMessageID, intervalNs, time.Since(start))
			waitTime = 0
		}

		r.currentMessageID++

		select {
		case signal := <-r.shutdown:
			klog.Info("Signal caught - exiting")
			r.shutdown <- signal
			return
		case <-time.After(waitTime):
		}
	}
	klog.Infof("Finished to publish requests on %s", RequestTopic)
}

func (r *ClientRequester) publishRequest(now *time.Time) {
	message := &serialization.Message{
		Id:              int32(r.currentMessageID),
		ClientTimestamp: timestamppb.New(*now),
		ServerTimestamp: &timestamppb.Timestamp{},
		Payload:         r.payload,
	}

	marshal, err := proto.Marshal(message)
	if err != nil {
		klog.Errorf("Failed to marshal message: %v", err)
	}
	klog.V(3).Infof("Publishing message %v", r.currentMessageID)
	token := r.client.Publish(RequestTopic, r.qos, false, marshal)
	go func(id uint) {
		<-token.Done()
		klog.V(3).Infof("Confirmed publication of message %v in %v", id, time.Since(*now))
		if token.Error() != nil {
			klog.Error("Failed to publish message %v: ", id, token.Error())
		}
	}(r.currentMessageID)
}
