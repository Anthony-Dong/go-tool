package pool

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	FuncIsNilError     = errors.New("func is nil")
	GPollIsClosedError = errors.New("go poll is closed")
	ArgsError          = errors.New("args has error")
)

type GPool interface {
	AddJob(run Fun) error
	Stop()
}

type Fun func()

type gPool struct {
	runningCount int64
	maxG         int64
	taskChannel  chan Fun
	cond         Condition
	down         chan struct{}
}

func New(maxG, maxWaitJobNum int) (GPool, error) {
	if maxG <= 0 {
		return nil, ArgsError
	}
	if maxWaitJobNum < 0 {
		return nil, ArgsError
	}
	pool := new(gPool)
	pool.maxG = int64(maxG)
	pool.cond = NewCond()
	pool.taskChannel = make(chan Fun, maxG+maxWaitJobNum) // 最大等待任务数量是除运行中的等待的任务数量
	pool.down = make(chan struct{}, 0)
	pool.run()
	return pool, nil
}

func (g *gPool) AddJob(run Fun) error {
	if run == nil {
		return FuncIsNilError
	}
	select {
	case <-g.down:
		return GPollIsClosedError
	case g.taskChannel <- run:
		return nil
	}
}

func (g *gPool) Stop() {
	close(g.down)
	close(g.taskChannel)
	g.cond.NotifyAll() // 释放run方法的主线程
}

func (g *gPool) run() {
	go func() {
		for {
			select {
			case job := <-g.taskChannel:
				if atomic.AddInt64(&g.runningCount, 1) > g.maxG {
					g.cond.Wait()
				}
				select {
				case <-g.down:
					return
				default:
					GoWithRecover(func() {
						defer func() {
							atomic.AddInt64(&g.runningCount, -1)
							g.cond.Notify()
						}()
						job()
					}, nil)
				}
			case <-g.down:
				return
			}
		}
	}()
}

type Condition interface {
	Wait()
	Notify()
	NotifyAll()
}

type condition struct {
	sync.Mutex
	c *sync.Cond
}

func NewCond() Condition {
	c := new(condition)
	c.c = sync.NewCond(c)
	return c
}
func (c *condition) Wait() {
	c.Lock()
	defer c.Unlock()
	c.c.Wait()
}
func (c *condition) Notify() {
	c.c.Signal()
}
func (c *condition) NotifyAll() {
	c.c.Broadcast()
}

func GoWithRecover(handler func(), recoverHandler func(r interface{})) {
	if handler == nil {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if recoverHandler == nil {
					PanicPrint(r, os.Stderr)
					return
				}
				go func() {
					defer func() {
						if r := recover(); r != nil {
							PanicPrint(r, os.Stderr)
						}
					}()
					recoverHandler(r)
				}()
			}
		}()
		handler()
	}()
}

func PanicPrint(recoverError interface{}, writer io.Writer) {
	stackBuffer := make([]byte, 64<<10) // 最多打印64k的堆栈信息
	if _, err := fmt.Fprintf(writer, "panic: %v\n\n%s\n", recoverError, stackBuffer[:runtime.Stack(stackBuffer, false)]); err != nil {
		fmt.Printf("panic: %v\n\n%s\n", recoverError, stackBuffer[:runtime.Stack(stackBuffer, false)]) // 出现异常用 std 输出
	}
}
