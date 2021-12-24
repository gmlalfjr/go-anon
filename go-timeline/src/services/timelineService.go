package services

import (
	"encoding/json"
	"time"

	"github.com/gmlalfjr/go-anon/go-timeline/src/domain"
	response "github.com/gmlalfjr/go_CommonResponse/utils"
	"github.com/google/uuid"
)

func PostTimeline(username string, timelineString string) (*domain.Timeline, *response.RestErr) {
	t := &domain.Timeline{}
	errUnmarshal := json.Unmarshal([]byte(timelineString), &t)
	if errUnmarshal != nil {
		return nil, &response.RestErr{
			Error: errUnmarshal.Error(),
			Message: "Error When unmarshalling",
			Status: 500,
		}
	}

	if errValidate := t .Validate(); errValidate != nil {
		return nil, &response.RestErr{
			Error: errUnmarshal.Error(),
			Message: "Error Validating Payload",
			Status: 500,
		}
	}
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
	t.Status = "OK"
	postErr := t.PostTimeline()
	if postErr != nil {
		return nil, postErr
	}

	return t, nil
}

func GetTimeline( limit int64, key map[string] string, postType string) ([]domain.Timeline, *domain.PaginationTimeline, *response.RestErr) {
	t := &domain.Timeline{
		Type: postType,
	}
	var exlusive *domain.ExlusiveStartKey
	if key["id"] != "" {
		exlusive = &domain.ExlusiveStartKey{
			Id:        key["id"],
			Type:      key["type"],
			CreatedAt: key["createdAt"],
			Status:    key["status"],
		}
	}
	res, pagination, getErr := t.GetTimeline(limit, exlusive)
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

func GetOwnPost(username string, limit int64) ([]domain.Timeline, *domain.PaginationTimeline, *response.RestErr) {
	exlusive := &domain.ExlusiveStartKeyUsername{}
	t := &domain.Timeline{}
	res, pagination, getErr := t.GetOwnUserPost(username, limit, exlusive)
	if getErr != nil {
		return nil, nil, getErr
	}
	return res, pagination, nil
}
