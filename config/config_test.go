package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	content := `
hosts:
  host1:
    address: 192.168.11.2
    port: 22
    user: user1
    key: keys/id_rsa_user1
  host2:
    address: 192.168.11.3
    port: 22
    user: user1
    key: keys/id_rsa_user1
templates:
  docker-compose:
    commands:
      start:
        file: scripts/start.sh
      stop:
        file: scripts/stop.sh
      restart:
        file: scripts/restart.sh
      update:
        file: scripts/update.sh  
  myservice:
    directory: /var/lib/${name}
    template:
      name: docker-compose
services:
  service1:
    host: host1
    template: 
      name: myservice
      vars:
        name: service1
  service2:
    host: host2
    directory: /var/lib/service2    
    template:
      name: docker-compose
`
	c, err := LoadConfig(strings.NewReader(content))
	if err != nil {
		t.Error(err)
	}
	if c.Templates != nil {
		t.Error("templates must be resolved")
	}
	c.Save(os.Stdout)
	service1, ok := c.Services["service1"]
	if !ok {
		t.Error("expected service not exist")
	}
	if service1.Template != nil {
		t.Error("template not resolved")
	}
	cmd, ok := service1.Commands["start"]
	if !ok {
		t.Error("command config not resolved")
	}
	if cmd.Command == "" {
		t.Error("command not resolved")
	}
	if cmd.Command != "scripts/start.sh" {
		t.Error("template not resolved")
	}
}
