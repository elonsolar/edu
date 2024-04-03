package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/util"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type RoleEntity struct {
	Id        int32 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Code      string
	Name      string

	Description string

	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32 `gorm:"default:1"`
}

func (t *RoleEntity) TableName() string {
	return "role"

}

type roleRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.RoleRepo.
func (*roleRepo) BatchSave(context.Context, []*domain.Role) error {
	panic("unimplemented")
}

// Delete implements domain.RoleRepo.
func (*roleRepo) Delete(context.Context, int32) (int32, error) {
	panic("unimplemented")
}

// ListAll implements domain.RoleRepo.
func (*roleRepo) ListAll(context.Context, *domain.Expression) ([]*domain.Role, error) {
	panic("unimplemented")
}

// ListByMap implements domain.RoleRepo.
func (pr *roleRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.Role, error) {
	var roleList = make([]*RoleEntity, 0)
	err := pr.repo.GetDBFromContext(ctx).Find(&roleList, params).Error
	if err != nil {
		return nil, err
	}
	var result = make([]*domain.Role, 0, len(roleList))
	_ = util.CopyProperties(&roleList, &result, util.IgnoreNotMatchedProperty())
	return result, nil
}

// FindByID implements domain.RoleRepo.
func (up *roleRepo) FindByID(ctx context.Context, id int32) (*domain.Role, error) {

	var role RoleEntity

	err := up.repo.GetDBFromContext(ctx).First(&role, id).Error

	if err != nil {
		return nil, err
	}
	var result domain.Role
	_ = util.CopyProperties(&role, &result, util.IgnoreNotMatchedProperty())

	return &result, nil
}

// List implements domain.RoleRepo.
func (up *roleRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.Role, error) {

	var roleList []*RoleEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&roleList).Error
	if err != nil {
		return nil, err
	}

	var result = make([]*domain.Role, 0, len(roleList))
	_ = util.CopyProperties(&roleList, &result, util.IgnoreNotMatchedProperty())

	return result, nil

}

// Count implements domain.RoleRepo.
func (up *roleRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&RoleEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.RoleReppo.
func (up *roleRepo) Save(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	var roleEntity RoleEntity

	_ = util.CopyProperties(role, &roleEntity, util.IgnoreNotMatchedProperty())

	err := up.repo.GetDBFromContext(ctx).Create(&roleEntity).Error

	if err != nil {
		return nil, err
	}
	role.Id = int32(roleEntity.Id)
	return role, nil

}

// Update implements domain.RoleRepo.
func (up *roleRepo) Update(ctx context.Context, role *domain.Role) (int32, error) {

	var roleEntity RoleEntity

	_ = util.CopyProperties(role, &roleEntity, util.IgnoreNotMatchedProperty())

	result := up.repo.GetDBFromContext(ctx).Model(&RoleEntity{Id: role.Id}).Updates(&roleEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewRoleRepo .
func NewRoleRepo(repo *BaseRepo, logger log.Logger) domain.RoleRepo {
	return &roleRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
