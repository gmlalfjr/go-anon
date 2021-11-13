package domain

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
)

const timeDurationPost = -1

func (t *Timeline) PostTimeline() *response.RestErr {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamo := dynamodb.New(sess)
	data, errGetDetail := t.GetDetail(t.Username)
	if errGetDetail != nil {

		log.Println("Error marshalling item: ", errGetDetail.Error)
		return response.Error("Error when marshalling dyanmodb item", 400, errors.New(errGetDetail.Error))
	}
	if data != nil {
		return response.Error("Error data", 400, errors.New("user Cant post at this time"))
	}

	item, err := dynamodbattribute.MarshalMap(&t)
	// newT2 := time.Now().Add(time.Minute * 15)
	// newT2.
	if err != nil {
		log.Println("Error marshalling item: ", err.Error())
		return response.Error("Error when marshalling create dyanmodb item", 400, err)
	}
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(os.Getenv("TIMELINE_TABLE_NAME")),
	}
	_, err = dynamo.PutItem(input)
	if err != nil {
		log.Println("Got error calling PutItem: ", err.Error())
		return response.Error("Error when Insert Item", 500, err)
	}
	return nil
}

func (t *Timeline) GetTimeline(limit int64, key *ExlusiveStartKey) ([]Timeline, *PaginationTimeline, *response.RestErr) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamo := dynamodb.New(sess)

	params := &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("TIMELINE_TABLE_NAME")),
		KeyConditionExpression: aws.String("#type = :type"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":type": {S: aws.String(t.Type)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#type": aws.String("type"),
		},
		Limit:            aws.Int64(int64(limit)),
		IndexName:        aws.String("typeGSI"),
		ScanIndexForward: aws.Bool(false),
	}

	if key.Id != "" {
		params.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(key.Id)},
			// "createdAt": {S: aws.String(key.CreatedAt)},
			"type": {S: aws.String(key.Type)},
		}
	}

	query, err := dynamo.Query(params)
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

func (t *Timeline) GetDetail(username string) (*Timeline, *response.RestErr) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamo := dynamodb.New(sess)
	dates := time.Now().UTC().String()
	params := &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("TIMELINE_TABLE_NAME")),
		KeyConditionExpression: aws.String("#username = :username and #createdAt > :createdAt"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username":  {S: aws.String(username)},
			":createdAt": {S: aws.String(dates)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#username":  aws.String("username"),
			"#createdAt": aws.String("createdAt"),
		},
		Limit:            aws.Int64(1),
		IndexName:        aws.String("usernameAndCreatedAtGSI"),
		ScanIndexForward: aws.Bool(false),
	}
	os.Getenv("SET_TIME_CREATE_TIMELINE")
	query, err := dynamo.Query(params)
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
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamo := dynamodb.New(sess)
	params := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TIMELINE_TABLE_NAME")),

		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	}
	get, err := dynamo.GetItem(params)
	if err != nil {
		return nil, response.Error("Failed Get Item", 500, err)
	}

	if len(get.Item) <= 0 {
		return nil, response.Error("Id Not Found", 404, errors.New("id Not Found"))
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
		return false, response.Error("Got Delete timeline", err.Status, errors.New("not authorize to delete this post"))
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamo := dynamodb.New(sess)
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv("TIMELINE_TABLE_NAME")),

		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	}
	_, errDel := dynamo.DeleteItem(params)
	if errDel != nil {
		return false, response.Error("Failed Delete Item", 500, errDel)
	}

	return true, nil
}
