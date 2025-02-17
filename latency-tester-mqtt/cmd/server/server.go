package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/richiMarchi/latency-tester/latency-tester-mqtt/pkg/logic"
	"k8s.io/klog/v2"
)

const id = "latency-tester-server"

func main() {
	broker := flag.String("broker", "", "The address to contact the broker")
	username := flag.String("username", "", "The broker username")
	password := flag.String("password", "", "The broker password")
	qos := flag.Uint("qos", 0, "mqtt QoS")
	log := flag.String("log", "./log.csv", "file to store latency results")
	klog.InitFlags(nil)
	flag.Parse()

	// Generate client id
	rand.Seed(time.Now().UnixNano())
	clientId := fmt.Sprintf("%v-%d", id, rand.Int())

	klog.Infof("Client id: %v", clientId)
	klog.Infof("Broker: %v", *broker)
	klog.Infof("QoS: %v", byte(*qos))

	logic.ConfigureLogging()

	opts := logic.BuildCommonConnectionOptions(*broker, clientId, *username, *password)
	client, err := logic.EstablishBrokerConnection(opts)
	if err != nil {
		klog.Fatal(err)
	}

	subscriber := logic.NewServerSubscriber(client, byte(*qos), *log)
	subscriber.Subscribe()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGTERM)

	<-shutdown
	klog.Info("Signal caught - exiting")
	client.Disconnect(logic.DisconnectQuiescence)
}
