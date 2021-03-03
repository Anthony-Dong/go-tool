package gstring

import (
	"fmt"
	"strconv"
	"strings"
)

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case uint8, uint16, uint32, uint64:
		convertUint64 := func(value interface{}) uint64 {
			switch v := value.(type) {
			case uint8:
				return uint64(v)
			case uint16:
				return uint64(v)
			case uint32:
				return uint64(v)
			case uint64:
				return v
			default:
				return 0
			}
		}
		return strconv.FormatUint(convertUint64(value), 10)
	case int, int8, int16, int32, int64:
		convertInt64 := func(value interface{}) int64 {
			switch v := value.(type) {
			case int8:
				return int64(v)
			case int16:
				return int64(v)
			case int32:
				return int64(v)
			case int64:
				return v
			case int:
				return int64(v)
			default:
				return 0
			}
		}
		return strconv.FormatInt(convertInt64(value), 10)
	case bool:
		if v {
			return "true"
		}
		return "false"
	// float 效率和直接fmt差距相差不大，为了保证不出问题，还是使用%v
	default:
		str, isOk := value.(fmt.Stringer)
		if isOk {
			return str.String()
		}
		return fmt.Sprintf("%+v", value)
	}
}

func NewString(len int, elem byte) string {
	if len == 0 {
		return ""
	}
	builder := strings.Builder{}
	for x := 0; x < len; x++ {
		builder.WriteByte(elem)
	}
	return builder.String()
}
