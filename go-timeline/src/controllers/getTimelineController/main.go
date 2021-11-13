package main

import (
	"errors"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gmlalfjr/go-anon/go-timeline/src/domain"
	"github.com/gmlalfjr/go-anon/go-timeline/src/services"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
)

func validate(lastEvaluatedKey string, typeKey string) error {

	if lastEvaluatedKey == "" || typeKey == "" {
		msg := errors.New("you must put all key for pagination")
		return msg
	}
	return nil
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	exlusive := &domain.ExlusiveStartKey{}
	limit := request.QueryStringParameters["limit"]
	postType := request.QueryStringParameters["type"]
	lastEvaluatedKey := request.QueryStringParameters["lastEvaluatedKey"]

	typeKey := request.QueryStringParameters["type"]
	if lastEvaluatedKey != "" {
		errValidate := validate(lastEvaluatedKey, typeKey)
		if errValidate != nil {
			return response.FailResponse(&response.ErrorWrapper{
				Error:      errValidate.Error(),
				Message:    "Validation Error",
				StatusCode: 400,
			})
		}

		exlusive = &domain.ExlusiveStartKey{
			Id:   lastEvaluatedKey,
			Type: typeKey,
		}
	}

	if limit == "" {
		limit = "10"
	}
	if postType == "" {
		postType = "ALL"
	}
	convertLimitDataType, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return response.FailResponse(&response.ErrorWrapper{
			Error:      err.Error(),
			Message:    "Failed conv query string",
			StatusCode: 400,
		})
	}
	timeline := &domain.Timeline{
		Type: postType,
	}

	res, pagination, errGetTimeline := services.GetTimeline(timeline, convertLimitDataType, exlusive)

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
