package sns

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/service/sns"

	intErr "github.com/easynetwork/aws-sdk-go-bindings/internal/error"
)

const (
	// MessageAttributesDataTypes
	messageAttributeString      = "String"
	messageAttributeStringArray = "String.Array"
	messageAttributeNumber      = "Number"
	messageAttributesBinary     = "Binary"
)

// Body is used to initialize a valid SNS message
type Body struct {
	Default string `json:"default"`
}

// NewPublishInput returns a new *PublishInput given a body and an endpoint
func NewPublishInput(input interface{}, messageAttributes map[string]interface{}, endpoint string) (*sns.PublishInput, error) {

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

	awsMessageAttributes, err := getAsSNSMessageAttributes(messageAttributes)
	if err != nil {
		return nil, err
	}

	out := &sns.PublishInput{}
	out = out.SetMessage(string(msgBytes))
	out = out.SetMessageAttributes(awsMessageAttributes)
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

func getAsSNSMessageAttributes(messageAttributes map[string]interface{}) (map[string]*sns.MessageAttributeValue, error) {

	output := make(map[string]*sns.MessageAttributeValue)

	for k, v := range messageAttributes {
		vBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		output[k] = new(sns.MessageAttributeValue)

		// Setting data type, since it is required
		switch t := v.(type) {
		case string:
			output[k] = output[k].
				SetStringValue(t).
				SetDataType(messageAttributeString)
		case int, float32, float64:
			output[k] = output[k].
				SetStringValue(string(vBytes)).
				SetDataType(messageAttributeNumber)
		case []string:
			output[k] = output[k].
				SetStringValue(string(vBytes)).
				SetDataType(messageAttributeStringArray)
		case interface{}:
			output[k] = output[k].
				SetBinaryValue(vBytes).
				SetDataType(messageAttributesBinary)
		default:
			return nil, fmt.Errorf("%v is not a supported type", t)
		}
		// Checking message attributes validity
		if err := output[k].Validate(); err != nil {
			return nil, err
		}
	}

	return output, nil
}
