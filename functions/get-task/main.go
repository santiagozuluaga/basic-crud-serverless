package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/santiagozuluaga/todoserverless/shared/apigateway"
	"github.com/santiagozuluaga/todoserverless/shared/logger"
	"github.com/santiagozuluaga/todoserverless/storage"
)

var (
	defaultLogger = logger.NewLogger("get-task")

	errMissingTaskID = errors.New("missing task id")
)

func apiGatewayHandler(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	taskID, ok := request.PathParameters["id"]
	if !ok || taskID == "" {
		defaultLogger.Log(logger.WarningLevel, logger.LogProperties{
			"event": "get_task_id_from_request_failed",
			"error": errMissingTaskID.Error(),
		})

		return apigateway.NewErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid request: %s.", errMissingTaskID.Error())), nil
	}

	task, err := storage.GetTaskFromID(ctx, taskID)
	if err != nil && errors.Is(err, storage.ErrTaskNotFound) {
		defaultLogger.Log(logger.WarningLevel, logger.LogProperties{
			"event": "get_task_from_id_failed",
			"error": err.Error(),
		})

		return apigateway.NewErrorResponse(http.StatusNotFound, "Task not found"), nil
	}

	if err != nil {
		defaultLogger.Log(logger.ErrorLevel, logger.LogProperties{
			"event": "get_task_from_id_failed",
			"error": err.Error(),
		})

		return apigateway.NewErrorResponse(http.StatusInternalServerError, "Internal server error, please try again later."), nil
	}

	defaultLogger.Log(logger.InfoLevel, logger.LogProperties{
		"event": "get_task_succeeded",
		"task":  task,
	})

	return apigateway.NewResponse(http.StatusOK, task), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
