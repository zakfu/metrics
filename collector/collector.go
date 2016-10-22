package main

import (
	"log"
	"os"
	"time"
	"github.com/zakfu/metrics"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
)

func output_metrics(c chan metrics.Metric) {
	for m := range c {
		log.Println("Output metric:", m)
		// Simulate write time
		time.Sleep(time.Millisecond * 1000)
	}
}

func main() {
	zero, _ := zmq.NewSocket(zmq.PULL)
	defer zero.Close()

	endpoint := "tcp://*:"
	if len(os.Args) == 2 {
		endpoint += os.Args[1]
	} else {
		endpoint += "5555"
	}
	zero.Bind(endpoint)
	log.Println("Listening on", endpoint)

	output_channel := make(chan metrics.Metric)
	go output_metrics(output_channel)

	for {
		data, err := zero.RecvBytes(0)
		if err != nil {
			log.Println("Failed to read from socket:", err)
			continue
		}
		metric := &metrics.Metric{}
		if err := proto.Unmarshal(data, metric); err != nil {
			log.Println("Failed to parse metric:", err)
			continue
		}	else {
			output_channel <- *metric
		}
	}
}
