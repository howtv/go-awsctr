package awsctr

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// ECS -
type ECS interface {
	GetContainerInstances(clusterName string) (instances *ecs.ListContainerInstancesOutput, err error)
	GetContainerInstancesStatus(clusterName string, instanceArns []*string) (status *ecs.DescribeContainerInstancesOutput, err error)
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

func (c *ecsImpl) GetContainerInstances(clusterName string) (instances *ecs.ListContainerInstancesOutput, err error) {
	stat := containerInstanceStatus
	instances, err = c.client.ListContainerInstances(&ecs.ListContainerInstancesInput{
		Cluster: &clusterName,
		Status:  &stat,
	})
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (c *ecsImpl) GetContainerInstancesStatus(clusterName string, instanceArns []*string) (status *ecs.DescribeContainerInstancesOutput, err error) {
	status, err = c.client.DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{
		Cluster:            &clusterName,
		ContainerInstances: instanceArns,
	})
	if err != nil {
		return nil, err
	}
	return status, nil
}
