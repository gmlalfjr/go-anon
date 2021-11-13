package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gmlalfjr/go-anon/go-timeline/src/domain"
	"github.com/gmlalfjr/go-anon/go-timeline/src/services"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	timelineString := request.Body
	timelineStruct := &domain.Timeline{}

	username := request.RequestContext.Authorizer["username"].(string)
	errUnmarshal := json.Unmarshal([]byte(timelineString), &timelineStruct)
	if errUnmarshal != nil {
		return response.FailResponse(&response.ErrorWrapper{
			Error:      errUnmarshal.Error(),
			Message:    "Failed Unmarshal Json",
			StatusCode: 400,
		})
	}

	if errValidate := timelineStruct.Validate(); errValidate != nil {
		return response.FailResponse(&response.ErrorWrapper{
			Error:      errValidate.Error(),
			Message:    "Error Validation Json Schema",
			StatusCode: 400,
		})
	}
	res, postTimelineErr := services.PostTimeline(timelineStruct, username)

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
