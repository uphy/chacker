package executor

import (
	"github.com/uphy/chacker/config"
)

type Executor struct {
}

type CommandResult struct {
	ExitCode int
	StdOut   string
	StdErr   string
}

func (e *Executor) Execute(host *config.HostConfig, command *config.CommandConfig, args []string) (*CommandResult, error) {
	return &CommandResult{0, "aaa", ""}, nil
}

func (e *Executor) generateCommandFile(command *config.CommandConfig) (file string, err error) {
	return "", nil
}

func (e *Executor) send(local string, remote string) error {
	return nil
}

func (e *Executor) executeRemote(remote string) error {
	return nil
}

func (e *Executor) deleteRemote(remote string) error {
	return nil
}
