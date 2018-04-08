package executor

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
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
		if err := ioutil.WriteFile(localFile, []byte(appendShebang(command.Script)), 0700); err != nil {
			return nil, err
		}
	} else {
		localFile = command.File
	}

	// send command file
	if err := c.Upload(localFile, tempFile, "0755"); err != nil {
		return nil, err
	}

	// execute the command
	return c.Exec("generateCommand(command, tempFile, args)")
}

func appendShebang(script string) string {
	if strings.HasPrefix(script, "#!") {
		return script
	}
	return fmt.Sprintf("#!/bin/sh\n%s", script)
}

func generateCommand(command *config.CommandConfig, shellScriptFile string, args []string) string {
	// command name
	cmd := fmt.Sprintf(`"%s"`, shellScriptFile)

	if len(command.Environment) > 0 {
		// append environment variables
		envs := ""
		// sort
		keys := make([]string, 0)
		for k := range command.Environment {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		// join
		for _, k := range keys {
			envs += fmt.Sprintf(`%s="%s" `, k, command.Environment[k])
		}
		cmd = envs + cmd
	}

	if len(args) > 0 {
		// append arguments
		s := ""
		for _, arg := range args {
			s += fmt.Sprintf(` "%s"`, arg)
		}
		cmd = cmd + s
	}

	// change directory
	if command.Directory != "" {
		cmd = fmt.Sprintf(`cd "%s";%s`, command.Directory, cmd)
	}

	return cmd
}
