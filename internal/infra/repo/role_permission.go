package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/util"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type RolePermissionEntity struct {
	Id           int32 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	RoleId       int32
	PermissionId int32

	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32 `gorm:"default:1"`
}

func (t *RolePermissionEntity) TableName() string {
	return "role_permission"

}

type rolePermissionRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.PermissionRepo.
func (pr *rolePermissionRepo) BatchSave(ctx context.Context, rolePermissionList []*domain.RolePermission) error {

	var rolePermissionEntityList = make([]*RolePermissionEntity, 0, len(rolePermissionList))
	_ = util.CopyProperties(&rolePermissionList, &rolePermissionEntityList, util.IgnoreNotMatchedProperty())
	return pr.repo.GetDBFromContext(ctx).Create(rolePermissionEntityList).Error
}

// Delete implements domain.PermissionRepo.
func (*rolePermissionRepo) Delete(context.Context, int32) (int32, error) {
	panic("unimplemented")
}

// ListAll implements domain.PermissionRepo.
func (*rolePermissionRepo) ListAll(context.Context, *domain.Expression) ([]*domain.RolePermission, error) {
	panic("unimplemented")
}

func (pr *rolePermissionRepo) DeleteRolePermissions(ctx context.Context, roleId int32, permissionId []int32) error {

	return pr.repo.GetDBFromContext(ctx).Where("role_id=?", roleId).Where("permission_id in ?", permissionId).Delete(&RolePermissionEntity{}).Error
}

func (pr *rolePermissionRepo) DeleteTenantPermissions(ctx context.Context, tenantId int32, permissionId []int32) error {

	return pr.repo.GetDBFromContext(ctx).Where("tenant_id=?", tenantId).Where("permission_id in ?", permissionId).Delete(&RolePermissionEntity{}).Error
}

// ListByMap implements domain.PermissionRepo.
func (pr *rolePermissionRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.RolePermission, error) {
	var roleList = make([]*RolePermissionEntity, 0)
	err := pr.repo.GetDBFromContext(ctx).Find(&roleList, params).Error
	if err != nil {
		return nil, err
	}
	var result = make([]*domain.RolePermission, 0, len(roleList))
	_ = util.CopyProperties(&roleList, &result, util.IgnoreNotMatchedProperty())
	return result, nil
}

// FindByID implements domain.PermissionRepo.
func (up *rolePermissionRepo) FindByID(ctx context.Context, id int32) (*domain.RolePermission, error) {

	var role RolePermissionEntity

	err := up.repo.GetDBFromContext(ctx).First(&role, id).Error

	if err != nil {
		return nil, err
	}
	var result domain.RolePermission
	_ = util.CopyProperties(&role, &result, util.IgnoreNotMatchedProperty())

	return &result, nil
}

// List implements domain.PermissionRepo.
func (up *rolePermissionRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.RolePermission, error) {

	var roleList []*PermissionEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&roleList).Error
	if err != nil {
		return nil, err
	}

	var result = make([]*domain.RolePermission, 0, len(roleList))
	_ = util.CopyProperties(&roleList, &result, util.IgnoreNotMatchedProperty())

	return result, nil

}

// Count implements domain.PermissionRepo.
func (up *rolePermissionRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&PermissionEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.PermissionReppo.
func (up *rolePermissionRepo) Save(ctx context.Context, role *domain.RolePermission) (*domain.RolePermission, error) {

	var roleEntity RolePermissionEntity

	_ = util.CopyProperties(role, &roleEntity, util.IgnoreNotMatchedProperty())

	err := up.repo.GetDBFromContext(ctx).Create(&roleEntity).Error

	if err != nil {
		return nil, err
	}
	role.Id = int32(roleEntity.Id)
	return role, nil

}

// Update implements domain.PermissionRepo.
func (up *rolePermissionRepo) Update(ctx context.Context, role *domain.RolePermission) (int32, error) {

	var roleEntity RolePermissionEntity

	_ = util.CopyProperties(role, &roleEntity, util.IgnoreNotMatchedProperty())

	result := up.repo.GetDBFromContext(ctx).Model(&RolePermissionEntity{Id: role.Id}).Updates(&roleEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewPermissionRepo .
func NewRolePermissionRepo(repo *BaseRepo, logger log.Logger) domain.RolePermissionRepo {
	return &rolePermissionRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
