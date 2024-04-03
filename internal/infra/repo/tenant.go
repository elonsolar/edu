package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/domain/enum"
	"edu/internal/util"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type TenantEntity struct {
	Id          int32 `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string
	Mobile      string
	Description string
	Status      int32
	Version     int32 `gorm:"default:1"`
}

func (u *TenantEntity) TableName() string {
	return "tenant"
}

type tenantRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.TenantRepo.
func (*tenantRepo) BatchSave(context.Context, []*domain.Tenant) error {
	panic("unimplemented")
}

// Delete implements domain.TenantRepo.
func (*tenantRepo) Delete(context.Context, int32) (int32, error) {
	panic("unimplemented")
}

// ListAll implements domain.TenantRepo.
func (*tenantRepo) ListAll(context.Context, *domain.Expression) ([]*domain.Tenant, error) {
	panic("unimplemented")
}

// ListByMap implements domain.TenantRepo.
func (up *tenantRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.Tenant, error) {
	var tenantEntityList = make([]*TenantEntity, 0)
	err := up.repo.GetDBFromContext(ctx).Model(&TenantEntity{}).Find(&tenantEntityList, params).Error
	if err != nil {
		return nil, err
	}
	var tenantList = make([]*domain.Tenant, 0, len(tenantEntityList))
	_ = util.CopyProperties(&tenantEntityList, &tenantList, util.IgnoreNotMatchedProperty())
	return tenantList, nil
}

// FindByID implements domain.TenantRepo.
func (up *tenantRepo) FindByID(ctx context.Context, id int32) (*domain.Tenant, error) {

	var tenant TenantEntity

	err := up.repo.GetDBFromContext(ctx).First(&tenant, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.Tenant{
		Id:          int32(tenant.Id),
		Name:        tenant.Name,
		Description: tenant.Description,
		Status:      tenant.Status,
		Version:     tenant.Version,
	}, nil
}

// List implements domain.TenantRepo.
func (up *tenantRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.Tenant, error) {

	var tenantList []*TenantEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&tenantList).Error
	if err != nil {
		return nil, err
	}

	return Map(tenantList, func(tenant *TenantEntity) *domain.Tenant {
		return &domain.Tenant{
			Id:          int32(tenant.Id),
			Name:        tenant.Name,
			Mobile:      tenant.Mobile,
			Description: tenant.Description,
			Status:      tenant.Status,
			UpdatedAt:   tenant.UpdatedAt.Format("2006-01-02 15:04:05"),
			Version:     tenant.Version,
		}
	}), nil

}

// Count implements domain.TenantRepo.
func (up *tenantRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&TenantEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.TenantRepo.
func (up *tenantRepo) Save(ctx context.Context, tenant *domain.Tenant) (*domain.Tenant, error) {
	var tenantEntity = &TenantEntity{
		Name:        tenant.Name,
		Mobile:      tenant.Mobile,
		Description: tenant.Description,
		Status:      int32(enum.EnableStatusEnabled),
	}
	err := up.repo.GetDBFromContext(ctx).Save(tenantEntity).Error
	if err != nil {
		return nil, err
	}
	tenant.Id = int32(tenantEntity.Id)
	return tenant, nil

}

// Update implements domain.TenantRepo.
func (up *tenantRepo) Update(ctx context.Context, tenant *domain.Tenant) (int32, error) {

	var tenantEntity = &TenantEntity{
		Id:          tenant.Id,
		Name:        tenant.Name,
		Mobile:      tenant.Mobile,
		Description: tenant.Description,
		Status:      tenant.Status,
		Version:     tenant.Version,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&TenantEntity{Id: tenant.Id}).Updates(tenantEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewTenantRepo .
func NewTenantRepo(repo *BaseRepo, logger log.Logger) domain.TenantRepo {
	return &tenantRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
