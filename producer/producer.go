package main

import (
	"log"
	"math/rand"
	"os"
	"time"
	"github.com/zakfu/metrics"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	
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

	for request_nbr := 0; request_nbr != 10; request_nbr++ {
		metric := &metrics.Metric {
			Measurement: "mock_metric",
			Value: rand.Int63n(1000),
		}

		data, err := proto.Marshal(metric)
		if err != nil {
			log.Fatalln("Failed to encode metric:", err)
		} else {
			log.Println("Sending", *metric)
			zero.SendBytes(data, 0)
		}
		
		time.Sleep(time.Millisecond * 500)
	}
}
