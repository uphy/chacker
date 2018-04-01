# Chacker

Remote service controller especially for ChatOps.

Chacker replaces routine tasks with simple and unified commands based on a config file.
The config file mainly consists of 2 parts; 'hosts' and 'services'.
'hosts' contains the SSH connection information and it will be used in commands.
'services' contains the service information.  A service has the commands.
'commands' can be defined with shell script or file.

## Getting Started

Create the first config file below.

```yaml
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
services:
  apache:
    host: host1
    commands:
      start:
        command: apachectl start
      stop:
        command: apachectl graceful-stop
      say:
        command: echo $1
```

Print the service list.

```bash
$ chacker service
name host
apache host1
```

Print the commands.

```bash
$ chacker service apache
name description
start 
stop 
```

Start the service.

```bash
$ chacker service apache start
```

Optionally, you can specify the hostname.

```bash
$ chacker service apache start --host host2
```

You can pass the arguments to the commands.

```bash
$ chacker service apache say "hello world"
```

## As a backend for Hubot

Start the chacker server.

```bash
$ chacker server
```

You can execute chacker command via HTTP.

```bash
$ curl -XPOST localhost:8080/run -F "command=service apache start"
{"body":{"exitCode":0,"stdout":""},"message":""}
```

Same as above example, you can send HTTP request from Hubot.
