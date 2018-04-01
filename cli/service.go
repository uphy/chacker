package cli

import (
	"errors"
	"fmt"

	"github.com/uphy/chacker/cli/result"

	"github.com/spf13/cobra"
	"github.com/uphy/chacker/config"
)

func (c *CLI) service() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "list services, list commands, execute command.",
		Long: `"service" command list services, list commands, and execute command.

If no args specified, print the list of services.
If you specify an argument, print the list of the available command in the service.
If you specify two or more arguments; service and command, execute the command.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var res *result.Result
			switch len(args) {
			case 0:
				res = c.listServices()
			case 1:
				res = c.listCommands(args[0])
			default:
				host, _ := cmd.Flags().GetString("host")
				res = c.executeCommand(args[0], args[1], args[2:], host)
			}
			c.setResult(res)
			return nil
		},
	}
	cmd.Flags().String("host", "", "target host")
	return cmd
}

func (c *CLI) listServices() *result.Result {
	body := result.NewGridResultBody("name", "host")
	for _, service := range c.config.Services {
		body.Append(service.Name, service.Host)
	}
	return result.New(body)
}

func (c *CLI) listCommands(serviceName string) *result.Result {
	body := result.NewGridResultBody("name", "description")
	service, err := c.findService(serviceName)
	if err != nil {
		return result.NewError(err)
	}
	for _, command := range service.Commands {
		body.Append(command.Name, command.Description)
	}
	return result.New(body)
}

func (c *CLI) findService(serviceName string) (*config.ServiceConfig, error) {
	service, exist := c.config.Services[serviceName]
	if !exist {
		return nil, fmt.Errorf("no such service: %s", serviceName)
	}
	return &service, nil
}

func (c *CLI) findCommand(serviceName string, commandName string) (*config.CommandConfig, error) {
	service, err := c.findService(serviceName)
	if err != nil {
		return nil, err
	}
	command, exist := service.Commands[commandName]
	if !exist {
		return nil, fmt.Errorf("no such command: %s", commandName)
	}
	return &command, nil
}

func (c *CLI) executeCommand(serviceName string, commandName string, args []string, hostName string) *result.Result {
	command, err := c.findCommand(serviceName, commandName)
	if err != nil {
		return result.NewError(err)
	}
	if hostName == "" {
		hostName = command.Host
		if hostName == "" {
			return result.NewError(errors.New("command requires --host option"))
		}
	}
	host, exist := c.config.Hosts[hostName]
	if !exist {
		return result.NewError(fmt.Errorf("no such host: %s", hostName))
	}

	commandResult, err := c.executor.Execute(&host, command, args)
	if err != nil {
		return result.NewError(fmt.Errorf("failed to execute command: %v", err))
	}
	body := result.NewCommandResultBody()
	body.ExitCode = commandResult.ExitCode
	body.StdOut = commandResult.StdOut
	body.StdErr = commandResult.StdErr
	return result.New(body)
}
