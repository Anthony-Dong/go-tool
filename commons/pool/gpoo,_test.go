package pool

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGoWithRecover(t *testing.T) {
	GoWithRecover(func() {
		panic("hello")
	}, nil)

	time.Sleep(time.Second)
}

func TestNew(t *testing.T) {
	pool := New(5, 0)

	fmt.Println(os.Getpid())

	for {
		if err := pool.AddJob(func() {

		}); err != nil {
			panic(err)
		}
	}

}

func TestNewCond(t *testing.T) {
	tests := []struct {
		name string
		want Condition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCond(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCond() = %v, want %v", got, tt.want)
			}
		})
	}
}
