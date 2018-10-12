package awsctr

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

// InvalidateInfo -
type InvalidateInfo struct {
	DistID string
	Path   string
}

// CloudFront -
type CloudFront interface {
	Invalidate(InvalidateInfo) error
}

// NewCloudFront -
func NewCloudFront(sess *session.Session) CloudFront {
	return &cloudFrontImpl{
		client: cloudfront.New(sess),
	}
}

type cloudFrontService interface {
	CreateInvalidation(input *cloudfront.CreateInvalidationInput) (*cloudfront.CreateInvalidationOutput, error)
}

type cloudFrontImpl struct {
	client cloudFrontService
}

// Invalidate -
// cli sample:
// aws cloudfront create-invalidation --distribution-id S11A16G5KZMEQD --paths /index.html /error.html
func (c *cloudFrontImpl) Invalidate(iv InvalidateInfo) error {
	path := aws.String(iv.Path)
	unixTime := time.Now().Unix()
	_, err := c.client.CreateInvalidation(&cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(iv.DistID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(fmt.Sprint(unixTime)),
			Paths: &cloudfront.Paths{
				Items:    []*string{path},
				Quantity: aws.Int64(1),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
