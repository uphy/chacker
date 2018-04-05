package executor

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/uphy/chacker/config"
)

type Executor struct {
}

type CommandResult struct {
	ExitStatus int
	StdOut     string
	StdErr     string
}

func New() *Executor {
	return &Executor{}
}

func (e *Executor) Execute(host *config.HostConfig, command *config.CommandConfig, args []string) (*CommandResult, error) {
	c, err := NewSSHClientFromHostConfig(host)
	if err != nil {
		return nil, err
	}

	// create remote temporary file
	tempFileResult, err := c.Exec("mktemp /tmp/chackerscript_XXXXX")
	if err != nil {
		return nil, err
	}
	tempFile := strings.Trim(tempFileResult.StdOut, "\r\n")
	defer c.Exec(fmt.Sprintf("rm -f %s", tempFile))
	if _, err := c.Exec(fmt.Sprint("chmod 700 ", tempFile)); err != nil {
		return nil, err
	}

	// prepare the local command file
	var localFile string
	if command.Script != "" {
		f, err := ioutil.TempFile("", "chackerscript")
		if err != nil {
			return nil, err
		}
		localFile = f.Name()
		defer os.Remove(localFile)
		f.Close()
		if err := ioutil.WriteFile(localFile, []byte(command.Script), 0700); err != nil {
			return nil, err
		}
	} else {
		localFile = command.File
	}

	// send command file
	if err := c.Copy(localFile, tempFile, "0755"); err != nil {
		return nil, err
	}

	// execute the command
	cmd := tempFile
	if len(command.Environment) > 0 {
		envs := ""
		for k, v := range command.Environment {
			envs += fmt.Sprint(k, "=", v, " ")
		}
		cmd = envs + cmd
	}
	if command.Directory != "" {
		cmd = fmt.Sprintf("cd %s; %s", command.Directory, cmd)
	}
	return c.Exec(cmd)
}
