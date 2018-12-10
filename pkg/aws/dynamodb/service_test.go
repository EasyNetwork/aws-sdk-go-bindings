package dynamodb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/easynetwork/aws-sdk-go-bindings/internal/configuration"
	"github.com/easynetwork/aws-sdk-go-bindings/testdata"
)

type TestDynamoDBDynamoPutItemType struct {
	SomeParam string `json:"some_param"`
}

func TestDynamoDB_Methods(t *testing.T) {

	cfg := testdata.MockConfiguration(t)

	testDynamoDBDynamoPutItem(t, cfg)
	testDynamoDBDynamoGetItem(t, cfg)

}

func testDynamoDBDynamoPutItem(t *testing.T, cfg *configuration.Configuration) {

	dynamoSvc := testdata.MockDynamoDB(t, cfg)

	tableName := cfg.DynamoDB.PkgTableName

	testdata.MockDynamoDBTable(t, dynamoSvc, tableName, cfg)

	var input TestDynamoDBDynamoPutItemType
	input.SomeParam = cfg.DynamoDB.PrimaryKey

	dynamoNewSvc := &DynamoDB{
		DynamoDB: dynamoSvc,
	}

	err := dynamoNewSvc.DynamoPutItem(input, tableName)

	assert.NoError(t, err)

}

func testDynamoDBDynamoGetItem(t *testing.T, cfg *configuration.Configuration) {

	dynamoSvc := testdata.MockDynamoDB(t, cfg)

	tableName := cfg.DynamoDB.PkgTableName
	primaryKey := cfg.DynamoDB.PrimaryKey
	keyValue := cfg.DynamoDB.PrimaryKey

	testdata.MockDynamoDBTable(t, dynamoSvc, tableName, cfg)

	var input TestDynamoDBDynamoPutItemType
	input.SomeParam = cfg.DynamoDB.PrimaryKey

	dynamoNewSvc := &DynamoDB{
		DynamoDB: dynamoSvc,
	}

	err := dynamoNewSvc.DynamoPutItem(input, tableName)

	assert.NoError(t, err)

	getItemOut, err := dynamoNewSvc.DynamoGetItem(tableName, primaryKey, keyValue)

	assert.NoError(t, err)
	assert.NotEmpty(t, getItemOut)

}
