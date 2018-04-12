package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/howtv/go-awsctr/awsctr"
	"github.com/juju/errors"
	"github.com/k0kubun/pp"
	"github.com/mitchellh/cli"
)

func logsWatchCommandFactory() (cli.Command, error) {
	c := new(logsWatchCommand)
	c.fs = NewFlagSet(&c.Config)
	c.fs.IntVar(&c.interval, "interval", 3, "api request interval seconds")
	c.fs.StringVar(&c.logFormat, "format", "", "log format currently supported (ecs)")
	c.fs.BoolVar(&c.noTrunc, "no-trunc", false, "output all log message")
	c.fs.StringVar(&c.filter, "filter", "", "filter string")

	return c, nil
}

type logsWatchCommand struct {
	Config
	fs           *flag.FlagSet
	interval     int
	noTrunc      bool
	filter       string
	logFormat    string
	logGroupName string
}

func (c *logsWatchCommand) Help() string {
	b := new(bytes.Buffer)
	b.WriteString("awsctr logs watch <logs-gorup> [options]\n")
	c.fs.SetOutput(b)
	c.fs.PrintDefaults()
	return b.String()
}

func (c *logsWatchCommand) Synopsis() string {
	return "watch cloudwatch logs"
}

func (c *logsWatchCommand) parse(args []string) {

	// care flag arg flag
	// ie. awsctr logs watch --interval=10 group-name --region=us-west-2
	c.fs.Parse(args)
	args = c.fs.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, c.Help())
		os.Exit(1)
	}
	c.logGroupName = args[0]
	c.fs.Parse(args[1:])
}

func (c *logsWatchCommand) Run(args []string) (exitStatus int) {

	c.parse(args)

	err := c.run(c.logGroupName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitStatus = 1
	}
	return exitStatus
}

func (c *logsWatchCommand) formatECSMessage(msg string) string {
	log, err := newECSLog(msg)
	if err != nil {
		// TODO Debug
		return ""
	}
	if c.noTrunc {
		panic("no-trunc option not implemented yet!")
	}
	innerLog := log.msg["log"].(string)
	return innerLog
}

func (c *logsWatchCommand) format(msg string) string {
	switch c.logFormat {
	case "ecs", "docker":
		return c.formatECSMessage(msg)
	}
	return msg
}

func (c *logsWatchCommand) filterp(msg string) bool {
	return !strings.Contains(msg, c.filter)
}

func (c *logsWatchCommand) run(logGroupName string) error {
	client, err := awsctr.NewCloudWatchLogs(c.Config.Region)
	if err != nil {
		return err
	}

	stream, err := client.FetchLatestStream(logGroupName)
	if err != nil {
		return err
	}

	logStreamName := aws.StringValue(stream.LogStreamName)
	if logStreamName == "" {
		return errors.Errorf("failed to fetch log stream from log group %s", logGroupName)
	}
	pp.Println("LogStreamName:", logStreamName)

	ctx := context.Background()
	defer ctx.Done()
	logsCh, errCh := client.OpenLogStream(ctx, logGroupName, logStreamName, time.Duration(c.interval)*time.Second)

	for {
		select {
		case err := <-errCh:
			return err
		case log := <-logsCh:
			msg := aws.StringValue(log.Message)
			if c.filterp(msg) {
				continue
			}
			msg = c.format(msg)
			if msg != "" {
				fmt.Println(msg)
			}
		}
	}
}

type ecsLog struct {
	raw               string
	timestampElements []string
	tag               string
	msgJSON           []string
	msg               map[string]interface{}
	innerLog          string
}

func newECSLog(msg string) (*ecsLog, error) {
	fields := strings.Fields(msg)
	if len(fields) <= 4 {
		return nil, errors.Errorf("invalid log format %s", msg)
	}

	l := &ecsLog{}
	// 2018-04-12 06:05;50.000000000 +0000
	l.timestampElements = fields[:3]
	l.tag = fields[3]
	l.msgJSON = fields[4:]
	l.raw = msg

	err := json.Unmarshal([]byte(strings.Join(l.msgJSON, "")), &l.msg)
	return l, err
}
