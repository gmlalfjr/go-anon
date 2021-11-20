package services

import (
	"strconv"
	"time"

	"github.com/gmlalfjr/go-anon/go-comments/src/domain"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
	"github.com/google/uuid"
)

func CreateComment(c *domain.Comments, username string, postId string) (*domain.Comments, *response.RestErr) {
	date := time.Now().UTC()
	id := uuid.New().String()

	c.Id = id
	c.Username = username
	c.CreatedAt = date
	c.ModifiedAt = date
	c.PostId = postId
	err := c.CreateComment(username, postId)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetComments(postId string, limit string, key *domain.ExlusiveStartKey) ([]domain.Comments, *domain.PaginationComments, *response.RestErr) {
	comment := &domain.Comments{}
	if limit == "" {
		limit = "10"
	}
	convertLimitDataType, errConv := strconv.ParseInt(limit, 10, 64)
	if errConv != nil {
		return nil, nil, response.Error("Failed conv query string", 500, errConv)
	}
	res, pagination, err := comment.GetComments(postId, convertLimitDataType, key)
	if err != nil {
		return nil, nil, err
	}
	return res, pagination, nil
}
