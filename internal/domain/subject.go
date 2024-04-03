package domain

import (
	"context"
)

type SubjectRepo interface {
	Save(context.Context, *Subject) (*Subject, error)
	Update(context.Context, *Subject) (int32, error)
	FindByID(context.Context, int32) (*Subject, error)
	List(context.Context, *Expression, *Page) ([]*Subject, error)
	Count(context.Context, *Expression) (int32, error)
}

type Subject struct {
	Id          int32
	Name        string
	Category    int32
	Description string
	Status      int32
	UpdatedAt   string
	Version     int32
}

type SubjectService struct {
	repo SubjectRepo
}

func NewSubjectService(repo SubjectRepo) *SubjectService {

	return &SubjectService{repo: repo}
}

func (t *SubjectService) Create(ctx context.Context, subject *Subject) (*Subject, error) {

	return t.repo.Save(ctx, subject)
}

func (t *SubjectService) Update(ctx context.Context, subject *Subject) (int32, error) {

	return t.repo.Update(ctx, subject)
}

func (t *SubjectService) FindByID(ctx context.Context, id int32) (*Subject, error) {

	return t.repo.FindByID(ctx, id)
}

func (t *SubjectService) List(ctx context.Context, query *Expression, page *Page) ([]*Subject, int32, error) {

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
