package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/easynetwork/aws-sdk-go-bindings/testdata"
)

func TestNew(t *testing.T) {

	cfg := testdata.MockConfiguration(t)

	in, inErr := NewSessionInput(cfg.Region)

	assert.NoError(t, inErr)
	assert.NotEmpty(t, in)

	svc, err := New(in)

	assert.NotEmpty(t, svc)
	assert.NoError(t, err)

	_, errNoRegionProvided := New(&SessionInput{
		region: "",
	})

	assert.Error(t, errNoRegionProvided)
	assert.Equal(t, ErrNoRegionProvided, errNoRegionProvided.Error())

}
