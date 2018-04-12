package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("awsctr", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"logs watch": logsWatchCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(exitStatus)
}
