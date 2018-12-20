package sns

// SnsPublish publishes an input on a given SNS targetArn
func (svc *SNS) SnsPublish(input interface{}, messageAttributes map[string]interface{}, targetArn string) (err error) {

	in, err := NewPublishInput(
		input,
		messageAttributes,
		targetArn,
	)
	if err != nil {
		return err
	}

	_, err = svc.SNS.Publish(in)
	if err != nil {
		return err
	}

	return nil

}
