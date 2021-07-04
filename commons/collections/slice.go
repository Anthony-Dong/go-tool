package collections

import (
	"reflect"
	"strings"

	"github.com/anthony-dong/go-tool/commons/gstring"

	"github.com/juju/errors"
)

// 将 map 的key 拿出来 转换成 string.
func GetMapKeysToString(obj interface{}) (list []string, err error) {
	if obj == nil {
		return nil, errors.New("the obj is nil can not get keys")
	}
	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
	}
	if value.Kind() != reflect.Map {
		return nil, errors.New("the obj kind is not map can not get keys")
	}
	values := value.MapKeys()
	result := make([]string, 0, len(values))
	for _, elem := range values {
		result = append(result, gstring.ToString(elem.Interface()))
	}
	return result, nil
}

// 转换成 cli命令的 多个条件描述文本，例如 k1,k2 => "k1"|"k2".
func ToCliMultiDescString(slice []string) string {
	if slice == nil || len(slice) == 0 {
		return ""
	}
	lastIndex := len(slice) - 1
	result := strings.Builder{}
	for index, elem := range slice {
		result.WriteByte('"')
		result.WriteString(elem)
		result.WriteByte('"')
		if index != lastIndex {
			result.WriteByte('|')
		}
	}
	return result.String()
}

func SliceLineToString(lines []string) string {
	if lines == nil || len(lines) == 0 {
		return ""
	}
	builder := strings.Builder{}
	for _, elem := range lines {
		builder.WriteString(elem)
		builder.WriteByte('\n')
	}
	return builder.String()
}

func SplitStringSlice(slice []string, length int) [][]string {
	if slice == nil || len(slice) == 0 {
		return [][]string{}
	}
	if len(slice) <= length {
		return [][]string{slice}
	}
	cut := 0
	if len(slice)%length == 0 {
		cut = len(slice) / length
	} else {
		cut = len(slice)/length + 1
	}
	result := make([][]string, 0, cut)
	for x := 0; x < cut; x++ {
		end := x*length + length
		if end > len(slice) {
			end = len(slice)
		}
		result = append(result, slice[x*length:end])
	}
	return result
}
