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
	storage.InitGetItemWithContextMock()
}

func TestHandler(t *testing.T) {
	c := require.New(t)

	buf := bytes.NewBufferString("")
	logger.InitMock(buf)

	defer func() {
		logger.ClearMock()
	}()

	storage.InsertTask(context.Background(), models.Task{
		ID:          "dummy-task-id",
		Title:       "dummy task title",
		Description: "dummy task description",
	})

	request := &events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"id": "dummy-task-id",
		},
	}

	response, err := apiGatewayHandler(context.Background(), request)
	c.NoError(err)
	c.NotNil(response)
	c.Equal(response.StatusCode, http.StatusOK)
	c.NotContains(buf.String(), `"level":"error"`)
	c.NotContains(buf.String(), `"level":"fatal"`)
	c.NotContains(buf.String(), `"level":"warning"`)
	c.Contains(buf.String(), `"level":"info"`)
	c.Contains(buf.String(), `"event":"get_task_succeeded"`)
}

func TestHandlerErrTaskNotFound(t *testing.T) {
	c := require.New(t)

	buf := bytes.NewBufferString("")
	logger.InitMock(buf)

	defer func() {
		logger.ClearMock()
	}()

	request := &events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"id": "dummy-task-id-1",
		},
	}

	response, err := apiGatewayHandler(context.Background(), request)
	c.NoError(err)
	c.NotNil(response)
	c.Equal(response.StatusCode, http.StatusNotFound)
	c.Contains(response.Body, "Task not found")

	c.NotContains(buf.String(), `"level":"info"`)
	c.NotContains(buf.String(), `"level":"fatal"`)
	c.NotContains(buf.String(), `"level":"error"`)
	c.Contains(buf.String(), `"level":"warning"`)
	c.Contains(buf.String(), storage.ErrTaskNotFound.Error())
	c.Contains(buf.String(), `"event":"get_task_from_id_failed"`)
}

func TestHandlerForceStorageFailure(t *testing.T) {
	c := require.New(t)

	buf := bytes.NewBufferString("")
	logger.InitMock(buf)
	storage.ActiveGetItemWithContextMockFailure()

	defer func() {
		logger.ClearMock()
		storage.DeactiveGetItemWithContextMockFailure()
	}()

	request := &events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"id": "dummy-task-id-1",
		},
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
	c.Contains(buf.String(), `"event":"get_task_from_id_failed"`)
}

func TestHandlerErrMissingTaskID(t *testing.T) {
	c := require.New(t)

	buf := bytes.NewBufferString("")
	logger.InitMock(buf)

	defer func() {
		logger.ClearMock()
	}()

	request := &events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"id": "",
		},
	}

	response, err := apiGatewayHandler(context.Background(), request)
	c.NoError(err)
	c.NotNil(response)
	c.Equal(response.StatusCode, http.StatusBadRequest)
	c.Contains(response.Body, "Invalid request: missing task id.")

	c.NotContains(buf.String(), `"level":"info"`)
	c.NotContains(buf.String(), `"level":"fatal"`)
	c.NotContains(buf.String(), `"level":"error"`)
	c.Contains(buf.String(), `"level":"warning"`)
	c.Contains(buf.String(), errMissingTaskID.Error())
	c.Contains(buf.String(), `"event":"get_task_id_from_request_failed"`)
}
