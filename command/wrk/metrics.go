package wrk

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

var (
	requestCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "请求次数",
		Help:        "请求接口的次数",
		ConstLabels: prometheus.Labels{},
	})
	requestBodyCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "请求体大小",
		Help:        "请求体大小.",
		ConstLabels: prometheus.Labels{},
	})
	requestDurations = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "请求时间统计",
			Help:    "请求时间统计。",
			Buckets: []float64{1, 5, 10, 50, 100, 500, 1000, 5000}, // 统计1ms,5ms,10ms,100ms
		},
	)
)

func writeHistogram(resource prometheus.Metric, time int) {
	metric := new(dto.Metric)
	if err := resource.Write(metric); err != nil {
		return
	}
	fmt.Println("==========请求时长分布图==============")
	for _, bucket := range metric.Histogram.Bucket {
		fmt.Printf("%.0fms  %d  %s\n", *bucket.UpperBound, *bucket.CumulativeCount, ca1(*bucket.CumulativeCount, *metric.Histogram.SampleCount))
		if bucket.Exemplar != nil {
			fmt.Printf("Exemplar: %v\n", *bucket.Exemplar)
		}
	}
	fmt.Println("==========请求次数分布图==============")
	fmt.Printf("请求总次数: %d\n", *metric.Histogram.SampleCount)
	fmt.Printf("请求次数/s: %d\n", (*metric.Histogram.SampleCount)/uint64(time))
}

func writeCounter(resource prometheus.Metric, time int) {
	metric := new(dto.Metric)
	if err := resource.Write(metric); err != nil {
		return
	}
	fmt.Println("==========请求体吞吐量==============")
	fmt.Printf("请求吞吐量(kb): %dkb\n", uint64(*metric.Counter.Value)>>10)
	fmt.Printf("请求每秒吞吐量(kb/s): %dkb\n", (uint64(*metric.Counter.Value)/uint64(time))>>10)
}

func ca1(num1, num2 uint64) string {
	result := (float64(num1) / float64(num2)) * 100
	return fmt.Sprintf("%.2f%s", result, "%")
}
