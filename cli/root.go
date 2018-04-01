package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/uphy/chacker/cli/result"
	"github.com/uphy/chacker/config"
	"github.com/uphy/chacker/executor"
)

type CLI struct {
	root     *cobra.Command
	config   *config.Config
	Out      io.Writer
	executor *executor.Executor
	res      *result.Result
}

func New(serverEnabled bool) *CLI {
	root := &cobra.Command{
		Use: "chacker",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	root.PersistentFlags().StringP("config", "c", "config.yml", "config file")
	root.PersistentFlags().Bool("pretty", false, "enable pretty print")

	cli := &CLI{
		root: root,
		Out:  os.Stdout,
		res:  result.NewError(errors.New("internal error: no result set")),
	}
	cobra.OnInitialize(func() {
		configFile, _ := root.Flags().GetString("config")
		conf, err := config.LoadConfigFile(configFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cli.config = conf
	})

	root.AddCommand(cli.service())
	if serverEnabled {
		root.AddCommand(cli.server())
	}
	root.AddCommand(cli.dumpConfig())
	return cli
}

func (c *CLI) setResult(result *result.Result) {
	if result == nil {
		panic("result == nil")
	}
	c.res = result
}

func (c *CLI) pretty() bool {
	pretty, _ := c.root.Flags().GetBool("pretty")
	return pretty
}

func (c *CLI) resultAsJSON() interface{} {
	v := c.res.JSON()
	message := new(bytes.Buffer)
	if c.pretty() {
		c.res.Pretty(message)
	} else {
		c.res.Plain(message)
	}
	messageStr := message.String()
	messageStr = strings.TrimRight(messageStr, "\r\n")
	v["message"] = messageStr
	return v
}

func (c *CLI) Execute(args []string) (*result.Result, error) {
	c.root.SetArgs(args)
	if err := c.root.Execute(); err != nil {
		return nil, err
	}
	if c.pretty() {
		c.res.Pretty(c.Out)
	} else {
		c.res.Plain(c.Out)
	}
	return c.res, nil
}
