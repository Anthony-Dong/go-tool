package pool

import (
	"io"
)

type ConnPool interface {
	Get() io.Closer
	Put(v io.Closer)
	Release()
}

type Option func(pool connPool)

type connPool struct {
	MaxOpenNum int
	buffer     chan io.Closer
	New        func() io.Closer
	down       chan struct{}
}

func NewNoBlockingSimplePool(max int) ConnPool {
	return NewNoBlockingPool(max, func() io.Closer {
		return nil
	})
}

func NewNoBlockingPool(max int, new func() io.Closer) ConnPool {
	return &connPool{
		MaxOpenNum: max,
		buffer:     make(chan io.Closer, max),
		New:        new,
		down:       make(chan struct{}, 0), // 关闭的时候直接close
	}
}

func (c *connPool) Get() io.Closer {
	select {
	case <-c.down: // 如果关闭，获取的始终为空
		return nil
	case data, isOpen := <-c.buffer: // 如果已经关闭，获取的为空，否则有buffer
		if isOpen {
			return data
		}
		return nil
	default:
		return c.New() // 默认是不阻塞
	}
}

func (c *connPool) Put(v io.Closer) {
	select {
	case <-c.down:
		return
	case c.buffer <- v:

	default:

	}
}

func (c *connPool) Release() {
	defer close(c.buffer)
	close(c.down)
	closeResource := func() error {
		for {
			select {
			case data := <-c.buffer:
				if err := data.Close(); err != nil {
					return err
				}
			default:
				return nil
			}
		}
	}
	if err := closeResource(); err != nil {

	}
}
