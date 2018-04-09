package executor

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

	// create remote temporary directory/file
	tempDirResult, err := c.Exec("mktemp -d /tmp/chacker_XXXXX")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %v", err)
	}
	tempDir := strings.Trim(tempDirResult.StdOut, "\r\n")
	scriptFile := filepath.ToSlash(filepath.Join(tempDir, "script.sh"))
	downloadDir := filepath.ToSlash(filepath.Join(tempDir, "downloads"))
	defer c.Exec(fmt.Sprintf(`rm -rf "%s"`, tempDir))
	if _, err := c.Exec(fmt.Sprintf(`mkdir -p "%s"`, downloadDir)); err != nil {
		return nil, fmt.Errorf("failed to create download directory: %v", err)
	}

	// prepare the local command file
	var localFile string
	if command.Script != "" {
		f, err := ioutil.TempFile("", "chackerscript")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp file: %v", err)
		}
		localFile = f.Name()
		defer os.Remove(localFile)
		f.Close()
		if err := ioutil.WriteFile(localFile, []byte(appendShebang(command.Script)), 0700); err != nil {
			return nil, fmt.Errorf("failed to write the shell script: %v", err)
		}
	} else {
		localFile = command.File
	}

	// send command file
	if err := c.Upload(localFile, scriptFile, "0755"); err != nil {
		return nil, fmt.Errorf("failed to upload the shell script: %v", err)
	}

	// execute the command
	env := getEnvironmentVariables(command, tempDir)
	cmd := generateCommand(env, command.Directory, scriptFile, args)
	result, err := c.Exec(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the command: err=%v, stdout=%s, stderr=%s", err, result.StdOut, result.StdErr)
	}

	// download files in CHACKER_DOWNLOAD
	if err := c.Download(downloadDir, "./"); err != nil {
		return nil, fmt.Errorf("failed to download the CHACKER_DOWNLOAD directory: %v", err)
	}

	return result, err
}

func appendShebang(script string) string {
	if strings.HasPrefix(script, "#!") {
		return script
	}
	return fmt.Sprintf("#!/bin/sh\n%s", script)
}

func getEnvironmentVariables(command *config.CommandConfig, tempDir string) map[string]string {
	environment := command.Environment
	if environment == nil {
		environment = map[string]string{}
	}
	environment["CHACKER_DOWNLOAD"] = filepath.ToSlash(filepath.Join(tempDir, "downloads"))
	return environment
}

func generateCommand(env map[string]string, directory string, shellScriptFile string, args []string) string {
	// command name
	cmd := fmt.Sprintf(`"%s"`, shellScriptFile)

	if len(env) > 0 {
		// append environment variables
		envs := ""
		// sort
		keys := make([]string, 0)
		for k := range env {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		// join
		for _, k := range keys {
			envs += fmt.Sprintf(`%s="%s" `, k, env[k])
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
	if directory != "" {
		cmd = fmt.Sprintf(`cd "%s";%s`, directory, cmd)
	}

	return cmd
}
