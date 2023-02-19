package storage

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/santiagozuluaga/todoserverless/models"
)

var (
	mockedTasks = map[string]models.Task{}

	backupPutItemWithContext func(ctx aws.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error)
	backupGetItemWithContext func(ctx context.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error)

	ErrForceMockFailure = errors.New("force mock failure")
)

// PutItemWithContext
func putItemWithContextMock(ctx aws.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
	newTasks := models.Task{}

	err := dynamodbattribute.UnmarshalMap(input.Item, &newTasks)
	if err != nil {
		return nil, err
	}

	mockedTasks[newTasks.ID] = newTasks

	return nil, nil
}

func InitPutItemWithContextMock() {
	backupPutItemWithContext = putItemWithContext

	putItemWithContext = putItemWithContextMock
}

func ClearPutItemWithContextMock() {
	putItemWithContext = backupPutItemWithContext
}

func ActivePutItemWithContextMockFailure() {
	backupPutItemWithContext = putItemWithContext

	putItemWithContext = func(ctx aws.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
		return nil, ErrForceMockFailure
	}
}

func DeactivePutItemWithContextMockFailure() {
	backupPutItemWithContext = putItemWithContext

	putItemWithContext = putItemWithContextMock
}

// GetItemWithContext
func getItemWithContextMock(ctx context.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
	task, ok := mockedTasks[*input.Key["id"].S]
	if !ok {
		item, err := dynamodbattribute.MarshalMap(models.Task{})
		if err != nil {
			return nil, err
		}

		return &dynamodb.GetItemOutput{
			Item: item,
		}, nil
	}

	item, err := dynamodbattribute.MarshalMap(task)
	if err != nil {
		return nil, err
	}

	return &dynamodb.GetItemOutput{
		Item: item,
	}, nil
}

func InitGetItemWithContextMock() {
	backupGetItemWithContext = getItemWithContext

	getItemWithContext = getItemWithContextMock
}

func ClearGetItemWithContextMock() {
	getItemWithContext = backupGetItemWithContext
}

func ActiveGetItemWithContextMockFailure() {
	backupGetItemWithContext = getItemWithContext

	getItemWithContext = func(ctx context.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
		return nil, ErrForceMockFailure
	}
}

func DeactiveGetItemWithContextMockFailure() {
	backupGetItemWithContext = getItemWithContext

	getItemWithContext = getItemWithContextMock
}
