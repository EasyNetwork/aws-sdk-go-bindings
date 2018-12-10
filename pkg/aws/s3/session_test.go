package s3

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/easynetwork/aws-sdk-go-bindings/pkg/aws"
	"github.com/easynetwork/aws-sdk-go-bindings/testdata"
)

func TestNew(t *testing.T) {

	cfg := testdata.MockConfiguration(t)

	in, err := aws.NewSessionInput(cfg.Region)

	assert.NoError(t, err)
	awsSvc, err := aws.New(in)

	assert.NoError(t, err)
	assert.NotEmpty(t, awsSvc)

	snsSvc, err := New(awsSvc, cfg.S3.Endpoint)

	assert.NoError(t, err)
	assert.NotEmpty(t, snsSvc)

}
