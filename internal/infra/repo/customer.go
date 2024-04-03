package repo

import (
	"context"
	"edu/internal/domain"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type CustomerEntity struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name        string
	Mobile      string
	Description string
	Community   string

	LessonNumber int32

	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32 `gorm:"default:1"`
}

func (t *CustomerEntity) TableName() string {
	return "customer"

}

type customerRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// FindByID implements domain.CustomerRepo.
func (up *customerRepo) FindByID(ctx context.Context, id int32) (*domain.Customer, error) {

	var customer CustomerEntity

	err := up.repo.GetDBFromContext(ctx).First(&customer, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.Customer{
		Id:           int32(customer.ID),
		Name:         customer.Name,
		Mobile:       customer.Mobile,
		Description:  customer.Description,
		Community:    customer.Community,
		LessonNumber: customer.LessonNumber,
		Version:      customer.Version,
	}, nil
}

// List implements domain.CustomerRepo.
func (up *customerRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.Customer, error) {

	var customerList []*CustomerEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&customerList).Error
	if err != nil {
		return nil, err
	}

	return Map(customerList, func(customer *CustomerEntity) *domain.Customer {
		return &domain.Customer{
			Id:           int32(customer.ID),
			Name:         customer.Name,
			Mobile:       customer.Mobile,
			Description:  customer.Description,
			Community:    customer.Community,
			LessonNumber: customer.LessonNumber,
			Version:      customer.Version,
			UpdatedAt:    customer.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}), nil

}

// Count implements domain.CustomerRepo.
func (up *customerRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&CustomerEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.CustomerRepo.
func (up *customerRepo) Save(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	var customerEntity = &CustomerEntity{
		Name:         customer.Name,
		Mobile:       customer.Mobile,
		Description:  customer.Description,
		Community:    customer.Community,
		LessonNumber: customer.LessonNumber,
	}
	err := up.repo.GetDBFromContext(ctx).Save(customerEntity).Error
	if err != nil {
		return nil, err
	}
	customer.Id = int32(customerEntity.ID)
	return customer, nil

}

// Update implements domain.CustomerRepo.
func (up *customerRepo) Update(ctx context.Context, customer *domain.Customer) (int32, error) {
	var customerEntity = &CustomerEntity{
		ID:          uint(customer.Id),
		Name:        customer.Name,
		Mobile:      customer.Mobile,
		Description: customer.Description,
		Community:   customer.Community,
		// LessonNumber: customer.LessonNumber,
		Version: customer.Version,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&CustomerEntity{ID: uint(customerEntity.ID)}).Updates(customerEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}
func (up *customerRepo) ChangeLessonNum(ctx context.Context, customer *domain.Customer) (int32, error) {
	var customerEntity = &CustomerEntity{
		ID:           uint(customer.Id),
		LessonNumber: customer.LessonNumber,
		Version:      customer.Version,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&CustomerEntity{ID: uint(customerEntity.ID)}).Updates(customerEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewCustomerRepo .
func NewCustomerRepo(repo *BaseRepo, logger log.Logger) domain.CustomerRepo {
	return &customerRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
