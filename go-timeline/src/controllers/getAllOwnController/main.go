package main

import (
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gmlalfjr/go-anon/go-timeline/src/services"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	username := request.RequestContext.Authorizer["username"].(string)
	limit := request.QueryStringParameters["limit"]
	lastEvaluatedKey := request.QueryStringParameters["lastEvaluatedKey"]
	createdAtKey := request.QueryStringParameters["createdAt"]
	var key map[string] string
	if lastEvaluatedKey != "" {
		key = map[string] string {
			"id": lastEvaluatedKey,
			"createdAt": createdAtKey,
			"username": username,
		}
	}
	if limit == ""{
		limit = "10"
	}
	convertLimitDataType, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return response.FailResponse(&response.ErrorWrapper{
			Error:      err.Error(),
			Message:    "Failed conv query string",
			StatusCode: 400,
		})
	}

	res, pagination, errGetTimeline := services.GetOwnPost(username, convertLimitDataType, key)

	if errGetTimeline != nil {
		return response.FailResponse(&response.ErrorWrapper{
			Error:      errGetTimeline.Error,
			Message:    errGetTimeline.Message,
			StatusCode: errGetTimeline.Status,
		})
	}

	return response.SuccessResponsePagination(&response.SuccessPaginationWrapper{
		Message:          "Success Get All Post Data",
		LastEvaluatedKey: pagination,
		Data:             res,
		StatusCode:       200,
	})
}

func main() {
	lambda.Start(Handler)
}
