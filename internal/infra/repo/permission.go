package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/util"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type PermissionEntity struct {
	Id        int32 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Code      string
	Name      string

	PermissionType int32

	description string
	ParentId    int32

	UserId int // 和登录账号绑定
	// TenantId int
	Version int32 `gorm:"default:1"`
}

func (t *PermissionEntity) TableName() string {
	return "sys_permission"

}

type permissionRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.PermissionRepo.
func (*permissionRepo) BatchSave(context.Context, []*domain.Permission) error {
	panic("unimplemented")
}

// Delete implements domain.PermissionRepo.
func (*permissionRepo) Delete(context.Context, int32) (int32, error) {
	panic("unimplemented")
}

// ListAll implements domain.PermissionRepo.
func (*permissionRepo) ListAll(context.Context, *domain.Expression) ([]*domain.Permission, error) {
	panic("unimplemented")
}

// ListByMap implements domain.PermissionRepo.
func (pr *permissionRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.Permission, error) {
	var permissionList = make([]*PermissionEntity, 0)
	err := pr.repo.GetDBFromContext(ctx).Find(&permissionList, params).Error
	if err != nil {
		return nil, err
	}
	var result = make([]*domain.Permission, 0, len(permissionList))
	_ = util.CopyProperties(&permissionList, &result, util.IgnoreNotMatchedProperty())
	return result, nil
}

// FindByID implements domain.PermissionRepo.
func (up *permissionRepo) FindByID(ctx context.Context, id int32) (*domain.Permission, error) {

	var permission PermissionEntity

	err := up.repo.GetDBFromContext(ctx).First(&permission, id).Error

	if err != nil {
		return nil, err
	}
	var result domain.Permission
	_ = util.CopyProperties(&permission, &result, util.IgnoreNotMatchedProperty())

	return &result, nil
}

// List implements domain.PermissionRepo.
func (up *permissionRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.Permission, error) {

	var permissionList []*PermissionEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&permissionList).Error
	if err != nil {
		return nil, err
	}

	var result = make([]*domain.Permission, 0, len(permissionList))
	_ = util.CopyProperties(&permissionList, &result, util.IgnoreNotMatchedProperty())

	return result, nil

}

// Count implements domain.PermissionRepo.
func (up *permissionRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&PermissionEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.PermissionReppo.
func (up *permissionRepo) Save(ctx context.Context, permission *domain.Permission) (*domain.Permission, error) {
	var permissionEntity PermissionEntity

	_ = util.CopyProperties(permission, &permissionEntity, util.IgnoreNotMatchedProperty())

	err := up.repo.GetDBFromContext(ctx).Create(&permissionEntity).Error

	if err != nil {
		return nil, err
	}
	permission.Id = int32(permissionEntity.Id)
	return permission, nil

}

// Update implements domain.PermissionRepo.
func (up *permissionRepo) Update(ctx context.Context, permission *domain.Permission) (int32, error) {

	var permissionEntity PermissionEntity

	_ = util.CopyProperties(permission, &permissionEntity, util.IgnoreNotMatchedProperty())

	result := up.repo.GetDBFromContext(ctx).Model(&PermissionEntity{Id: permission.Id}).Updates(permissionEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewPermissionRepo .
func NewPermissionRepo(repo *BaseRepo, logger log.Logger) domain.PermissionRepo {
	return &permissionRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
