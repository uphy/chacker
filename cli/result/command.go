package result

import (
	"fmt"
	"io"
)

type (
	CommandResultBody struct {
		ExitCode int    `json:"exitCode"`
		StdOut   string `json:"stdout,omitempty"`
		StdErr   string `json:"stderr,omitempty"`
	}
)

func NewCommandResultBody() *CommandResultBody {
	return &CommandResultBody{}
}

func (c *CommandResultBody) JSON() interface{} {
	return c
}

func (c *CommandResultBody) Pretty(writer io.Writer) error {
	if c.ExitCode == 0 {
		if c.StdOut != "" {
			fmt.Fprintln(writer, c.StdOut)
		}
		if c.StdErr != "" {
			fmt.Fprintln(writer, c.StdErr)
		}
	} else {
		_, err := fmt.Fprintf(writer, "Failed to execute command.  (exitCode=%d, stdout=%s, stderr=%s)\n", c.ExitCode, c.StdOut, c.StdErr)
		return err
	}
	return nil
}
func (c *CommandResultBody) Plain(writer io.Writer) error {
	return c.Pretty(writer)
}
