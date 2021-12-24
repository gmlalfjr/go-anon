package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/gmlalfjr/go-anon/go-timeline/src/services"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	timelineString := request.Body


	username := request.RequestContext.Authorizer["username"].(string)


	res, postTimelineErr := services.PostTimeline(username, timelineString)

	if postTimelineErr != nil {
		return response.FailResponse(&response.ErrorWrapper{
			Error:      postTimelineErr.Error,
			Message:    postTimelineErr.Message,
			StatusCode: postTimelineErr.Status,
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
