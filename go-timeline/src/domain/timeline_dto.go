package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type PaginationTimeline struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type ExlusiveStartKey struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type Timeline struct {
	Id           string    `json:"id"`
	TimelineId   string    `json:"timelineId"`
	Username     string    `json:"username" validate:"required"`
	Post         string    `json:"post" validate:"required"`
	Type         string    `json:"type" validate:"required,oneof='CURHAT' 'CARI_PASANGAN'"`
	TotalLike    int8      `json:"totalLike"`
	TotalComment int8      `json:"totalComment"`
	IsPrivate    bool      `json:"isPrivate"`
	TotalReport  int8      `json:"totalReport"`
	CreatedAt    time.Time `json:"createdAt"`
	ModifiedAt   time.Time `json:"modifiedAt"`
}

func (t *Timeline) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}
