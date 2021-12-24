package domain

import (
	"errors"
	"github.com/gmlalfjr/go-anon/go-timeline/src/repositories"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
)

const timeDurationPost = -1

func (t *Timeline) PostTimeline() *response.RestErr {
	data, errGetDetail := t.GetDetailLasUserPost(t.Username)
	if errGetDetail != nil {
		log.Println("Error marshalling item: ", errGetDetail.Error)
		return response.Error("Error when marshalling dyanmodb item", 400, errors.New(errGetDetail.Error))
	}
	if data != nil {
		return response.Error("Error data", 400, errors.New("user Cant post at this time"))
	}

	item, err := dynamodbattribute.MarshalMap(&t)

	if err != nil {
		log.Println("Error marshalling item: ", err.Error())
		return response.Error("Error when marshalling create dyanmodb item", 400, err)
	}
	err = repositories.CreateTimeline(item)
	if err != nil {
		log.Println("Got error calling PutItem: ", err.Error())
		return response.Error("Error when Insert Item", 500, err)
	}
	return nil
}

func (t *Timeline) GetTimeline(limit int64, key map[string] string) ([]Timeline, *PaginationTimeline, *response.RestErr) {
	var query *dynamodb.QueryOutput
	var err error
	if t.Type == "ALL" {
		query, err =repositories.GetAllTimeline(limit, key)
		if err != nil {
			log.Println("Got Error When Get All Timeline Data", err)
			return nil, nil, response.Error("Got Error When Get All Timeline Data", 500, err)
		}

	} else {
		query, err =repositories.GetAllTimelineByType(limit,t.Type, key)
		if err != nil {
			log.Println("Got Error When Get All Timeline Data", err)
			return nil, nil, response.Error("Got Error When Get All Timeline Data", 500, err)
		}

	}
	if err != nil {
		return nil, nil, response.Error("Failed Query List Timeline", 500, err)
	}

	var results []Timeline
	for _, i := range query.Items {
		timeline := Timeline{}

		err = dynamodbattribute.UnmarshalMap(i, &timeline)
		if err != nil {
			log.Println("Got error unmarshalling", err)
			return nil, nil, response.Error("Got error unmarshallin", 500, err)
		}

		results = append(results, timeline)
	}
	pagination := PaginationTimeline{}
	err = dynamodbattribute.UnmarshalMap(query.LastEvaluatedKey, &pagination)
	if err != nil {
		log.Println("Got error unmarshalling", err)
		return nil, nil, response.Error("Got error unmarshallin", 500, err)
	}
	if pagination.Id == "" && pagination.Type == "" {
		return results, nil, nil
	}

	return results, &pagination, nil
}

func (t *Timeline) GetDetailLasUserPost(username string) (*Timeline, *response.RestErr) {
	query, err := repositories.GetLastUserPost(username)
	if err != nil {
		return nil, response.Error("Failed Query Get Last Item", 500, err)
	}

	timeline := Timeline{}

	if len(query.Items) > 0 {
		err = dynamodbattribute.UnmarshalMap(query.Items[0], &timeline)
		if err != nil {
			log.Println("Got error unmarshalling", err)
			return nil, response.Error("Got error unmarshallin", 500, err)
		}
		dateNow := time.Now()

		dates := dateNow.Add(timeDurationPost * time.Minute)

		checkLastTimeUserPost := dates.After(timeline.CreatedAt)

		if !checkLastTimeUserPost {
			return &timeline, nil
		}
		return nil, nil
	}
	return nil, nil

}

func (t *Timeline) GetTimelineDetail(id string) (*Timeline, *response.RestErr) {
	get, err := repositories.GetDetail(id)
	if err != nil {
		return nil, response.Error("Failed Get Item", 500, err)
	}

	if len(get.Item) <= 0 {
		return nil, response.Error("Not Found Error", 404, errors.New("post not found"))
	}
	timeline := &Timeline{}
	err = dynamodbattribute.UnmarshalMap(get.Item, &timeline)
	if err != nil {
		log.Println("Got error unmarshalling", err)
		return nil, response.Error("Got error unmarshallin", 500, err)
	}
	return timeline, nil
}

func (t *Timeline) DeleteUserPost(id string, username string) (bool, *response.RestErr) {
	data, err := t.GetTimelineDetail(id)
	if err != nil {
		return false, response.Error("not found id", err.Status, errors.New(err.Error))
	}
	if data.Username != username {
		return false, response.Error("Cant delete this post", 404, errors.New("not authorize to delete this post"))
	}
	errDelete := repositories.Delete(id)
	if errDelete != nil {
		return false, response.Error("Failed Delete Item", 500,errDelete)
	}

	return true, nil
}

func (t *Timeline) GetOwnUserPost(username string, limit int64, key map[string] string) ([]Timeline, *PaginationTimelineById, *response.RestErr) {
	query, err :=repositories.GetAllById(username, limit, key)

	if err != nil {
		return nil, nil, response.Error("Failed Query Own List Timeline", 500, err)
	}

	var results []Timeline
	for _, i := range query.Items {
		timeline := Timeline{}

		err = dynamodbattribute.UnmarshalMap(i, &timeline)
		if err != nil {
			log.Println("Got error unmarshalling", err)
			return nil, nil, response.Error("Got error unmarshallin", 500, err)
		}

		results = append(results, timeline)
	}
	pagination := PaginationTimelineById{}
	err = dynamodbattribute.UnmarshalMap(query.LastEvaluatedKey, &pagination)
	if err != nil {
		log.Println("Got error unmarshalling", err)
		return nil, nil, response.Error("Got error unmarshalling", 500, err)
	}
	if pagination.Id == "" && pagination.Username == "" {
		return results, nil, nil
	}

	return results, &pagination, nil
}
