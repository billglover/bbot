package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

// DynamoDB represents a DynamoDB table.
type DynamoDB struct {
	Region string
	Table  string
}

// Save stores a record in DynamoDB. It takes an interface and returns an error
// if unable to save the record.
func (d *DynamoDB) Save(v interface{}) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(d.Region)},
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
		TableName: aws.String(d.Table),
	}

	if _, err := ddb.PutItem(record); err != nil {
		return errors.Wrap(err, "unable to save record")
	}

	return nil
}

// Retrieve returns a record from DynamoDb. It takes a region, table name, key,
// an ID, and an interface. It returns an error if unable to retrieve the value.
func (d *DynamoDB) Retrieve(k, id string, v interface{}) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(d.Region)},
	)
	if err != nil {
		return errors.Wrap(err, "unable to open session")
	}

	ddb := dynamodb.New(sess)

	request := &dynamodb.GetItemInput{
		TableName: aws.String(d.Table),
		Key:       map[string]*dynamodb.AttributeValue{k: {S: aws.String(id)}},
	}

	record, err := ddb.GetItem(request)
	if err != nil {
		return errors.Wrap(err, "unable to retrieve record")
	}

	if len(record.Item) == 0 {
		return errors.New("no auth token exists for team: " + id)
	}

	if err := dynamodbattribute.UnmarshalMap(record.Item, v); err != nil {
		return errors.Wrap(err, "unable to unmarshal value")
	}

	return nil
}
