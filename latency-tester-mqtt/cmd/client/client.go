package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/richiMarchi/latency-tester/latency-tester-mqtt/pkg/logic"
	"k8s.io/klog/v2"
)

const id = "latency-tester-client"

func main() {
	broker := flag.String("broker", "", "The address to contact the broker")
	username := flag.String("username", "", "The broker username")
	password := flag.String("password", "", "The broker password")
	messageRate := flag.Uint("messageRate", 100, "number of messages per second")
	experimentDuration := flag.Uint("duration", 120,
									"Duration in s of th experiment (keep 0 if unlimited)")
	requestSize := flag.Uint("requestSize", 1024, "bytes of the payload")
	qos := flag.Uint("qos", 0, "mqtt QoS")
	klog.InitFlags(nil)
	flag.Parse()

	klog.Infof("Broker: %v", *broker)
	klog.Infof("Initial mps: %v", *messageRate)
	klog.Infof("Request Size: %v Bytes", *requestSize)
	klog.Infof("QoS: %v", byte(*qos))

	logic.ConfigureLogging()

	opts := logic.BuildCommonConnectionOptions(*broker, id, *username, *password)
	client, err := logic.EstablishBrokerConnection(opts)
	if err != nil {
		klog.Fatal(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGTERM)

	requester := logic.NewClientRequester(client, *messageRate, *experimentDuration, *requestSize, byte(*qos), shutdown)
	requester.PublishRequests()

	klog.Info("Exiting")
	client.Disconnect(logic.DisconnectQuiescence)
}
