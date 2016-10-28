package main

import (
	"log"
	"math/rand"
	"os"
	"time"
	"strconv"
	"github.com/zakfu/metrics"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
)

func main() {
	zero, _ := zmq.NewSocket(zmq.PUSH)
	defer zero.Close()

	endpoint := "tcp://"
	if len(os.Args) == 2 {
		endpoint += os.Args[1]
	} else {
		endpoint += "127.0.0.1:5555"
	}
	zero.Connect(endpoint)
	log.Println("Sending to", endpoint)

	for {
		metric := GetMockMetric()

		data, err := proto.Marshal(metric)
		if err != nil {
			log.Fatalln("Failed to encode metric:", err)
		}

		log.Println("Sending", *metric)
		zero.SendBytes(data, 0)

		time.Sleep(time.Second)
	}
}

func GetMockMetric() *metrics.Metric {
	rand.Seed(time.Now().UTC().UnixNano())
	return &metrics.Metric {
		Measurement: "cpu_temp",
		Tags: []*metrics.Metric_Tag {
			{Key: "host", Value: "host"+strconv.FormatInt(rand.Int63n(3), 10)},
			{Key: "cpu", Value: "cpu"+strconv.FormatInt(rand.Int63n(3), 10)},
		},
		Fields: []*metrics.Metric_Field {
			{Key: "external", Value: rand.Int63n(100)},
			{Key: "internal", Value: rand.Int63n(100)},
		},
		Timestamp: time.Now().UTC().UnixNano(),
	}
}