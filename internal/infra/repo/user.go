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

type UserEntity struct {
	Id           int32 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Username     string
	Password     string
	Salt         string
	Mobile       string
	Avatar       string
	Description  string
	RoleId       int32
	RoleName     string
	Status       int32
	TenantId     int32
	Version      int32 `gorm:"default:1"`
	IsSuperAdmin bool
	IsTenant     bool
}

func (u *UserEntity) TableName() string {
	return "user"
}

type userRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.UserRepo.
func (*userRepo) BatchSave(context.Context, []*domain.User) error {
	panic("unimplemented")
}

// Delete implements domain.UserRepo.
func (*userRepo) Delete(context.Context, int32) (int32, error) {
	panic("unimplemented")
}

// ListAll implements domain.UserRepo.
func (*userRepo) ListAll(context.Context, *domain.Expression) ([]*domain.User, error) {
	panic("unimplemented")
}

// ListByMap implements domain.UserRepo.
func (up *userRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.User, error) {
	var userEntityList = make([]*UserEntity, 0)
	err := up.repo.GetDBFromContext(ctx).Model(&UserEntity{}).Find(&userEntityList, params).Error
	if err != nil {
		return nil, err
	}
	var userList = make([]*domain.User, 0, len(userEntityList))
	_ = util.CopyProperties(&userEntityList, &userList, util.IgnoreNotMatchedProperty())
	return userList, nil
}

// FindByID implements domain.UserRepo.
func (up *userRepo) FindByID(ctx context.Context, id int32) (*domain.User, error) {

	var user UserEntity

	err := up.repo.GetDBFromContext(ctx).First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.User{
		Id:           int32(user.Id),
		Username:     user.Username,
		Password:     user.Password,
		Mobile:       user.Mobile,
		Description:  user.Description,
		Avatar:       user.Avatar,
		Status:       user.Status,
		RoleId:       user.RoleId,
		RoleName:     user.RoleName,
		Version:      user.Version,
		IsSuperAdmin: user.IsSuperAdmin,
		IsTenant:     user.IsTenant,
	}, nil
}

// List implements domain.UserRepo.
func (up *userRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.User, error) {

	var userList []*UserEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&userList).Error
	if err != nil {
		return nil, err
	}

	return Map(userList, func(user *UserEntity) *domain.User {
		return &domain.User{
			Id:          int32(user.Id),
			Username:    user.Username,
			Mobile:      user.Mobile,
			Avatar:      user.Avatar,
			Description: user.Description,
			RoleId:      user.RoleId,
			RoleName:    user.RoleName,
			Status:      user.Status,
			UpdatedAt:   user.UpdatedAt.Format("2006-01-02 15:04:05"),
			Version:     user.Version,
		}
	}), nil

}

// Count implements domain.UserRepo.
func (up *userRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&UserEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.UserRepo.
func (up *userRepo) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
	var userEntity = &UserEntity{
		Username:    user.Username,
		Password:    user.Password,
		Salt:        user.Salt,
		Mobile:      user.Mobile,
		Avatar:      user.Avatar,
		Description: user.Description,
		Status:      int32(enum.EnableStatusEnabled),
		RoleId:      user.RoleId,
		RoleName:    user.RoleName,
		TenantId:    0,
		Version:     0,
	}
	err := up.repo.GetDBFromContext(ctx).Save(userEntity).Error
	if err != nil {
		return nil, err
	}
	user.Id = int32(userEntity.Id)
	return user, nil

}

// Update implements domain.UserRepo.
func (up *userRepo) Update(ctx context.Context, user *domain.User) (int32, error) {

	var userEntity = &UserEntity{
		Id:          user.Id,
		Username:    user.Username,
		Password:    "",
		Mobile:      user.Mobile,
		Avatar:      user.Avatar,
		Description: user.Description,
		Status:      user.Status,
		RoleId:      user.RoleId,
		RoleName:    user.RoleName,
		TenantId:    0,
		Version:     user.Version,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&UserEntity{Id: user.Id}).Updates(userEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewUserRepo .
func NewUserRepo(repo *BaseRepo, logger log.Logger) domain.UserRepo {
	return &userRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
