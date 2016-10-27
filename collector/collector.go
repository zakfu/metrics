package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"github.com/zakfu/metrics"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
)

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

	ch := make(chan metrics.Metric)
	go OutputMetrics(ch)

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
		} else {
			ch <- *metric
		}
	}
}

func OutputMetrics(ch chan metrics.Metric) {
	conn, err := net.Dial("udp", "localhost:8089")
	if err != nil {
		log.Fatalln("UDP error")
	}
	defer conn.Close()

	for m := range ch {
		im := metrics.InfluxMetric{&m}
		fmt.Println(im)
		conn.Write([]byte(im.String()))
	}
}
