package metrics

import (
	"bytes"
	"fmt"
)

type InfluxMetric struct {
	*Metric
}

func (im InfluxMetric) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(im.Measurement)
	for _, t := range im.Tags {
		buffer.WriteString(fmt.Sprintf(",%s=%s", t.Key, t.Value))
	}
	buffer.WriteString(" ")
	for i, f := range im.Fields {
		buffer.WriteString(fmt.Sprintf("%s=%d", f.Key, f.Value))
		if i < len(im.Fields)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(fmt.Sprintf(" %d", im.Timestamp))
	return buffer.String()
}