package config

import (
	"fmt"
	"strings"
)

type (
	ServiceConfig struct {
		Name      string                   `yaml:"name"`
		Host      string                   `yaml:"host"`
		Directory string                   `yaml:"directory"`
		Commands  map[string]CommandConfig `yaml:"commands"`
		Template  *TemplateRefConfig       `yaml:"template,omitempty"`
	}
)

func (s *ServiceConfig) resolve(templates map[string]ServiceConfig) *ServiceConfig {
	resolved := *s
	if s.Template != nil {
		templateName := s.Template.Name

		t := templates[templateName]
		t.resolve(templates)
		templates[templateName] = t

		resolved = *mergeService(&t, s)
		resolved.substituteVariables(s.Template.Vars)
		resolved.Template = nil
	}
	resolved.Template = nil
	return &resolved
}

func (s *ServiceConfig) injectDefault() {
	for name, command := range s.Commands {
		if command.Directory == "" {
			command.Directory = s.Directory
		}
		if command.Host == "" {
			command.Host = s.Host
		}
		command.Name = name
		s.Commands[name] = command
	}
}

func (s *ServiceConfig) substituteVariables(vars map[string]string) {
	s.Directory = substitute(s.Directory, vars)
	s.Host = substitute(s.Host, vars)
	for commandName, command := range s.Commands {
		command.Script = substitute(command.Script, vars)
		command.Directory = substitute(command.Directory, vars)
		command.Host = substitute(command.Host, vars)
		s.Commands[commandName] = command
	}
}

func substitute(template string, vars map[string]string) string {
	s := template
	for name, value := range vars {
		s = strings.Replace(s, fmt.Sprintf("${%s}", name), value, -1)
	}
	return s
}

func mergeService(base *ServiceConfig, v *ServiceConfig) *ServiceConfig {
	merged := *v
	if merged.Directory == "" {
		merged.Directory = base.Directory
	}
	if merged.Host == "" {
		merged.Host = base.Host
	}
	merged.Commands = mergeCommands(base.Commands, v, merged.Commands)
	merged.injectDefault()
	return &merged
}

func mergeCommands(base map[string]CommandConfig, service *ServiceConfig, v map[string]CommandConfig) map[string]CommandConfig {
	merged := map[string]CommandConfig{}
	for name, command := range v {
		merged[name] = command
	}
	for baseCommandName, baseCommand := range base {
		command, exist := v[baseCommandName]
		if exist {
			mergedCommand := mergeCommand(&baseCommand, service, &command)
			merged[baseCommandName] = *mergedCommand
		} else {
			merged[baseCommandName] = baseCommand
		}
	}
	return merged
}

func mergeCommand(base *CommandConfig, service *ServiceConfig, cmd *CommandConfig) *CommandConfig {
	merged := *cmd

	if merged.Script == "" {
		merged.Script = base.Script
	}
	if merged.File == "" {
		merged.File = base.File
	}
	if merged.Directory == "" {
		merged.Directory = base.Directory
	}
	if merged.Host == "" {
		merged.Host = base.Host
	}
	if merged.Description == "" {
		merged.Description = base.Description
	}
	for k, v := range base.Environment {
		if _, exists := merged.Environment[k]; !exists {
			merged.Environment[k] = v
		}
	}
	return &merged
}
