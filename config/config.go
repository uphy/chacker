package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Hosts     map[string]HostConfig    `yaml:"hosts"`
		Templates map[string]ServiceConfig `yaml:"templates,omitempty"`
		Services  map[string]ServiceConfig `yaml:"services"`
	}
	HostConfig struct {
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Key      string `yaml:"key"`
		Password string `yaml:"password"`
	}
	TemplateRefConfig struct {
		Name string            `yaml:"name"`
		Vars map[string]string `yaml:"vars,omitempty"`
	}
	CommandConfig struct {
		Name        string            `yaml:"name"`
		Script      string            `yaml:"script,omitempty"`
		File        string            `yaml:"file,omitempty"`
		Host        string            `yaml:"host"`
		Directory   string            `yaml:"directory"`
		Description string            `yaml:"description"`
		Environment map[string]string `yaml:"environment"`
	}
)

func (c *Config) Save(writer io.Writer) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}
	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}
	return nil
}

func LoadConfigFile(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()
	return LoadConfig(f)
}

func LoadConfig(reader io.Reader) (*Config, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	// set default values
	for name, host := range c.Hosts {
		if host.Port <= 0 {
			host.Port = 22
		}
		c.Hosts[name] = host
	}
	// resolve template of template
	for name, template := range c.Templates {
		template.Name = name
		s := &template
		s.injectDefault()
		s = s.resolve(c.Templates)
		c.Templates[name] = *s
	}
	// resolve template of services
	for name, service := range c.Services {
		service.Name = name
		s := &service
		s.injectDefault()
		s = s.resolve(c.Templates)
		c.Services[name] = *s
	}
	// templates no longer required
	c.Templates = nil

	// verify config
	if err := c.verify(); err != nil {
		return nil, fmt.Errorf("verification error: %v", err)
	}

	return &c, nil
}

func (c *Config) verify() error {
	for _, service := range c.Services {
		for _, command := range service.Commands {
			_, exist := c.Hosts[command.Host]
			if !exist {
				return fmt.Errorf("host doesn't exist. (service=%s, command=%s, host=%s)", service.Name, command.Name, command.Host)
			}
			if command.Script == "" && command.File == "" {
				return fmt.Errorf("Specify either 'command' or 'file'. (service=%s, command=%s, host=%s)", service.Name, command.Name, command.Host)
			}
			if command.Script != "" && command.File != "" {
				return fmt.Errorf("Can not specify both 'command' and 'file'. (service=%s, command=%s, host=%s)", service.Name, command.Name, command.Host)
			}
		}
	}
	return nil
}
