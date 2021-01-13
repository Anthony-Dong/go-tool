package util

import (
	"encoding/json"

	"github.com/anthony-dong/go-tool/util/prettyjson"
)

func ToJsonString(v interface{}) []byte {
	if v == nil {
		return []byte{}
	}
	jsonByte, err := json.Marshal(v)
	if err != nil {
		return []byte{}
	}
	jsonByte, err = JsonFormat(jsonByte)
	if err != nil {
		return []byte{}
	}
	return jsonByte
}

var (
	JsonFormat = prettyjson.Format
)
