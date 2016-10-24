package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
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
	for m := range ch {
		im := InfluxMetric{&m}
		fmt.Println(im)
		// Simulate write time
		time.Sleep(time.Millisecond * 1000)
	}
}

type InfluxMetric struct {
	Metric *metrics.Metric
}

func (im InfluxMetric) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(im.Metric.Measurement)
	for _, t := range im.Metric.Tags {
		buffer.WriteString(",")
		buffer.WriteString(t.Key)
		buffer.WriteString("=")
		buffer.WriteString(t.Value)
	}
	buffer.WriteString(" ")
	for i, f := range im.Metric.Fields {
		buffer.WriteString(f.Key)
		buffer.WriteString("=")
		buffer.WriteString(strconv.FormatInt(f.Value, 10))
		if i < len(im.Metric.Fields)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(" ")
	buffer.WriteString(strconv.FormatInt(im.Metric.Timestamp, 10))
	return buffer.String()
}