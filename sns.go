package awsctr

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// PublishInput -
type PublishInput = sns.PublishInput

// SNS -
type SNS interface {
	PublishSNS(input *PublishInput) (output *sns.PublishOutput, err error)
}

// NewSNS -
func NewSNS(sess *session.Session) SNS {
	return &snsImpl{
		client: sns.New(sess),
	}
}

type snsService interface {
	Publish(*sns.PublishInput) (*sns.PublishOutput, error)
}

type snsImpl struct {
	client snsService
}

func (c *snsImpl) PublishSNS(input *PublishInput) (output *sns.PublishOutput, err error) {
	output, err = c.client.Publish(input)
	if err != nil {
		return nil, err
	}
	return output, nil
}
