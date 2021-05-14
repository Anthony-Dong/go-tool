package wrk

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/anthony-dong/go-tool/commons/codec/gjson"

	"github.com/anthony-dong/go-tool/command"
	"github.com/anthony-dong/go-tool/command/api"
	"github.com/anthony-dong/go-tool/commons/collections"
	"github.com/anthony-dong/go-tool/commons/ghttp"
	"github.com/juju/errors"
	"github.com/panjf2000/ants/v2"
	"github.com/urfave/cli/v2"
)

type Config struct {
	Duration    time.Duration   `json:"duration"` // 处理时间
	Url         string          `json:"url"`      // 请求路径
	HeaderSlice cli.StringSlice `json:"-"`
	Header      http.Header     `json:"header"` // 请求头
	Body        string          `json:"body"`   // 请求体

	Threads     int `json:"threads"`     // 最大并发数量
	Connections int `json:"connections"` // 最大连接数

	Timeout time.Duration `json:"timeout"` // 连接/请求超时时间
	Method  string        `json:"method"`  // 请求方法
}

func NewConfig(time time.Duration, url string, op ...Option) *Config {
	config := new(Config)
	config.Duration = time
	config.Url = url
	initOp(config, op...)
	return config
}

func NewCommand() command.Command {
	return new(Config)
}

func (c *Config) InitConfig(context *cli.Context, config api.CommonConfig) ([]byte, error) {
	headers, err := ghttp.NewHeader(c.HeaderSlice.Value())
	if err != nil {
		return nil, err
	}
	c.Header = headers
	c.Method = strings.ToUpper(c.Method)
	return gjson.ToJsonString(c), nil
}

func (c *Config) Flag() []cli.Flag {
	return []cli.Flag{
		&cli.DurationFlag{
			Name:        "duration",
			Aliases:     []string{"d"},
			Usage:       fmt.Sprintf("Set the request duration for wrk (%s)", collections.ToCliMultiDescString([]string{"1ms", "1s", "1m", "1h"})),
			Destination: &c.Duration,
			Required:    true,
		},
		&cli.IntFlag{
			Name:        "connections",
			Aliases:     []string{"c"},
			Usage:       fmt.Sprintf("Connections to keep open"),
			Destination: &c.Connections,
			Required:    false,
			Value:       defaultConfig.Connections,
		},
		&cli.IntFlag{
			Name:        "threads",
			Aliases:     []string{"t"},
			Usage:       fmt.Sprintf("Number of threads to use"),
			Destination: &c.Threads,
			Required:    false,
			Value:       defaultConfig.Threads,
		},
		&cli.StringFlag{
			Name:        "method",
			Aliases:     []string{"m"},
			Usage:       fmt.Sprintf("Set the http request method method"),
			Destination: &c.Method,
			Required:    false,
			Value:       "GET",
		},
		&cli.StringFlag{
			Name:        "url",
			Aliases:     []string{"u"},
			Usage:       "Set the request url",
			Destination: &c.Url,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "body",
			Aliases:     []string{"b"},
			Usage:       fmt.Sprintf("Set the request body"),
			Destination: &c.Body,
			Required:    false,
		},
		&cli.StringSliceFlag{
			Name:        "header",
			Aliases:     []string{"H"},
			Usage:       fmt.Sprintf("Set the request header"),
			Destination: &c.HeaderSlice,
			Required:    false,
		},
		&cli.DurationFlag{
			Name:        "timeout",
			Usage:       fmt.Sprintf("Socket request timeout"),
			Destination: &c.Timeout,
			Required:    false,
			Value:       defaultConfig.Timeout,
		},
	}
}

func (c *Config) Run(context *cli.Context) error {
	start := time.Now()
	if err := run(c); err != nil {
		return err
	}
	// 打印信息
	writeHistogram(requestDurations, int(c.Duration/time.Second))
	writeCounter(requestBodyCount, int(c.Duration/time.Second))
	fmt.Printf("程序结束，一共花费 %.3fs\n", time.Now().Sub(start).Seconds())
	return nil
}

var (
	httpClientPool = sync.Pool{New: func() interface{} {
		return &http.Client{}
	}}
	bufferPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	defaultTransport http.RoundTripper
	errChan          = make(chan error, 0)
)

func run(config *Config) error {
	defaultTransport = newTransport(config) // 创建连接池
	defer func() {
	}()
	// 线程池
	gpool, err := ants.NewPool(int(config.Threads), func(opts *ants.Options) {
		opts.MaxBlockingTasks = 1 << 20 // 最大 1024*1024个任务
		opts.PreAlloc = true
	}) // 创建线程池
	if err != nil {
		return errors.Trace(err)
	}
	defer gpool.Release() // 最后释放

	// 创建timer
	ctx, cancelFunc := context.WithTimeout(context.Background(), config.Duration) // 创建超时时间
	defer cancelFunc()

	// 调度
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				sendError(gpool.Submit(func() {
					if err := request(config); err != nil {
						sendError(err)
					}
				}))
			}
		}
	}()

	// wait
	select {
	case <-ctx.Done():

	case err := <-errChan:
		return err
	}
	return nil
}

func newTransport(config *Config) http.RoundTripper {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   config.Timeout,
			KeepAlive: config.Duration,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxConnsPerHost:       config.Connections,
		MaxIdleConnsPerHost:   config.Connections,
		MaxIdleConns:          config.Connections,
		IdleConnTimeout:       config.Duration,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func request(config *Config) (err error) {
	// 创建客户端
	client := resetClient()
	defer httpClientPool.Put(client)
	client.Timeout = config.Timeout
	client.Transport = defaultTransport

	// 获取请求体
	buffer, err := newRequestReader(config.Body)
	if err != nil {
		return err
	}
	defer recycleBuffer(buffer)

	// 获取请求
	request, err := http.NewRequestWithContext(context.Background(), config.Method, config.Url, ioutil.NopCloser(buffer))
	if err != nil {
		return err
	}
	// 添加header
	for headerKey, headerValue := range config.Header {
		for _, value := range headerValue {
			request.Header.Set(headerKey, value)
		}
	}
	start := time.Now()
	defer func() {
		requestDurations.Observe(float64(time.Now().Sub(start).Milliseconds()))
	}()

	// 发送请求
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	// 关闭连接
	defer response.Body.Close()
	if err := readResponseBuffer(response, func(responseBody []byte) {
		requestBodyCount.Add(float64(len(responseBody)))
	}); err != nil { // 如果读取不完就复用不了连接！
		return err
	}
	return nil
}

func resetClient() (client *http.Client) {
	client = httpClientPool.Get().(*http.Client)
	client.CheckRedirect = nil
	client.Jar = nil
	client.Transport = nil
	return client
}

func newRequestReader(body string) (*bytes.Buffer, error) {
	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	if _, err := buffer.WriteString(body); err != nil {
		return nil, err
	}
	return buffer, nil
}

func recycleBuffer(buf *bytes.Buffer) {
	if buf != nil {
		bufferPool.Put(buf)
	}
}

func readResponseBuffer(resp *http.Response, foo func(responseBody []byte)) error {
	if resp.Body == nil {
		foo([]byte{})
		return nil
	}
	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	if _, err := buffer.ReadFrom(resp.Body); err != nil {
		return err
	}
	if foo != nil {
		foo(buffer.Bytes())
	}
	return nil
}

func sendError(err error) {
	if err != nil {
		errChan <- err
	}
}

//// 统计流量
