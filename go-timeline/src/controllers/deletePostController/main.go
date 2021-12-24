package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gmlalfjr/go-anon/go-timeline/src/services"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	username := request.RequestContext.Authorizer["username"].(string)
	id := request.PathParameters["idPostDelete"]

	_, errDeleteUserPost := services.DeleteUserPost(id, username)

	if errDeleteUserPost != nil {
		return response.FailResponse(&response.ErrorWrapper{
			Error:      errDeleteUserPost.Error,
			Message:    errDeleteUserPost.Message,
			StatusCode: errDeleteUserPost.Status,
		})
	}

	return response.SuccessResponse(&response.SuccessWrapper{
		Message: "Success Delete Data",
		StatusCode: 200,
	})
}

func main() {
	lambda.Start(Handler)
}
