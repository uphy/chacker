package cli

import (
	"github.com/spf13/cobra"
	"github.com/uphy/chacker/cli/result"
)

func (c *CLI) dumpConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump",
		Short: "dump the resolved config",
		RunE: func(cmd *cobra.Command, args []string) error {
			body := result.NewMessageResultBody()
			c.config.Save(body.Out)
			c.setResult(result.New(body))
			return nil
		},
	}
	return cmd
}
