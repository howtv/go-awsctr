package awsctr

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codepipeline"
)

// JobInfo -
type JobInfo struct {
	ID        string
	IsSuccess bool
}

// CodePipeline -
type CodePipeline interface {
	SendJobSuccess(JobInfo) error
	SendJobFailure(JobInfo) error
}

// NewCodePipeline -
func NewCodePipeline(sess *session.Session) CodePipeline {
	return &codePipelineImpl{
		client: codepipeline.New(sess),
	}
}

type codePipelineService interface {
	PutJobSuccessResult(input *codepipeline.PutJobSuccessResultInput) (*codepipeline.PutJobSuccessResultOutput, error)
	PutJobFailureResult(input *codepipeline.PutJobFailureResultInput) (*codepipeline.PutJobFailureResultOutput, error)
}

type codePipelineImpl struct {
	client codePipelineService
}

// SendJobSuccess -
// cli sample:
// aws codepipeline put-job-success-result --job-id e930bc23-49f3-442d-8189-8c33889a1791
func (c *codePipelineImpl) SendJobSuccess(job JobInfo) error {

	var input codepipeline.PutJobSuccessResultInput
	input = *input.SetJobId(job.ID)
	if err := input.Validate(); err != nil {
		return err
	}
	if _, err := c.client.PutJobSuccessResult(&input); err != nil {
		return err
	}

	return nil
}

// SendJobFailure -
// cli sample:
// aws codepipeline put-job-failure-result --job-id e930bc23-49f3-442d-8189-8c33889a1791
func (c *codePipelineImpl) SendJobFailure(job JobInfo) error {

	var input codepipeline.PutJobFailureResultInput
	input = *input.SetJobId(job.ID)
	if err := input.Validate(); err != nil {
		return err
	}
	if _, err := c.client.PutJobFailureResult(&input); err != nil {
		return err
	}

	return nil
}
