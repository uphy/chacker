package executor

import (
	"testing"

	"github.com/uphy/chacker/config"
)

func TestExecutorRun(t *testing.T) {
	e := New()
	if r, err := e.Execute(&config.HostConfig{
		Address:  "192.168.100.138",
		Port:     22,
		User:     "ishikura",
		Password: "ishikura",
	}, &config.CommandConfig{
		Script: `#!/bin/bash
ls -al /tmp
env
`,
		Environment: map[string]string{
			"foo": "1",
			"bar": "aaa",
		},
	}, []string{"arg1", "arg2"}); err != nil {
		t.Fatal(err)
	} else {
		if r == nil {
			t.Fatal("result == nil")
		}
	}
}
