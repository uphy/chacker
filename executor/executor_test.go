package executor

import (
	"testing"

	"github.com/uphy/chacker/executor/sshtest"

	"github.com/uphy/chacker/config"
)

func TestExecutorExecute(t *testing.T) {
	sshtest.SSHServer(t, &sshtest.Result{
		StdOut: "hello\n",
	}, func() {
		e := New()
		if r, err := e.Execute(&config.HostConfig{
			Address:  sshtest.Host,
			Port:     sshtest.Port,
			User:     "user1",
			Password: "user1",
		}, &config.CommandConfig{
			Script: `#!/bin/bash
echo hello
`}, []string{}); err != nil {
			t.Fatal(err)
		} else {
			if r == nil {
				t.Fatal("result == nil")
			}
		}
	})
}

func TestExecutorExecuteDialFailure(t *testing.T) {
	e := New()
	if _, err := e.Execute(&config.HostConfig{
		Address:  "255.255.255.255",
		Port:     sshtest.Port,
		User:     "user1",
		Password: "user1",
	}, nil, nil); err == nil {
		t.Error("error should be returned")
	}
}

func TestGenerateCommand(t *testing.T) {
	if cmd := generateCommand(&config.CommandConfig{}, "test.sh", []string{}); cmd != `"test.sh"` {
		t.Error("unexpected command generated: ", cmd)
	}
}

func TestGenerateCommandWithEnvironment(t *testing.T) {
	if cmd := generateCommand(&config.CommandConfig{
		Environment: map[string]string{
			"foo": "1",
			"bar": "aaa",
		},
	}, "test.sh", []string{}); cmd != `bar="aaa" foo="1" "test.sh"` {
		t.Error("unexpected command generated: ", cmd)
	}
}

func TestGenerateCommandWithDirectory(t *testing.T) {
	if cmd := generateCommand(&config.CommandConfig{
		Directory: "/home/user1/dir",
	}, "test.sh", []string{}); cmd != `cd "/home/user1/dir";"test.sh"` {
		t.Error("unexpected command generated: ", cmd)
	}
}

func TestGenerateCommandWithArguments(t *testing.T) {
	if cmd := generateCommand(&config.CommandConfig{}, "test.sh", []string{"arg 1", "arg 2"}); cmd != `"test.sh" "arg 1" "arg 2"` {
		t.Error("unexpected command generated: ", cmd)
	}
}

func TestGenerateCommandWithAll(t *testing.T) {
	if cmd := generateCommand(&config.CommandConfig{
		Directory: "/home/user1/dir",
		Environment: map[string]string{
			"foo": "1",
			"bar": "aaa",
		},
	}, "test.sh", []string{"arg 1", "arg 2"}); cmd != `cd "/home/user1/dir";bar="aaa" foo="1" "test.sh" "arg 1" "arg 2"` {
		t.Error("unexpected command generated: ", cmd)
	}
}

func TestAppendShebang(t *testing.T) {
	if s := appendShebang("echo hello"); s != "#!/bin/sh\necho hello" {
		t.Error()
	}
	if s := appendShebang("#!/bin/sh\necho hello"); s != "#!/bin/sh\necho hello" {
		t.Error()
	}
	if s := appendShebang("#!/bin/bash\necho hello"); s != "#!/bin/bash\necho hello" {
		t.Error()
	}
}
