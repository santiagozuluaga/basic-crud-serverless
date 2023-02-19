package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/santiagozuluaga/todoserverless/models"
	"github.com/santiagozuluaga/todoserverless/shared/apigateway"
	"github.com/santiagozuluaga/todoserverless/shared/logger"
	"github.com/santiagozuluaga/todoserverless/storage"
)

var (
	defaultLogger = logger.NewLogger("create-task")
)

func apiGatewayHandler(ctx context.Context, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	body, err := url.ParseQuery(request.Body)
	if err != nil {
		defaultLogger.Log(logger.ErrorLevel, logger.LogProperties{
			"event":        "parse_request_body_failed",
			"error":        err.Error(),
			"request_body": request.Body,
		})

		return apigateway.NewErrorResponse(http.StatusBadRequest, "Invalid request: unexpected request body."), nil
	}

	task, err := models.NewTask(body.Get("title"), body.Get("description"))
	if err != nil {
		defaultLogger.Log(logger.ErrorLevel, logger.LogProperties{
			"event":        "create_new_task_failed",
			"error":        err.Error(),
			"request_body": request.Body,
		})

		return apigateway.NewErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid request: %s.", err.Error())), nil
	}

	err = storage.InsertTask(ctx, *task)
	if err != nil {
		defaultLogger.Log(logger.ErrorLevel, logger.LogProperties{
			"event":        "insert_new_task_in_dynamo_failed",
			"error":        err.Error(),
			"request_body": task,
		})

		return apigateway.NewErrorResponse(http.StatusInternalServerError, "Internal server error, please try again later."), nil
	}

	defaultLogger.Log(logger.InfoLevel, logger.LogProperties{
		"event": "create_new_task_succeeded",
		"task":  task,
	})

	return apigateway.NewResponse(http.StatusOK, task), nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
