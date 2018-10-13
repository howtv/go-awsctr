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
	var paths cloudfront.Paths
	paths = *paths.SetQuantity(1)
	paths = *paths.SetItems([]*string{path})
	if err := paths.Validate(); err != nil {
		return err
	}

	unixTime := time.Now().Unix()
	var batch cloudfront.InvalidationBatch
	batch = *batch.SetCallerReference(fmt.Sprint(unixTime))
	batch = *batch.SetPaths(&paths)
	if err := batch.Validate(); err != nil {
		return err
	}

	var input cloudfront.CreateInvalidationInput
	input = *input.SetDistributionId(iv.DistID)
	input = *input.SetInvalidationBatch(&batch)
	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := c.client.CreateInvalidation(&input); err != nil {
		return err
	}

	return nil
}
