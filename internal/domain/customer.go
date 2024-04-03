package domain

import (
	"context"
)

type CustomerRepo interface {
	Save(context.Context, *Customer) (*Customer, error)
	Update(context.Context, *Customer) (int32, error)
	FindByID(context.Context, int32) (*Customer, error)
	List(context.Context, *Expression, *Page) ([]*Customer, error)
	Count(context.Context, *Expression) (int32, error)
	ChangeLessonNum(context.Context, *Customer) (int32, error)
}

type Customer struct {
	Id           int32
	Name         string
	Mobile       string
	Description  string
	Community    string
	LessonNumber int32
	UpdatedAt    string
	Version      int32
}

type CustomerService struct {
	repo CustomerRepo
}

func NewCustomerService(repo CustomerRepo) *CustomerService {

	return &CustomerService{repo: repo}
}

func (t *CustomerService) Create(ctx context.Context, customer *Customer) (*Customer, error) {

	return t.repo.Save(ctx, customer)
}

func (t *CustomerService) Update(ctx context.Context, customer *Customer) (int32, error) {

	return t.repo.Update(ctx, customer)
}

func (t *CustomerService) FindByID(ctx context.Context, id int32) (*Customer, error) {

	return t.repo.FindByID(ctx, id)
}

func (t *CustomerService) List(ctx context.Context, query *Expression, page *Page) ([]*Customer, int32, error) {

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

func (t *CustomerService) ChangeLessonNum(ctx context.Context, id int32, numChange int32, version int32) error {

	customer, err := t.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if customer.Version != version {
		return ErrCurrencyUpdate
	}

	count, err := t.repo.ChangeLessonNum(ctx, &Customer{
		Id:           customer.Id,
		LessonNumber: customer.LessonNumber + numChange,
		Version:      customer.Version,
	})
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrCurrencyUpdate
	}
	return nil
}
