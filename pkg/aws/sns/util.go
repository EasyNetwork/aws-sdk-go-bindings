package sns

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/service/sns"

	intErr "github.com/easynetwork/aws-sdk-go-bindings/internal/error"
)

// Body is used to initialize a valid SNS message
type Body struct {
	Default string `json:"default"`
}

// NewPublishInput returns a new *PublishInput given a body and an endpoint
func NewPublishInput(input interface{}, endpoint string) (*sns.PublishInput, error) {

	if endpoint == "" {
		return nil, intErr.Format(Endpoint, ErrEmptyParameter)
	}

	if reflect.ValueOf(input).Kind() == reflect.Ptr {
		return nil, intErr.Format(Input, ErrPointerParameterNotAllowed)
	}

	inBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	// Mandatory since SNS needs escaped `"`. So we need to escape them to `\"`
	unquote := strings.Replace(string(inBytes), `"`, "\"", -1)

	// Mandatory since SNS needs bodies like:
	// {
	// 		"default" : {
	// 			\"par1\" : \"some value\"
	// 		}
	// }
	snsBody := Body{
		Default: unquote,
	}

	// Mandatory since we want to get a string out of encoded bytes
	msgBytes, err := json.Marshal(snsBody)
	if err != nil {
		return nil, err
	}

	out := &sns.PublishInput{}
	out = out.SetMessage(string(msgBytes))
	out = out.SetMessageStructure(MessageStructure)
	out = out.SetTargetArn(endpoint)

	return out, nil

}

// UnmarshalMessage unmarshal an SNS Message to a given interface
func UnmarshalMessage(message string, input interface{}) error {

	if message == "" {
		return intErr.Format(Message, ErrEmptyParameter)
	}

	if reflect.ValueOf(input).Kind() != reflect.Ptr {
		return intErr.Format(Input, ErrNoPointerParameter)
	}

	uS := unescapeMessageString(message)

	err := json.Unmarshal([]byte(uS), input)
	if err != nil {
		return err
	}

	return nil

}

// unescapeMessageString takes a SNS message string like
// `"{\"stuff\" : \"somevalue\"}"` and outputs `"{"stuff" : "somevalue"}"`
func unescapeMessageString(in string) string {
	return strings.Replace(in, `\"`, `"`, -1)
}
