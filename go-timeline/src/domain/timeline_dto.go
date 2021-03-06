package domain

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type PaginationTimeline struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	CreatedAt string `json:"createdAt"`
	Status    string `json:"status"`
}

type PaginationTimelineById struct {
	Id   string `json:"id"`
	CreatedAt string `json:"createdAt"`
	Username string `json:"username"`
}

type ExlusiveStartKey struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt"`
	Status    string `json:"status"`
}

type ExlusiveStartKeyUsername struct {
	Id       string `json:"id"`
	Username string `json:"username"`
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
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	ModifiedAt   time.Time `json:"modifiedAt"`
}

func (t *Timeline) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}

func (t * Timeline) ValidateMapType(evaluatedKey map[string] string, timelineType string) error {
	if evaluatedKey["id"] != "" {
		if timelineType == "ALL" {
			if evaluatedKey["status"] == "" || evaluatedKey["createdAt"] == "" {
				return errors.New("must Input all parameter for pagination")
			}
		} else {
			if evaluatedKey["type"] == "" || evaluatedKey["createdAt"] == "" {
				return errors.New("must Input all parameter for pagination")
			}
		}
	}
	return nil
}


func (t * Timeline) ValidateGetOwnAllPost(evaluatedKey map[string] string, username string) error {
	if evaluatedKey["id"] != "" {
		if evaluatedKey["username"] == "" || evaluatedKey["createdAt"] == "" {
			return errors.New("must Input all parameter for pagination")
		}
	}
	return nil
}
