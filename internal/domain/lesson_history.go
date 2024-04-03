package domain

import (
	"context"
)

type LessonHistoryRepo interface {
	Save(context.Context, *LessonHistory) (*LessonHistory, error)
	Update(context.Context, *LessonHistory) (int32, error)
	FindByID(context.Context, int32) (*LessonHistory, error)
	List(context.Context, *Expression, *Page) ([]*LessonHistory, error)
	Count(context.Context, *Expression) (int32, error)
}

type LessonHistory struct {
	Id          int32
	Mobile      string
	SourceType  int32
	OriginNum   int32
	NumChange   int32
	DetailId    int32
	Description string
	UpdatedAt   string
	Version     int32
}

type LessonHistoryService struct {
	repo LessonHistoryRepo
}

func NewLessonHistoryService(repo LessonHistoryRepo) *LessonHistoryService {

	return &LessonHistoryService{repo: repo}
}

func (t *LessonHistoryService) Create(ctx context.Context, lessonHistory *LessonHistory) (*LessonHistory, error) {

	return t.repo.Save(ctx, lessonHistory)
}

func (t *LessonHistoryService) Update(ctx context.Context, lessonHistory *LessonHistory) (int32, error) {

	return t.repo.Update(ctx, lessonHistory)
}

func (t *LessonHistoryService) FindByID(ctx context.Context, id int32) (*LessonHistory, error) {

	return t.repo.FindByID(ctx, id)
}

func (t *LessonHistoryService) List(ctx context.Context, query *Expression, page *Page) ([]*LessonHistory, int32, error) {

	count, err := t.repo.Count(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	if count == 0 {
		return nil, 0, nil
	}

	list, err := t.repo.List(ctx, query, page)
	if err != nil {
		return nil, 0, err
	}

	return list, count, nil
}
