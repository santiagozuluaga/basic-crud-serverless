package main

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/santiagozuluaga/todoserverless/models"
	"github.com/santiagozuluaga/todoserverless/shared/logger"
	"github.com/santiagozuluaga/todoserverless/storage"
	"github.com/stretchr/testify/require"
)

func init() {
	storage.InitPutItemWithContextMock()
}

func TestHandler(t *testing.T) {
	c := require.New(t)

	buf := bytes.NewBufferString("")
	logger.InitMock(buf)

	defer func() {
		logger.ClearMock()
	}()

	request := &events.APIGatewayProxyRequest{
		Body: "title=Research how to do a sandwich&description=It's very important know how to do a sandwich",
	}

	response, err := apiGatewayHandler(context.Background(), request)
	c.NoError(err)
	c.NotNil(response)
	c.Equal(response.StatusCode, http.StatusOK)
	c.NotContains(buf.String(), `"level":"error"`)
	c.NotContains(buf.String(), `"level":"fatal"`)
	c.NotContains(buf.String(), `"level":"warning"`)
	c.Contains(buf.String(), `"level":"info"`)
	c.Contains(buf.String(), `"event":"create_new_task_succeeded"`)
}

func TestHandlerForceStorageFailure(t *testing.T) {
	c := require.New(t)

	buf := bytes.NewBufferString("")
	logger.InitMock(buf)
	storage.ActivePutItemWithContextMockFailure()

	defer func() {
		logger.ClearMock()
		storage.DeactivePutItemWithContextMockFailure()
	}()

	request := &events.APIGatewayProxyRequest{
		Body: "title=Research how to do a sandwich&description=It's very important know how to do a sandwich",
	}

	response, err := apiGatewayHandler(context.Background(), request)
	c.NoError(err)
	c.NotNil(response)
	c.Equal(response.StatusCode, http.StatusInternalServerError)
	c.Contains(response.Body, "Internal server error, please try again later.")

	c.NotContains(buf.String(), `"level":"info"`)
	c.NotContains(buf.String(), `"level":"fatal"`)
	c.NotContains(buf.String(), `"level":"warning"`)
	c.Contains(buf.String(), `"level":"error"`)
	c.Contains(buf.String(), storage.ErrForceMockFailure.Error())
	c.Contains(buf.String(), `"event":"insert_new_task_in_dynamo_failed"`)
}

func TestHandlerErrMissingTitle(t *testing.T) {
	c := require.New(t)

	buf := bytes.NewBufferString("")
	logger.InitMock(buf)

	defer func() {
		logger.ClearMock()
	}()

	request := &events.APIGatewayProxyRequest{
		Body: "description=It's very important know how to do a sandwich",
	}

	response, err := apiGatewayHandler(context.Background(), request)
	c.NoError(err)
	c.NotNil(response)
	c.Equal(response.StatusCode, http.StatusBadRequest)
	c.Contains(response.Body, models.ErrMissingTitle.Error())

	c.NotContains(buf.String(), `"level":"info"`)
	c.NotContains(buf.String(), `"level":"fatal"`)
	c.NotContains(buf.String(), `"level":"warning"`)
	c.Contains(buf.String(), `"level":"error"`)
	c.Contains(buf.String(), `"event":"create_new_task_failed"`)
}

func TestHandlerErrMissingDescription(t *testing.T) {
	c := require.New(t)

	buf := bytes.NewBufferString("")
	logger.InitMock(buf)

	defer func() {
		logger.ClearMock()
	}()

	request := &events.APIGatewayProxyRequest{
		Body: "title=Research how to do a sandwich",
	}

	response, err := apiGatewayHandler(context.Background(), request)
	c.NoError(err)
	c.NotNil(response)
	c.Equal(response.StatusCode, http.StatusBadRequest)
	c.Contains(response.Body, models.ErrMissingDescription.Error())

	c.NotContains(buf.String(), `"level":"info"`)
	c.NotContains(buf.String(), `"level":"fatal"`)
	c.NotContains(buf.String(), `"level":"warning"`)
	c.Contains(buf.String(), `"level":"error"`)
	c.Contains(buf.String(), `"event":"create_new_task_failed"`)
}
