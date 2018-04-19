package awsctr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/juju/errors"
)

// ContainerInstance -
type ContainerInstance struct {
	ID        string
	Connected bool
}

// ContainerInstances -
type ContainerInstances []ContainerInstance

// ECS -
type ECS interface {
	GetContainerInstances(clusterName string) (insts ContainerInstances, err error)
}

// NewECS -
func NewECS(sess *session.Session) ECS {
	return &ecsImpl{
		client: ecs.New(sess),
	}
}

type ecsService interface {
	ListContainerInstances(*ecs.ListContainerInstancesInput) (*ecs.ListContainerInstancesOutput, error)
	DescribeContainerInstances(*ecs.DescribeContainerInstancesInput) (*ecs.DescribeContainerInstancesOutput, error)
}

type ecsImpl struct {
	client ecsService
}

func (c *ecsImpl) GetContainerInstances(clusterName string) (insts ContainerInstances, err error) {
	instanceArns, err := c.getContainerInstanceArns(clusterName)
	if err != nil {
		return nil, err
	}
	status, err := c.getContainerInstancesStatus(clusterName, instanceArns)
	if err != nil {
		return nil, err
	}

	insts, err = getInstances(status)
	if err != nil {
		return nil, err
	}

	return insts, nil
}

func (c *ecsImpl) getContainerInstanceArns(clusterName string) (instanceArns []*string, err error) {
	stat := containerInstanceStatus
	instances, err := c.client.ListContainerInstances(&ecs.ListContainerInstancesInput{
		Cluster: &clusterName,
		Status:  &stat,
	})
	if err != nil {
		return nil, err
	}

	for _, arn := range instances.ContainerInstanceArns {
		instanceArns = append(instanceArns, arn)
	}

	return instanceArns, nil
}

func (c *ecsImpl) getContainerInstancesStatus(clusterName string, instanceArns []*string) (status *ecs.DescribeContainerInstancesOutput, err error) {
	status, err = c.client.DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{
		Cluster:            &clusterName,
		ContainerInstances: instanceArns,
	})
	if err != nil {
		return nil, err
	}
	return status, nil
}

func getInstances(status *ecs.DescribeContainerInstancesOutput) (insts ContainerInstances, err error) {
	for _, stat := range status.ContainerInstances {
		inst := ContainerInstance{
			ID:        aws.StringValue(stat.Ec2InstanceId),
			Connected: aws.BoolValue(stat.AgentConnected),
		}
		if inst.ID != "" {
			insts = append(insts, inst)
		} else {
			return nil, errors.Errorf("connection is failure but missing instance id: %v", inst)
		}
	}
	return insts, nil
}
