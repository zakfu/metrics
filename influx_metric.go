package metrics

import (
  "bytes"
	"strconv"
)

type InfluxMetric struct {
	Metric *Metric
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