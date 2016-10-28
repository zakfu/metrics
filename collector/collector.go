package main

import (
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	"github.com/zakfu/metrics"
	"log"
	"net"
)

func main() {
	bind := flag.String("bind", "*:5555", "host:port to listen on")
	dump := flag.Bool("dump", false, "show metrics on STDOUT")
	flag.Parse()

	zero, _ := zmq.NewSocket(zmq.PULL)
	defer zero.Close()

	endpoint := "tcp://" + *bind
	zero.Bind(endpoint)
	log.Println("Listening on", endpoint)

	ch := make(chan metrics.Metric)
	go OutputMetrics(ch, *dump)

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

func OutputMetrics(ch chan metrics.Metric, dump bool) {
	conn, err := net.Dial("udp", "localhost:8089")
	if err != nil {
		log.Fatalln("UDP error")
	}
	defer conn.Close()

	for m := range ch {
		im := metrics.InfluxMetric{&m}
		if dump {
			fmt.Println(im)
		}
		conn.Write(im.Bytes())
	}
}
