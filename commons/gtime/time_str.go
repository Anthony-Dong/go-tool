package gtime

import (
	"strconv"
	"time"
)

// 时间之差 s 输出 0.100010s.
func TimeToSeconds(duration time.Duration) string {
	// 1s=1000ms 1ms=1000us  保留6位到us
	return strconv.FormatInt(int64(duration/time.Second), 10) + "s"
}

// 除固定值，保留固定小数位.
func Float642String(num float64, saveDecimalPoint int) string {
	return strconv.FormatFloat(num, 'f', saveDecimalPoint, 64)
}
