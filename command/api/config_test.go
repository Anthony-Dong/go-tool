package api

import (
	"testing"
)

func TestCommonConfig_ReadConfig(t *testing.T) {
	config := CommonConfig{
		Config: "/Users/fanhaodong/go/bin/upload-config.json",
	}
	bytes, err := config.ReadConfig("default")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n", bytes)
}
