package pool

import (
	"fmt"
	"testing"
)

type TestType uint8

func (TestType) Close() error {
	fmt.Println("close")
	return nil
}

func newTestType(x uint8) TestType {
	return TestType(x)
}

func Test_connPool_Get(t *testing.T) {
	pool := NewNoBlockingSimplePool(2)
	pool.Put(newTestType(1))
	fmt.Println("put 1 and get", pool.Get())
	pool.Put(newTestType(2))
	fmt.Println("put 2 and get", pool.Get())
	pool.Put(newTestType(3))
	fmt.Println("put 3")
	pool.Put(newTestType(4))
	fmt.Println("put 4")
	pool.Put(newTestType(5))
	fmt.Println("put 5")
}
