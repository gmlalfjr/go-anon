package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type ExlusiveStartKey struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type PaginationComments struct {
	Id string `json:"id"`
}

type Comments struct {
	Id            string    `json:"id"`
	Username      string    `json:"username"`
	PostId        string    `json:"postId"`
	Comment       string    `json:"comment" validate:"required"`
	GeneratedName string    `json:"generatedName"`
	TotalReplies  int8      `json:"totalReplies"`
	TotalLike     int8      `json:"totalLike"`
	CreatedAt     time.Time `json:"createdAt"`
	ModifiedAt    time.Time `json:"modifiedAt"`
}

func (c *Comments) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
