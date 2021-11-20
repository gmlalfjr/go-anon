package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gmlalfjr/go-anon/go-comments/src/domain"
	"github.com/gmlalfjr/go-anon/go-comments/src/services"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	commentString := request.Body
	commentStruct := &domain.Comments{}
	errUnmarshal := json.Unmarshal([]byte(commentString), &commentStruct)
	if errUnmarshal != nil {
		return response.FailResponse(&response.ErrorWrapper{
			Error:      errUnmarshal.Error(),
			Message:    "Failed Unmarshal Json",
			StatusCode: 400,
		})
	}
	username := request.RequestContext.Authorizer["username"].(string)
	postId := request.PathParameters["postId"]
	res, err := services.CreateComment(commentStruct, username, postId)

	if err != nil {
		return response.FailResponse(&response.ErrorWrapper{
			Error:      err.Error,
			Message:    err.Message,
			StatusCode: err.Status,
		})
	}
	return response.SuccessResponse(&response.SuccessWrapper{
		Message:    "Success Create Data",
		Data:       res,
		StatusCode: 202,
	})
}

func main() {
	lambda.Start(Handler)
}
