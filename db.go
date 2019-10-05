package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

// getItem gets an item from the database. If nothing is found false is returned
// as the first argument. Otherwise true is returned
func getItem(keyName string, keyValue string, tableName string, out interface{}) (bool, error) {
	// Prepare the input for the query.
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(keyValue),
			},
		},
	}

	// Retrieve the item from DynamoDB. If no matching item is found
	// return nil.
	result, err := db.GetItem(input)
	if err != nil {
		return false, err
	}
	if result.Item == nil {
		return false, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, out)
	if err != nil {
		return false, err
	}

	return true, nil
}
