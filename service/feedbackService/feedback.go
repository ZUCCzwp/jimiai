package feedbackService

import (
	"errors"
	"jiyu/model/contextModel"
	"jiyu/model/globalModel"
	"jiyu/repo/feedbackRepo"
	"time"
)

func SaveFeedback(ctx contextModel.Context, data globalModel.Feedback) error {
	return feedbackRepo.SaveFeedback(data)
}

func GetFeedback(ctx contextModel.Context, feedbackId int) (globalModel.Feedback, error) {
	return feedbackRepo.GetFeedback(int(ctx.User.ID), feedbackId)
}

func ReadFeedback(ctx contextModel.Context, feedbackId int) error {
	feedback, err := feedbackRepo.GetFeedback(int(ctx.User.ID), feedbackId)
	if err != nil {
		return err
	}

	if feedback.UserId != int(ctx.User.ID) {
		return errors.New("无权限")
	}

	return feedbackRepo.ReadFeedback(int(ctx.User.ID), feedbackId)
}

func IsNewFeedbackReply(ctx contextModel.Context) (bool, error) {
	return feedbackRepo.IsNewFeedbackReply(int(ctx.User.ID))
}

func GetQaList(ctx contextModel.Context, page, pageSize int) ([]globalModel.QAListResponse, int, error) {
	reports, err := feedbackRepo.GetQaList(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	count, err := feedbackRepo.CountQa()
	if err != nil {
		return nil, 0, err
	}

	return reports, count, nil
}

func GetList(ctx contextModel.Context, page, pageSize int) ([]globalModel.FeedbackListResponse, int, error) {
	list, err := feedbackRepo.GetList(int(ctx.User.ID), page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	count, err := feedbackRepo.Count(int(ctx.User.ID))
	if err != nil {
		return nil, 0, err
	}

	result := make([]globalModel.FeedbackListResponse, 0)
	for _, v := range list {
		result = append(result, globalModel.FeedbackListResponse{
			ID:           int(v.ID),
			Nickname:     v.Nickname,
			FeedbackType: v.FeedbackType,
			Content:      v.Content,
			Status:       v.Status,
			ProcessTime:  v.ProcessTime.Format(time.DateTime),
			Replay:       v.Replay,
			IsRead:       v.IsRead,
		})
	}

	return result, count, nil
}
