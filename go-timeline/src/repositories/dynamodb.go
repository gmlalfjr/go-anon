package repositories

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"time"

	"os"
)

func CreateTimeline(item map[string]*dynamodb.AttributeValue) error{
	dynamodbSession := Session()
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(os.Getenv("TIMELINE_TABLE_NAME")),
	}

	_, err := dynamodbSession.PutItem(input)
	if err != nil  {
		return err
	}
	return nil
}


func GetDetail(id string) (*dynamodb.GetItemOutput, error) {
	dynamodbSession := Session()
	params := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TIMELINE_TABLE_NAME")),

		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	}
	get, err := dynamodbSession.GetItem(params)

	if err != nil {
		return nil, err
	}

	return get, nil
}

func GetLastUserPost(username string) (*dynamodb.QueryOutput, error) {
	dynamodbSession := Session()

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

	query, err := dynamodbSession.Query(params)
	if err != nil {
		return nil, err
	}
	return query, nil
}

func GetAllTimeline(limit int64, key map[string] string) (*dynamodb.QueryOutput, error){
	dynamodbSession := Session()

	dates := time.Now().UTC().String()
	params := &dynamodb.QueryInput{
		TableName:        aws.String(os.Getenv("TIMELINE_TABLE_NAME")),
		Limit:            aws.Int64(int64(limit)),
		ScanIndexForward: aws.Bool(false),
	}
	params.KeyConditionExpression = aws.String("#status = :status and #createdAt > :createdAt")
	params.IndexName = aws.String("statusAndCreatedAtGSI")
	params.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
		":status":    {S: aws.String("OK")},
		":createdAt": {S: aws.String(dates)},
	}
	params.ExpressionAttributeNames = map[string]*string{
		"#status":    aws.String("status"),
		"#createdAt": aws.String("createdAt"),
	}

	if key != nil && key["id"] != "" {
		params.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id":        {S: aws.String(key["id"])},
			"status":    {S: aws.String(key["status"])},
			"createdAt": {S: aws.String(key["createdAt"])},
		}
	}
	query, err := dynamodbSession.Query(params)
	if err != nil {
		return nil, err
	}

	return query, nil
}

func GetAllTimelineByType(limit int64, postType string, key map[string] string) (*dynamodb.QueryOutput, error){
	dynamodbSession := Session()

	dates := time.Now().UTC().String()
	params := &dynamodb.QueryInput{
		TableName:        aws.String(os.Getenv("TIMELINE_TABLE_NAME")),
		Limit:            aws.Int64(int64(limit)),
		ScanIndexForward: aws.Bool(false),
	}
	params.KeyConditionExpression = aws.String("#type = :type and #createdAt > :createdAt")
	params.IndexName = aws.String("typeAndCreatedAtGSI")
	params.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
		":type":      {S: aws.String(postType)},
		":createdAt": {S: aws.String(dates)},
	}
	params.ExpressionAttributeNames = map[string]*string{
		"#type":      aws.String("type"),
		"#createdAt": aws.String("createdAt"),
	}
	if key != nil && key["id"] != "" {
		params.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id":        {S: aws.String(key["id"])},
			"type":      {S: aws.String(key["type"])},
			"createdAt": {S: aws.String(key["createdAt"])},
		}
	}
	query, err := dynamodbSession.Query(params)
	if err != nil {
		return nil, err
	}

	return query, nil
}

func GetAllById(username string, limit int64, key map[string] string)(*dynamodb.QueryOutput, error) {
	dynamodbSession := Session()
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
		Limit:            aws.Int64(int64(limit)),
		IndexName:        aws.String("usernameAndCreatedAtGSI"),
		ScanIndexForward: aws.Bool(false),
	}

	if key["id"] != "" {
		params.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id":   {S: aws.String(key["id"])},
			"createdAt":   {S: aws.String(key["createdAt"])},
			"username": {S: aws.String(key["username"])},
		}
	}

	query, err := dynamodbSession.Query(params)
	if err != nil {
		return nil, err
	}

	return query, nil

}

func Delete(id string) error {
	dynamodbSession := Session()
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv("TIMELINE_TABLE_NAME")),

		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	}

	_, err := dynamodbSession.DeleteItem(params)

	if err != nil {
		return err
	}
	return nil
}

func Session() *dynamodb.DynamoDB{
	sessionOptions := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamoSession := dynamodb.New(sessionOptions)

	return dynamoSession
}