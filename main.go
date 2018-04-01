package main

import (
	"os"

	"github.com/uphy/chacker/cli"
)

func main() {
	c := cli.New(true)
	c.Execute(os.Args[1:])
}
