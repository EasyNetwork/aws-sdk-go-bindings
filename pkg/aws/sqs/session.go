package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	pkgAws "github.com/easynetwork/aws-sdk-go-bindings/pkg/aws"
)

// SQS embeds *sns.SNS
type SQS struct {
	*sqs.SQS
}

// New returns a new *SQS
func New(svc *pkgAws.Session, endpoint string) (*SQS, error) {

	if len(endpoint) > 0 {
		svc.Config.Endpoint = aws.String(endpoint)
	}

	newSvc, newSvcErr := session.NewSession(svc.Config)
	if newSvcErr != nil {
		return nil, newSvcErr
	}

	sqsSvc := new(SQS)
	sqsSvc.SQS = sqs.New(newSvc)

	return sqsSvc, nil

}
