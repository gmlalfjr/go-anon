package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gmlalfjr/go-anon/go-comments/src/domain"
	"github.com/gmlalfjr/go-anon/go-comments/src/services"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	exlusive := &domain.ExlusiveStartKey{}
	postId := request.PathParameters["postId"]
	limit := request.QueryStringParameters["limit"]
	lastEvaluatedKey := request.QueryStringParameters["lastEvaluatedKey"]

	if lastEvaluatedKey != "" {

		exlusive = &domain.ExlusiveStartKey{
			Id: lastEvaluatedKey,
		}
	}

	res, pag, err := services.GetComments(postId, limit, exlusive)

	if err != nil {
		return response.FailResponse(&response.ErrorWrapper{
			Error:      err.Error,
			Message:    err.Message,
			StatusCode: err.Status,
		})
	}
	return response.SuccessResponsePagination(&response.SuccessPaginationWrapper{
		Message:          "Success Get Data",
		Data:             res,
		StatusCode:       200,
		LastEvaluatedKey: pag,
	})
}

func main() {
	lambda.Start(Handler)
}
