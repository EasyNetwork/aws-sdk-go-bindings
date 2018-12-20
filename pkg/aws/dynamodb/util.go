package dynamodb

import (
	"encoding/json"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	intError "github.com/easynetwork/aws-sdk-go-bindings/internal/error"
)

// NewPutItemInput returns a new *dynamodb.PutItemInput
func NewPutItemInput(input interface{}, table string) (*dynamodb.PutItemInput, error) {

	if reflect.DeepEqual(input, reflect.Zero(reflect.TypeOf(input)).Interface()) {
		return nil, intError.Format(ErrEmptyParameter, Input)
	}
	if table == "" {
		return nil, intError.Format(ErrEmptyParameter, Table)
	}

	dynamoInput, err := dynamodbattribute.MarshalMap(input)
	if err != nil {
		return nil, err
	}

	out := &dynamodb.PutItemInput{}
	out = out.SetItem(dynamoInput)
	out = out.SetTableName(table)

	return out, nil

}

// NewGetItemInput returns a new *GetItemInput
func NewGetItemInput(table, keyName, keyValue string) (*dynamodb.GetItemInput, error) {

	if table == "" {
		return nil, intError.Format(ErrEmptyParameter, Table)
	}
	if keyName == "" {
		return nil, intError.Format(ErrEmptyParameter, KeyName)
	}
	if keyValue == "" {
		return nil, intError.Format(ErrEmptyParameter, KeyValue)
	}

	out := &dynamodb.GetItemInput{}
	out = out.SetTableName(table)
	out = out.SetKey(
		map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(keyValue),
			},
		},
	)

	return out, nil

}

func NewScanInput(table, keyName string, keyValue interface{}) (*dynamodb.ScanInput, error) {

	if table == "" {
		return nil, intError.Format(ErrEmptyParameter, table)
	}
	if keyName == "" {
		return nil, intError.Format(ErrEmptyParameter, KeyName)
	}
	if keyValue == nil {
		return nil, intError.Format(ErrEmptyParameter, KeyValue)
	}

	filter := expression.Name(KeyName).Equal(expression.Value(KeyValue))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return nil, err
	}

	return &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(table),
	}, nil

}

// UnmarshalStreamImage unmarshals a events.DynamoDBEventRecord in a pointer to an interface
func UnmarshalStreamImage(input events.DynamoDBEventRecord, output interface{}) error {

	img := input.Change.NewImage

	if reflect.DeepEqual(input, reflect.Zero(reflect.TypeOf(input)).Interface()) {
		return intError.Format(ErrEmptyParameter, Input)
	}

	if reflect.ValueOf(output).Kind() != reflect.Ptr {
		return intError.Format(ErrNoPointerParameter, Output)
	}

	if len(img) == 0 {
		return intError.Format(ErrNoPointerParameter, NewImage)
	}

	dbAttrMap := make(map[string]*dynamodb.AttributeValue)

	for k, v := range img {

		bytes, err := v.MarshalJSON()
		if err != nil {
			return err
		}

		var dbAttr dynamodb.AttributeValue

		json.Unmarshal(bytes, &dbAttr)
		dbAttrMap[k] = &dbAttr

	}

	return dynamodbattribute.UnmarshalMap(dbAttrMap, output)

}

// UnmarshalGetItemOutput unmarshals a *dynamodb.GetItemOutput into a passed interface reference
func UnmarshalGetItemOutput(input *dynamodb.GetItemOutput, out interface{}) error {

	if reflect.ValueOf(out).Kind() != reflect.Ptr {
		return intError.Format(ErrNoPointerParameter, Input)
	}

	unmarshalError := dynamodbattribute.UnmarshalMap(input.Item, out)
	if unmarshalError != nil {
		return unmarshalError
	}

	return nil

}
