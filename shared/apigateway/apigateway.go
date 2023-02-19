package apigateway

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func NewErrorResponse(statusCode int, message string) *events.APIGatewayProxyResponse {
	bytes, _ := json.Marshal(ErrorMessage{
		Message: message,
	})

	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(bytes),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}

func NewResponse(statusCode int, body interface{}) *events.APIGatewayProxyResponse {
	if body == nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: statusCode,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}
	}

	bytes, _ := json.Marshal(body)

	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(bytes),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}
