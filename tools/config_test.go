package tools

import (
	"testing"
)

func Test_GetConfig(t *testing.T) {
	config, err := GetConfig()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v\n", config)
}
