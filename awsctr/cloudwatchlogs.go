package awsctr

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/juju/errors"
)

const (
	logStreamOrderEventTime  = "LastEventTime"
	logStreamOrderStreamName = "logStreamName"

	getLogEventsLimit      = 10000
	describeLogGroupsLimit = 50
	logStreamBuffer        = 100
)

// CloudWatchLogs -
type CloudWatchLogs interface {
	FetchLatestStream(logGroupName string) (*cloudwatchlogs.LogStream, error)
	OpenLogStream(ctx context.Context, logGroupName, logStreamName string, interval time.Duration) (<-chan *cloudwatchlogs.OutputLogEvent, <-chan error)
	ListLogGroups() ([]string, error)
}

// NewCloudWatchLogs -
func NewCloudWatchLogs(region string) (CloudWatchLogs, error) {
	sess, err := NewSession(region)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &cloudWatchLogs{
		client: cloudwatchlogs.New(sess),
	}, nil
}

type cloudWatchLogsService interface {
	DescribeLogStreams(*cloudwatchlogs.DescribeLogStreamsInput) (*cloudwatchlogs.DescribeLogStreamsOutput, error)
	GetLogEvents(*cloudwatchlogs.GetLogEventsInput) (*cloudwatchlogs.GetLogEventsOutput, error)
	DescribeLogGroups(*cloudwatchlogs.DescribeLogGroupsInput) (*cloudwatchlogs.DescribeLogGroupsOutput, error)
}

type cloudWatchLogs struct {
	client cloudWatchLogsService
}

// interfaceの実装
func (c *cloudWatchLogs) fetchLatestStreamInput(logGroupName string) (*cloudwatchlogs.DescribeLogStreamsInput, error) {
	if logGroupName == "" {
		return nil, InvalidParamErr
	}
	return &cloudwatchlogs.DescribeLogStreamsInput{
		Descending:   aws.Bool(true),
		Limit:        aws.Int64(1),
		LogGroupName: aws.String(logGroupName),
		OrderBy:      aws.String(logStreamOrderEventTime),
	}, nil
}

func (c *cloudWatchLogs) fetchLatestStreamOutput(output *cloudwatchlogs.DescribeLogStreamsOutput) (*cloudwatchlogs.LogStream, error) {
	if output == nil {
		return nil, LogicErr
	}

	streams := output.LogStreams
	if len(streams) == 0 {
		return nil, NotFoundErr
	}

	return streams[0], nil
}

// FetchLatestStream -
func (c *cloudWatchLogs) FetchLatestStream(logGroupName string) (*cloudwatchlogs.LogStream, error) {
	input, err := c.fetchLatestStreamInput(logGroupName)
	if err != nil {
		return nil, err
	}
	output, err := c.client.DescribeLogStreams(input)
	if err != nil {
		return nil, err
	}
	return c.fetchLatestStreamOutput(output)
}

func (c *cloudWatchLogs) openLogStreamInput(logGroupName, logStreamName string, startTime time.Time) (*cloudwatchlogs.GetLogEventsInput, error) {
	if logGroupName == "" || logStreamName == "" {
		return nil, InvalidParamErr
	}
	return &cloudwatchlogs.GetLogEventsInput{
		Limit:         aws.Int64(getLogEventsLimit),
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
		StartFromHead: aws.Bool(true),
		StartTime:     aws.Int64(ToMilliseconds(startTime)),
	}, nil
}

func (c *cloudWatchLogs) openLogStream(
	ctx context.Context,
	logsCh chan *cloudwatchlogs.OutputLogEvent,
	errCh chan error,
	input *cloudwatchlogs.GetLogEventsInput,
	ticker *time.Ticker) {

	defer ticker.Stop()

	fetch := func() (next *string, err error) {
		output, err := c.client.GetLogEvents(input)
		if err != nil {
			return nil, err
		}

		events := output.Events
		for _, event := range events {
			logsCh <- event
		}
		return output.NextForwardToken, nil
	}

	for {
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		case <-ticker.C:
			next, err := fetch()
			if err != nil {
				errCh <- err
				return
			}
			input.NextToken = next
		}
	}
}

// OpenLogStream keep sending logEvents to channel
func (c *cloudWatchLogs) OpenLogStream(ctx context.Context, logGroupName, logStreamName string, interval time.Duration) (<-chan *cloudwatchlogs.OutputLogEvent, <-chan error) {
	var logsCh = make(chan *cloudwatchlogs.OutputLogEvent, logStreamBuffer)
	var errCh = make(chan error, 1)

	input, err := c.openLogStreamInput(logGroupName, logStreamName, time.Now())
	if err != nil {
		errCh <- err
		return nil, errCh
	}

	ticker := time.NewTicker(interval)
	go c.openLogStream(ctx, logsCh, errCh, input, ticker)

	return logsCh, errCh
}

func (c *cloudWatchLogs) listLogGroupsInput() (*cloudwatchlogs.DescribeLogGroupsInput, error) {
	return &cloudwatchlogs.DescribeLogGroupsInput{
		Limit: aws.Int64(describeLogGroupsLimit),
	}, nil
}

func (c *cloudWatchLogs) listLogGroupsOutput(input *cloudwatchlogs.DescribeLogGroupsInput) ([]string, error) {
	if input == nil {
		return nil, LogicErr
	}

	var groups []string
	for {
		output, err := c.client.DescribeLogGroups(input)
		if err != nil {
			return nil, err
		}
		for _, g := range output.LogGroups {
			groups = append(groups, aws.StringValue(g.LogGroupName))
		}
		if output.NextToken == nil {
			break
		}
		input.NextToken = output.NextToken
	}
	return groups, nil
}

func (c *cloudWatchLogs) ListLogGroups() ([]string, error) {
	input, err := c.listLogGroupsInput()
	if err != nil {
		return nil, err
	}
	return c.listLogGroupsOutput(input)
}
