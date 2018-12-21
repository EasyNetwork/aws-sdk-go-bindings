package dynamodb

// GetItemOutput embeds *dynamodb.GetItemOutput
type GetItemOutput interface{}

// ScanOutput embeds *dynamodb.ScanOutput
type ScanOutput interface{}

// DynamoPutItem puts a given input in a dynamodb table
func (svc *DynamoDB) DynamoPutItem(input interface{}, table string) error {

	newPutItemIn, err := NewPutItemInput(input, table)
	if err != nil {
		return err
	}

	_, err = svc.PutItem(newPutItemIn)
	if err != nil {
		return err
	}

	return nil

}

// DynamoGetItem gets an item from DynamoDB given a key and its value.
// A *GetItemOutput will be returned
func (svc *DynamoDB) DynamoGetItem(table, keyName, keyValue string) (*GetItemOutput, error) {

	in, err := NewGetItemInput(
		table,
		keyName,
		keyValue,
	)
	if err != nil {
		return nil, err
	}

	item, err := svc.GetItem(in)
	if err != nil {
		return nil, err
	}

	out := new(GetItemOutput)

	err = UnmarshalGetItemOutput(item, &out)
	if err != nil {
		return nil, err
	}

	return out, nil

}

// DynamoScan gets items from DynamoDB given a key and its value.
// A *ScanOutput will be returned
func (svc *DynamoDB) DynamoScan(table, keyName string, keyValue interface{}) (*ScanOutput, error) {

	scanInput, err := NewScanInput(table, keyName, keyValue)
	if err != nil {
		return nil, err
	}
	scanOutput, err := svc.Scan(scanInput)
	if err != nil {
		return nil, err
	}

	result := new(ScanOutput)

	err = UnmarshalScanOutput(scanOutput, &result)

	if err != nil {
		return nil, err
	}
	return result, nil

}
