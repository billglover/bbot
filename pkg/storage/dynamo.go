package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

// Save stores a record in DynamoDB. It takes a region and table name. It
// returns an error if unable to save the record.
func Save(r, t string, v interface{}) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(r)},
	)
	if err != nil {
		return errors.Wrap(err, "unable to open session")
	}

	ddb := dynamodb.New(sess)

	value, err := dynamodbattribute.MarshalMap(v)
	if err != nil {
		return errors.Wrap(err, "unable to marshal value")
	}

	record := &dynamodb.PutItemInput{
		Item:      value,
		TableName: aws.String(t),
	}

	if _, err := ddb.PutItem(record); err != nil {
		return errors.Wrap(err, "unable to save record")
	}

	return nil
}

// Retrieve returns a record from DynamoDb. It takes a region, table name, key,
// an ID, and an interface. It returns an error if unable to retrieve the value.
func Retrieve(r, t, k, id string, v interface{}) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(r)},
	)
	if err != nil {
		return errors.Wrap(err, "unable to open session")
	}

	ddb := dynamodb.New(sess)

	request := &dynamodb.GetItemInput{
		TableName: aws.String(t),
		Key:       map[string]*dynamodb.AttributeValue{k: {S: aws.String(id)}},
	}

	record, err := ddb.GetItem(request)
	if err != nil {
		return errors.Wrap(err, "unable to retrieve record")
	}

	if err := dynamodbattribute.UnmarshalMap(record.Item, v); err != nil {
		return errors.Wrap(err, "unable to unmarshal value")
	}

	return nil
}
