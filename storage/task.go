package storage

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/santiagozuluaga/todoserverless/models"
)

const (
	tableName = "task"
)

var (
	dynamoClient *dynamodb.DynamoDB

	putItemWithContext func(ctx aws.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error)
	getItemWithContext func(ctx context.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error)

	ErrTaskNotFound = errors.New("task not found")
)

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dynamoClient = dynamodb.New(sess)
	putItemWithContext = dynamoClient.PutItemWithContext
	getItemWithContext = dynamoClient.GetItemWithContext
}

func InsertTask(ctx context.Context, task models.Task) error {
	item, err := dynamodbattribute.MarshalMap(task)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	_, err = putItemWithContext(ctx, input)

	return err
}

func GetTaskFromID(ctx context.Context, taskID string) (*models.Task, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(taskID),
			},
		},
	}

	output, err := getItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	task := &models.Task{}

	err = dynamodbattribute.UnmarshalMap(output.Item, &task)
	if err != nil {
		return nil, err
	}

	if task.ID == "" {
		return nil, ErrTaskNotFound
	}

	return task, err
}
