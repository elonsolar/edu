package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/util"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type TenantPermissionEntity struct {
	Id           int32 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	TenantId     int32
	PermissionId int32
	Version      int32 `gorm:"default:1"`
}

func (t *TenantPermissionEntity) TableName() string {
	return "tenant_permission"

}

type tenantPermissionRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.PermissionRepo.
func (pr *tenantPermissionRepo) BatchSave(ctx context.Context, tenantPermissionList []*domain.TenantPermission) error {

	var tenantPermissionEntityList = make([]*TenantPermissionEntity, 0, len(tenantPermissionList))
	_ = util.CopyProperties(&tenantPermissionList, &tenantPermissionEntityList, util.IgnoreNotMatchedProperty())
	return pr.repo.GetDBFromContext(ctx).Create(tenantPermissionEntityList).Error
}

// Delete implements domain.PermissionRepo.
func (*tenantPermissionRepo) Delete(context.Context, int32) (int32, error) {
	panic("unimplemented")
}

// ListAll implements domain.PermissionRepo.
func (*tenantPermissionRepo) ListAll(context.Context, *domain.Expression) ([]*domain.TenantPermission, error) {
	panic("unimplemented")
}

func (pr *tenantPermissionRepo) BatchDelete(ctx context.Context, tenantId int32, permissionId []int32) error {

	return pr.repo.GetDBFromContext(ctx).Where("tenant_id=?", tenantId).Where("permission_id in ?", permissionId).Delete(&TenantPermissionEntity{}).Error
}

// ListByMap implements domain.PermissionRepo.
func (pr *tenantPermissionRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.TenantPermission, error) {
	var tenantList = make([]*TenantPermissionEntity, 0)
	err := pr.repo.GetDBFromContext(ctx).Find(&tenantList, params).Error
	if err != nil {
		return nil, err
	}
	var result = make([]*domain.TenantPermission, 0, len(tenantList))
	_ = util.CopyProperties(&tenantList, &result, util.IgnoreNotMatchedProperty())
	return result, nil
}

// FindByID implements domain.PermissionRepo.
func (up *tenantPermissionRepo) FindByID(ctx context.Context, id int32) (*domain.TenantPermission, error) {

	var tenant TenantPermissionEntity

	err := up.repo.GetDBFromContext(ctx).First(&tenant, id).Error

	if err != nil {
		return nil, err
	}
	var result domain.TenantPermission
	_ = util.CopyProperties(&tenant, &result, util.IgnoreNotMatchedProperty())

	return &result, nil
}

// List implements domain.PermissionRepo.
func (up *tenantPermissionRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.TenantPermission, error) {

	var tenantList []*PermissionEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&tenantList).Error
	if err != nil {
		return nil, err
	}

	var result = make([]*domain.TenantPermission, 0, len(tenantList))
	_ = util.CopyProperties(&tenantList, &result, util.IgnoreNotMatchedProperty())

	return result, nil

}

// Count implements domain.PermissionRepo.
func (up *tenantPermissionRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&PermissionEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.PermissionReppo.
func (up *tenantPermissionRepo) Save(ctx context.Context, tenant *domain.TenantPermission) (*domain.TenantPermission, error) {

	var tenantEntity TenantPermissionEntity

	_ = util.CopyProperties(tenant, &tenantEntity, util.IgnoreNotMatchedProperty())

	err := up.repo.GetDBFromContext(ctx).Create(&tenantEntity).Error

	if err != nil {
		return nil, err
	}
	tenant.Id = int32(tenantEntity.Id)
	return tenant, nil

}

// Update implements domain.PermissionRepo.
func (up *tenantPermissionRepo) Update(ctx context.Context, tenant *domain.TenantPermission) (int32, error) {

	var tenantEntity TenantPermissionEntity

	_ = util.CopyProperties(tenant, &tenantEntity, util.IgnoreNotMatchedProperty())

	result := up.repo.GetDBFromContext(ctx).Model(&TenantPermissionEntity{Id: tenant.Id}).Updates(&tenantEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewPermissionRepo .
func NewTenantPermissionRepo(repo *BaseRepo, logger log.Logger) domain.TenantPermissionRepo {
	return &tenantPermissionRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
