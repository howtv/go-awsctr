package awsctr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// CountMetricsInfo -
type CountMetricsInfo struct {
	NameSpace      string
	DimensionName  string
	DimensionValue string
	MetricName     string
	Value          float64
}

// AlarmInfo -
type AlarmInfo struct {
	Name string
}

// CloudWatch -
type CloudWatch interface {
	PutCountMetrics(CountMetricsInfo) error
	AlarmOn(AlarmInfo) error
	AlarmOff(AlarmInfo) error
}

// NewCloudWatch -
func NewCloudWatch(sess *session.Session) CloudWatch {
	return &cloudWatchImpl{
		client: cloudwatch.New(sess),
	}
}

type cloudWatchService interface {
	PutMetricData(input *cloudwatch.PutMetricDataInput) (*cloudwatch.PutMetricDataOutput, error)
	EnableAlarmActions(input *cloudwatch.EnableAlarmActionsInput) (*cloudwatch.EnableAlarmActionsOutput, error)
	DisableAlarmActions(input *cloudwatch.DisableAlarmActionsInput) (*cloudwatch.DisableAlarmActionsOutput, error)
}

type cloudWatchImpl struct {
	client cloudWatchService
}

// PutCountMetrics -
// cli sample:
// aws cloudwatch put-metric-data --dimensions PluginId=bq_nginx --namespace fluentd --metric-name buffer_queue_length --value 1 --unit Count
func (c *cloudWatchImpl) PutCountMetrics(info CountMetricsInfo) error {
	_, err := c.client.PutMetricData(&cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String(info.DimensionName),
						Value: aws.String(info.DimensionValue),
					},
				},
				MetricName: aws.String(info.MetricName),
				Unit:       aws.String("Count"),
				Value:      aws.Float64(info.Value),
			},
		},
		Namespace: aws.String(info.NameSpace),
	})
	if err != nil {
		return err
	}

	return nil
}

// AlarmOn -
// cli sample:
// aws cloudwatch enable-alarm-actions --alarm-names <value>
func (c *cloudWatchImpl) AlarmOn(alarm AlarmInfo) error {
	_, err := c.client.EnableAlarmActions(&cloudwatch.EnableAlarmActionsInput{
		AlarmNames: []*string{
			aws.String(alarm.Name),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// AlarmOff -
// cli sample:
// aws cloudwatch disable-alarm-actions --alarm-names <value>
func (c *cloudWatchImpl) AlarmOff(alarm AlarmInfo) error {
	_, err := c.client.DisableAlarmActions(&cloudwatch.DisableAlarmActionsInput{
		AlarmNames: []*string{
			aws.String(alarm.Name),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
