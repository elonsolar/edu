package domain

import (
	"context"
)

type RoomRepo interface {
	Save(context.Context, *Room) (*Room, error)
	Update(context.Context, *Room) (int32, error)
	FindByID(context.Context, int32) (*Room, error)
	List(context.Context, *Expression, *Page) ([]*Room, error)
	Count(context.Context, *Expression) (int32, error)
}

type Room struct {
	Id          int32
	Code        string
	Subjects    []int32
	Description string
	Status      int32
	UpdatedAt   string
	Version     int32
}

type RoomService struct {
	repo RoomRepo
}

func NewRoomService(repo RoomRepo) *RoomService {

	return &RoomService{repo: repo}
}

func (t *RoomService) Create(ctx context.Context, teacher *Room) (*Room, error) {

	return t.repo.Save(ctx, teacher)
}

func (t *RoomService) Update(ctx context.Context, teacher *Room) (int32, error) {

	return t.repo.Update(ctx, teacher)
}

func (t *RoomService) FindByID(ctx context.Context, id int32) (*Room, error) {

	return t.repo.FindByID(ctx, id)
}

func (t *RoomService) List(ctx context.Context, query *Expression, page *Page) ([]*Room, int32, error) {

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
