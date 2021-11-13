package services

import (
	"time"

	"github.com/gmlalfjr/go-anon/go-timeline/src/domain"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
	"github.com/google/uuid"
)

func PostTimeline(t *domain.Timeline, username string) (*domain.Timeline, *response.RestErr) {
	date := time.Now().UTC()
	id := uuid.New().String()
	timelineId := uuid.New().String()

	t.Id = id
	t.CreatedAt = date
	t.ModifiedAt = date
	t.Username = username
	t.TotalComment = 0
	t.TotalLike = 0
	t.TotalReport = 0
	t.TimelineId = timelineId

	postErr := t.PostTimeline()
	if postErr != nil {
		return nil, postErr
	}

	return t, nil
}

func GetTimeline(t *domain.Timeline, limit int64, key *domain.ExlusiveStartKey) ([]domain.Timeline, *domain.PaginationTimeline, *response.RestErr) {
	res, pagination, getErr := t.GetTimeline(limit, key)
	if getErr != nil {
		return nil, nil, getErr
	}
	return res, pagination, nil

}

func GetTimelineDetail(id string) (*domain.Timeline, *response.RestErr) {
	t := &domain.Timeline{}
	res, getErr := t.GetTimelineDetail(id)
	if getErr != nil {
		return nil, getErr
	}
	return res, nil

}

func DeleteUserPost(id string, username string) (bool, *response.RestErr) {
	t := &domain.Timeline{}
	res, getErr := t.DeleteUserPost(id, username)
	if getErr != nil {
		return false, getErr
	}
	return res, nil

}
