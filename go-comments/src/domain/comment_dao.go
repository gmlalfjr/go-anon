package domain

import (
	"errors"
	"log"

	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gmlalfjr/go-anon/go-comments/src/constant"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
)

func dynamoSession() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return dynamodb.New(sess)
}

func (c *Comments) CreateComment(username string, postId string) *response.RestErr {
	dynamo := dynamoSession()
	errGet := c.GetTimeline(postId)
	if errGet != nil {
		return response.Error(errGet.Message, errGet.Status, errors.New(errGet.Error))
	}

	errGetComment := c.GetDetailCommentToGetGeneratedName(postId)
	if errGetComment != nil {
		return response.Error(errGetComment.Message, errGetComment.Status, errors.New(errGetComment.Error))
	}
	item, err := dynamodbattribute.MarshalMap(&c)
	if err != nil {
		log.Println("Error marshalling item: ", err.Error())
		return response.Error("Error when marshalling create dyanmodb item", 400, err)
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(os.Getenv("COMMENT_TABLE_NAME")),
	}
	_, errPut := dynamo.PutItem(input)
	if errPut != nil {
		log.Println("Got error calling PutItem: ", errPut.Error())
		return response.Error("Error when Insert Item", 500, errPut)
	}

	return nil
}

func (c *Comments) GetTimeline(postId string) *response.RestErr {
	dynamo := dynamoSession()
	timeline := &Timeline{}
	params := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TIMELINE_TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(postId)},
		},
	}

	get, err := dynamo.GetItem(params)
	if err != nil {
		return response.Error("Failed Get Item", 500, err)
	}
	if len(get.Item) <= 0 {
		return response.Error("Id Not Found", 404, errors.New("id Not Found"))
	}
	err = dynamodbattribute.UnmarshalMap(get.Item, timeline)
	if err != nil {
		log.Println("Got error unmarshalling", err)
		return response.Error("Got error unmarshallin", 500, err)
	}
	if timeline.IsPrivate {
		return response.Error("This Post is private", 403, errors.New("cant comment on this post"))
	}
	return nil
}

func (c *Comments) GetDetailCommentToGetGeneratedName(postId string) *response.RestErr {
	dynamo := dynamoSession()

	params := &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("COMMENT_TABLE_NAME")),
		KeyConditionExpression: aws.String("#postId = :postId and #username = :username"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {S: aws.String(c.Username)},
			":postId":   {S: aws.String(postId)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#username": aws.String("username"),
			"#postId":   aws.String("postId"),
		},
		Limit:            aws.Int64(int64(1)),
		IndexName:        aws.String("postIdAndUsernameGSI"),
		ScanIndexForward: aws.Bool(false),
	}
	data, err := dynamo.Query(params)

	if err != nil {
		return response.Error("Failed Query Comments", 500, err)
	}

	if len(data.Items) > 0 {
		comment := Comments{}
		err = dynamodbattribute.UnmarshalMap(data.Items[0], &comment)
		if err != nil {
			log.Println("Got error unmarshalling", err)
			return response.Error("Got error unmarshallin", 500, err)
		}
		c.GeneratedName = comment.GeneratedName
	} else {
		generateName := constant.GenerateName()
		c.GeneratedName = generateName
	}

	return nil
}

func (c *Comments) GetComments(postId string, limit int64, key *ExlusiveStartKey) ([]Comments, *PaginationComments, *response.RestErr) {
	dynamo := dynamoSession()

	params := &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("COMMENT_TABLE_NAME")),
		KeyConditionExpression: aws.String("#postId = :postId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":postId": {S: aws.String(postId)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#postId": aws.String("postId"),
		},
		Limit:            aws.Int64(int64(limit)),
		IndexName:        aws.String("postIdGSI"),
		ScanIndexForward: aws.Bool(false),
	}
	if key.Id != "" {
		params.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(key.Id)},
		}
	}

	query, err := dynamo.Query(params)
	if err != nil {
		return nil, nil, response.Error("Failed Query List Comments", 500, err)
	}
	var results []Comments
	for _, i := range query.Items {
		comments := Comments{}

		err = dynamodbattribute.UnmarshalMap(i, &comments)
		if err != nil {
			log.Println("Got error unmarshalling", err)
			return nil, nil, response.Error("Got error unmarshallin", 500, err)
		}

		results = append(results, comments)
	}
	pagination := PaginationComments{}
	err = dynamodbattribute.UnmarshalMap(query.LastEvaluatedKey, &pagination)
	if err != nil {
		log.Println("Got error unmarshalling", err)
		return nil, nil, response.Error("Got error unmarshallin", 500, err)
	}
	if pagination.Id == "" {
		return results, nil, nil
	}
	return results, &pagination, nil
}
