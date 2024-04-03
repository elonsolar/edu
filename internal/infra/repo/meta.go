package repo

import (
	"context"
	"edu/internal/domain"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type MetaEntity struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	DataType int32
	Name     string
	Value    string
	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32
}

func (t *MetaEntity) TableName() string {
	return "meta"

}

type metaRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.MetaRepo.
func (*metaRepo) BatchSave(context.Context, []*domain.Meta) error {
	panic("unimplemented")
}

// Delete implements domain.MetaRepo.
func (*metaRepo) Delete(context.Context, int32) (int32, error) {
	panic("unimplemented")
}

// ListAll implements domain.MetaRepo.
func (*metaRepo) ListAll(context.Context, *domain.Expression) ([]*domain.Meta, error) {
	panic("unimplemented")
}

// ListByMap implements domain.MetaRepo.
func (tp *metaRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.Meta, error) {
	var metaList = make([]*MetaEntity, 0)
	err := tp.repo.GetDBFromContext(ctx).Model(&MetaEntity{}).Find(&metaList, params).Error
	if err != nil {
		return nil, err
	}

	result := Map(metaList, func(meta *MetaEntity) *domain.Meta {
		return &domain.Meta{
			Id:       int32(meta.ID),
			DataType: meta.DataType,
			Name:     meta.Name,
			Value:    meta.Value,
			Version:  meta.Version,
		}
	})
	return result, nil
}

// FindByID implements domain.metaRepo.
func (up *metaRepo) FindByID(ctx context.Context, id int32) (*domain.Meta, error) {

	var meta MetaEntity

	err := up.repo.GetDBFromContext(ctx).First(&meta, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.Meta{
		Id:       int32(meta.ID),
		DataType: meta.DataType,
		Name:     meta.Name,
		Value:    meta.Value,
		Version:  meta.Version,
	}, nil
}

// List implements domain.metaRepo.
func (up *metaRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.Meta, error) {

	var metaList []*MetaEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&metaList).Error
	if err != nil {
		return nil, err
	}

	return Map(metaList, func(meta *MetaEntity) *domain.Meta {
		return &domain.Meta{
			Id:        int32(meta.ID),
			DataType:  meta.DataType,
			Name:      meta.Name,
			Value:     meta.Value,
			Version:   meta.Version,
			UpdatedAt: meta.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}), nil

}

// Count implements domain.metaRepo.
func (up *metaRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&MetaEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.metaRepo.
func (up *metaRepo) Save(ctx context.Context, meta *domain.Meta) (*domain.Meta, error) {
	var metaEntity = &MetaEntity{

		DataType: meta.DataType,
		Name:     meta.Name,
		Value:    meta.Value,
		Version:  meta.Version,
	}
	err := up.repo.GetDBFromContext(ctx).Save(metaEntity).Error
	if err != nil {
		return nil, err
	}
	meta.Id = int32(metaEntity.ID)
	return meta, nil

}

// Update implements domain.metaRepo.
func (up *metaRepo) Update(ctx context.Context, meta *domain.Meta) (int32, error) {
	var metaEntity = &MetaEntity{
		ID:       uint(meta.Id),
		DataType: meta.DataType,
		Name:     meta.Name,
		Value:    meta.Value,
		Version:  meta.Version,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&MetaEntity{ID: uint(metaEntity.ID)}).Updates(metaEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewmetaRepo .
func NewMetaRepo(repo *BaseRepo, logger log.Logger) domain.MetaRepo {
	return &metaRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
