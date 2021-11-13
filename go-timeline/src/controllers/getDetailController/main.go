package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gmlalfjr/go-anon/go-timeline/src/services"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	id := request.PathParameters["idDetail"]

	res, errGetTimeline := services.GetTimelineDetail(id)

	if errGetTimeline != nil {
		return response.FailResponse(&response.ErrorWrapper{
			Error:      errGetTimeline.Error,
			Message:    errGetTimeline.Message,
			StatusCode: errGetTimeline.Status,
		})
	}

	return response.SuccessResponse(&response.SuccessWrapper{
		Message:    "Success Get Data",
		Data:       res,
		StatusCode: 200,
	})
}

func main() {
	lambda.Start(Handler)
}
