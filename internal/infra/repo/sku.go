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

type SkuEntity struct {
	Id        int32 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Code      string
	Name      string
	Category  int32 // 课时，耗材, 组合商品

	Specifications   string //规格
	Price            float64
	OccupiedQuantity int32
	Quantity         int32  // 库存
	Unit             string //单位
	Description      string
	Status           int32 //已下架，在售
	UserId           int   // 和登录账号绑定
	TenantId         int
	Version          int32 `gorm:"default:1"`
}

func (t *SkuEntity) TableName() string {
	return "sku"

}

type skuRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.SkuRepo.
func (*skuRepo) BatchSave(context.Context, []*domain.Sku) error {
	panic("unimplemented")
}

// Delete implements domain.SkuRepo.
func (*skuRepo) Delete(context.Context, int32) (int32, error) {
	panic("unimplemented")
}

// ListAll implements domain.SkuRepo.
func (*skuRepo) ListAll(context.Context, *domain.Expression) ([]*domain.Sku, error) {
	panic("unimplemented")
}

// ListByMap implements domain.SkuRepo.
func (pr *skuRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.Sku, error) {
	var skuList = make([]*SkuEntity, 0)
	err := pr.repo.GetDBFromContext(ctx).Find(&skuList, params).Error
	if err != nil {
		return nil, err
	}
	var result = make([]*domain.Sku, 0, len(skuList))
	_ = util.CopyProperties(&skuList, &result, util.IgnoreNotMatchedProperty())
	return result, nil
}

// FindByID implements domain.SkuRepo.
func (up *skuRepo) FindByID(ctx context.Context, id int32) (*domain.Sku, error) {

	var sku SkuEntity

	err := up.repo.GetDBFromContext(ctx).First(&sku, id).Error

	if err != nil {
		return nil, err
	}
	var result domain.Sku
	_ = util.CopyProperties(&sku, &result, util.IgnoreNotMatchedProperty())

	return &result, nil
}

// List implements domain.SkuRepo.
func (up *skuRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.Sku, error) {

	var skuList []*SkuEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&skuList).Error
	if err != nil {
		return nil, err
	}

	var result = make([]*domain.Sku, 0, len(skuList))
	_ = util.CopyProperties(&skuList, &result, util.IgnoreNotMatchedProperty())

	return result, nil

}

// Count implements domain.SkuRepo.
func (up *skuRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&SkuEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.SkuReppo.
func (up *skuRepo) Save(ctx context.Context, sku *domain.Sku) (*domain.Sku, error) {
	var skuEntity SkuEntity

	_ = util.CopyProperties(sku, &skuEntity, util.IgnoreNotMatchedProperty())

	skuEntity.Status = int32(enum.SkuStatusType_AVAILABLE)
	err := up.repo.GetDBFromContext(ctx).Create(&skuEntity).Error

	if err != nil {
		return nil, err
	}
	sku.Id = int32(skuEntity.Id)
	return sku, nil

}

// Update implements domain.SkuRepo.
func (up *skuRepo) Update(ctx context.Context, sku *domain.Sku) (int32, error) {

	var skuEntity SkuEntity

	_ = util.CopyProperties(sku, &skuEntity, util.IgnoreNotMatchedProperty())

	result := up.repo.GetDBFromContext(ctx).Model(&SkuEntity{Id: sku.Id}).Updates(&skuEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewSkuRepo .
func NewSkuRepo(repo *BaseRepo, logger log.Logger) domain.SkuRepo {
	return &skuRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
