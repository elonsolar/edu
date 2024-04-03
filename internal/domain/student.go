package domain

import (
	"context"
	"time"
)

type StudentRepo interface {
	Save(context.Context, *Student) (*Student, error)
	Update(context.Context, *Student) (int32, error)
	FindByID(context.Context, int32) (*Student, error)
	List(context.Context, *Expression, *Page) ([]*Student, error)
	Count(context.Context, *Expression) (int32, error)
}

type Student struct {
	Id          int32
	Name        string
	Mobile      string
	Birthday    time.Time
	Description string
	UpdatedAt   string
	Version     int32
}

type StudentService struct {
	repo StudentRepo
}

func NewStudentService(repo StudentRepo) *StudentService {

	return &StudentService{repo: repo}
}

func (t *StudentService) Create(ctx context.Context, teacher *Student) (*Student, error) {

	return t.repo.Save(ctx, teacher)
}

func (t *StudentService) Update(ctx context.Context, teacher *Student) (int32, error) {

	return t.repo.Update(ctx, teacher)
}

func (t *StudentService) FindByID(ctx context.Context, id int32) (*Student, error) {

	return t.repo.FindByID(ctx, id)
}

func (t *StudentService) List(ctx context.Context, query *Expression, page *Page) ([]*Student, int32, error) {

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
