package awsctr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewSession(region string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
}
