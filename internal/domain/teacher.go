package domain

import (
	"context"
)

type TeacherRepo interface {
	Save(context.Context, *Teacher) (*Teacher, error)
	Update(context.Context, *Teacher) (int32, error)
	FindByID(context.Context, int32) (*Teacher, error)
	List(context.Context, *Expression, *Page) ([]*Teacher, error)
	Count(context.Context, *Expression) (int32, error)
}

type Teacher struct {
	Id          int32
	Name        string
	Mobile      string
	Status      int32
	Description string
	UpdatedAt   string
	Version     int32
}

type TeacherService struct {
	repo TeacherRepo
}

func NewTeacherService(repo TeacherRepo) *TeacherService {

	return &TeacherService{repo: repo}
}

func (t *TeacherService) Create(ctx context.Context, teacher *Teacher) (*Teacher, error) {

	return t.repo.Save(ctx, teacher)
}

func (t *TeacherService) Update(ctx context.Context, teacher *Teacher) (int32, error) {

	return t.repo.Update(ctx, teacher)
}

func (t *TeacherService) FindByID(ctx context.Context, id int32) (*Teacher, error) {

	return t.repo.FindByID(ctx, id)
}

func (t *TeacherService) List(ctx context.Context, query *Expression, page *Page) ([]*Teacher, int32, error) {

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
