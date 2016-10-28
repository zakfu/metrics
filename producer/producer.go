package main

import (
	"flag"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	"github.com/zakfu/metrics"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	dest := flag.String("dest", "127.0.0.1:5555", "host:port to send metrics to")
	interval := flag.Int("interval", 1000, "interval (in milliseconds) at which to send metrics")
	dump := flag.Bool("dump", false, "show metrics on STDOUT")
	flag.Parse()

	zero, _ := zmq.NewSocket(zmq.PUSH)
	defer zero.Close()

	endpoint := "tcp://" + *dest
	zero.Connect(endpoint)
	log.Println("Sending to", endpoint)

	for {
		metric := GetMockMetric()

		data, err := proto.Marshal(metric)
		if err != nil {
			log.Fatalln("Failed to encode metric:", err)
		}

		if *dump {
			log.Println("Sending", *metric)
		}

		zero.SendBytes(data, 0)

		time.Sleep(time.Millisecond * time.Duration(*interval))
	}
}

func GetMockMetric() *metrics.Metric {
	rand.Seed(time.Now().UTC().UnixNano())
	return &metrics.Metric{
		Measurement: "cpu_temp",
		Tags: []*metrics.Metric_Tag{
			{Key: "host", Value: "host" + strconv.FormatInt(rand.Int63n(3), 10)},
			{Key: "cpu", Value: "cpu" + strconv.FormatInt(rand.Int63n(3), 10)},
		},
		Fields: []*metrics.Metric_Field{
			{Key: "external", Value: rand.Int63n(100)},
			{Key: "internal", Value: rand.Int63n(100)},
		},
		Timestamp: time.Now().UTC().UnixNano(),
	}
}
