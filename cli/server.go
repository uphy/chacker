package cli

import (
	"bytes"
	"fmt"

	"github.com/labstack/echo"
	"github.com/spf13/cobra"
	"github.com/uphy/chacker/handlers"
)

func (c *CLI) server() *cobra.Command {
	cmd := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			e := echo.New()
			e.HideBanner = true
			runner := func(args []string) (interface{}, error) {
				cli := New(false)
				buffer := new(bytes.Buffer)
				cli.Out = buffer
				_, err := cli.Execute(args)
				if err != nil {
					return nil, err
				}
				return cli.resultAsJSON(), nil
			}
			runHandler := handlers.NewRunHandler(runner)
			e.POST("run", runHandler.Run)

			port, _ := cmd.Flags().GetInt("port")
			if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
				panic(err)
			}
			return nil
		},
	}
	cmd.Flags().IntP("port", "p", 8080, "server port")
	return cmd
}
