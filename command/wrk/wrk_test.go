package wrk

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestName(t *testing.T) {

	// 统计流量
	Counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "请求次数",
		Help: "请求接口的次数",
		ConstLabels: prometheus.Labels{
			"path":   "/",
			"method": "get",
		},
	})

	for x := 0; x < 1000; x++ {
		Counter.Inc()
	}

	metric := dto.Metric{}
	err := Counter.Write(&metric)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(metric)

}

func TestSummary(t *testing.T) {
	cool := make(chan prometheus.Metric)
	Summary := prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "请求次数",
		Help: "请求接口的次数",
		ConstLabels: prometheus.Labels{
			"path":   "/",
			"method": "get",
		},
	})
	go func() {
		for x := 0; x < 100000; x++ {
			time.Sleep(time.Millisecond * 10)
			Summary.Observe(float64(x))
		}
	}()
	go func() {
		for {
			<-time.After(time.Second)
			Summary.Collect(cool)
		}
	}()

	for {
		select {
		case x := <-cool:
			metric := dto.Metric{}
			if err := x.Write(&metric); err != nil {
				t.Fatal(err)
			}
			fmt.Println(metric)
		}
	}
}

func TestTimeOutAndKeepAlive(t *testing.T) {
	dialer := net.Dialer{
		Timeout:   time.Second,
		KeepAlive: time.Second,
	}
	if _, err := dialer.DialContext(context.Background(), "tcp", ":8888"); err != nil {
		t.Fatal(err)
	}
}

func TestRun(t *testing.T) {
	config := NewConfig(time.Second*10, "http://localhost:8888", func(config *Config) {
	})
	if err := run(config); err != nil {
		return
	}
}

func TestMaxCPU(t *testing.T) {
	t.Log(runtime.NumCPU())
}

func TestRequest(t *testing.T) {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxConnsPerHost:       10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Millisecond,
	}
	wg := sync.WaitGroup{}
	for x := 0; x < 1000; x++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := http.Client{
				Transport: transport,
			}
			newRequest, err := http.NewRequest(http.MethodGet, "http://localhost:8888", nil)
			if err != nil {
				t.Fatal(err)
			}
			response, err := client.Do(newRequest)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			bytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("body: %s\n", bytes)
		}()
	}
	wg.Wait()
}
