package gtime

import (
	"testing"
	"time"
)

func TestTimeToSeconds(t *testing.T) {
	t.Log(TimeToSeconds(time.Second*2))
}
