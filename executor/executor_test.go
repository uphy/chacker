package executor

import (
	"testing"

	"github.com/uphy/chacker/config"
)

func TestExecutorRun(t *testing.T) {
	e := New()
	if _, err := e.Execute(&config.HostConfig{
		Address:  "192.168.100.138",
		Port:     22,
		User:     "ishikura",
		Password: "ishikura",
	}, &config.CommandConfig{
		Script: `#!/bin/bash
ls -al /tmp
`,
	}, []string{}); err != nil {
		t.Fatal(err)
	}
	t.Error()
}
