package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	awsctr "github.com/howtv/go-awsctr"
	"github.com/juju/errors"
	"github.com/mitchellh/cli"
)

func logsListCommandFactory() (cli.Command, error) {
	c := new(logsListCommand)
	c.fs = NewFlagSet(&c.Config)

	return c, nil
}

type logsListCommand struct {
	Config
	fs           *flag.FlagSet
	logGroupName string
}

func (c *logsListCommand) Help() string {
	b := new(bytes.Buffer)
	line := "awsctr logs list"
	return HelpMsg(b, line, c.fs)
}

func (c *logsListCommand) Synopsis() string {
	return "list log groups/log streams"
}

func (c *logsListCommand) parse(args []string) {
	c.fs.Parse(args)
	args = c.fs.Args()
	if len(args) > 0 {
		c.logGroupName = args[0]
		c.fs.Parse(args[1:])
	}
}

func (c *logsListCommand) Run(args []string) (exitStatus int) {
	c.parse(args)

	err := c.run(c.logGroupName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitStatus = 1
	}
	return exitStatus
}

func (c *logsListCommand) run(logGroupName string) error {
	if logGroupName == "" {
		return c.listLogGroups()
	}
	return errors.NotImplementedf("list log streams")
}

func (c *logsListCommand) listLogGroups() error {
	sess := awsctr.NewSession(c.Config.Region)
	client := awsctr.NewCloudWatchLogs(sess)

	groups, err := client.ListLogGroups()
	if err != nil {
		return err
	}
	c.printLogGroups(groups)
	return nil
}

func (c *logsListCommand) printLogGroups(groups []string) {
	for _, g := range groups {
		fmt.Println(g)
	}
}
