package db

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var _db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

// DB is the service to interact with the database
type DB struct{}

// GetItem gets an item from the database. If nothing is found false is returned
// as the first argument. Otherwise true is returned
func (db *DB) GetItem(key interface{}, tableName string, dst interface{}) (bool, error) {
	_key, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		log.Println(fmt.Sprintf("failed to DynamoDB marshal getItem key, %v", err))
		return false, err
	}
	// Prepare the input for the query.
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       _key,
	}

	// Retrieve the item from DynamoDB. If no matching item is found
	// return nil.
	result, err := _db.GetItem(input)
	if err != nil {
		return false, err
	}
	if result.Item == nil {
		return false, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, dst)
	if err != nil {
		return false, err
	}

	return true, nil
}

// PutItem adds a new record to db
func (db *DB) PutItem(tableName string, item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Println(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}
	_, err = _db.PutItem(input)
	return err
}

// UpdateItem update an specific item in the db
func (db *DB) UpdateItem(tableName string, key interface{}, update interface{}, updateExp string) error {
	_key, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		log.Println(fmt.Sprintf("failed to DynamoDB marshal update key, %v", err))
		return err
	}
	_update, err := dynamodbattribute.MarshalMap(update)
	if err != nil {
		log.Println(fmt.Sprintf("failed to DynamoDB marshal update value, %v", err))
		return err
	}
	input := &dynamodb.UpdateItemInput{
		Key:                       _key,
		TableName:                 aws.String(tableName),
		UpdateExpression:          aws.String(updateExp),
		ExpressionAttributeValues: _update,
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	_, err = _db.UpdateItem(input)
	if err != nil {
		log.Println(fmt.Sprintf("failed to DynamoDB update item, %v", err))
		return err
	}
	return nil
}

func (db *DB) GetItems(tableName string, key interface{}, keyCondExp string, projectionExp string, expAttNames map[string]*string, dst interface{}) error {
	_key, err := dynamodbattribute.MarshalMap(key)
	// Prepare the input for the query.
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: _key,
		KeyConditionExpression:    aws.String(keyCondExp),
		ProjectionExpression:      aws.String(projectionExp),
		ExpressionAttributeNames:  expAttNames,
		TableName:                 aws.String(tableName),
	}
	result, err := _db.Query(input)
	if err != nil {
		fmt.Println("Failed to query table", tableName)
		return err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, dst)
	if err != nil {
		return err
	}
	return nil
}
