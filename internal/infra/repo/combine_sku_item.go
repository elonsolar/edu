package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/util"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type CombineSkuItemEntity struct {
	Id          int32 `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	CombineId   int32          //组合商品ID
	SkuId       int32          //商品ID
	SkuCode     string
	SkuName     string
	SkuPrice    float64
	SkuQuantity int32 // 数量
	UserId      int   // 和登录账号绑定
	TenantId    int
	Version     int32 `gorm:"default:1"`
}

func (t *CombineSkuItemEntity) TableName() string {
	return "combine_sku_item"

}

type combineSkuItemRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.CombineSkuItemRepo.
func (cr *combineSkuItemRepo) BatchSave(ctx context.Context, itemList []*domain.CombineSkuItem) error {

	var itemEntityList = make([]*CombineSkuItemEntity, len(itemList))
	_ = util.CopyProperties(&itemList, &itemEntityList, util.IgnoreNotMatchedProperty())
	err := cr.repo.GetDBFromContext(ctx).Save(&itemEntityList).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete implements domain.CombineSkuItemRepo.
func (cr *combineSkuItemRepo) Delete(ctx context.Context, id int32) (int32, error) {
	result := cr.repo.GetDBFromContext(ctx).Delete(&CombineSkuItemEntity{}, id)
	if result.Error != nil {
		return 0, result.Error
	}
	return int32(result.RowsAffected), nil
}

// ListAll implements domain.CombineSkuItemRepo.
func (*combineSkuItemRepo) ListAll(context.Context, *domain.Expression) ([]*domain.CombineSkuItem, error) {
	panic("unimplemented")
}

// ListByMap implements domain.CombineSkuItemRepo.
func (pr *combineSkuItemRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.CombineSkuItem, error) {
	var skuList = make([]*CombineSkuItemEntity, 0)
	err := pr.repo.GetDBFromContext(ctx).Find(&skuList, params).Error
	if err != nil {
		return nil, err
	}
	var result = make([]*domain.CombineSkuItem, 0, len(skuList))
	_ = util.CopyProperties(&skuList, &result, util.IgnoreNotMatchedProperty())
	return result, nil
}

// FindByID implements domain.CombineSkuItemRepo.
func (up *combineSkuItemRepo) FindByID(ctx context.Context, id int32) (*domain.CombineSkuItem, error) {

	var sku CombineSkuItemEntity

	err := up.repo.GetDBFromContext(ctx).First(&sku, id).Error

	if err != nil {
		return nil, err
	}
	var result domain.CombineSkuItem
	_ = util.CopyProperties(&sku, &result, util.IgnoreNotMatchedProperty())

	return &result, nil
}

// List implements domain.CombineSkuItemRepo.
func (up *combineSkuItemRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.CombineSkuItem, error) {

	var skuList []*CombineSkuItemEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&skuList).Error
	if err != nil {
		return nil, err
	}

	var result = make([]*domain.CombineSkuItem, 0, len(skuList))
	_ = util.CopyProperties(&skuList, &result, util.IgnoreNotMatchedProperty())

	return result, nil

}

// Count implements domain.CombineSkuItemRepo.
func (up *combineSkuItemRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&CombineSkuItemEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.CombineSkuItemReppo.
func (up *combineSkuItemRepo) Save(ctx context.Context, sku *domain.CombineSkuItem) (*domain.CombineSkuItem, error) {
	var combineSkuEntity CombineSkuItemEntity

	_ = util.CopyProperties(sku, &combineSkuEntity, util.IgnoreNotMatchedProperty())

	err := up.repo.GetDBFromContext(ctx).Create(&combineSkuEntity).Error

	if err != nil {
		return nil, err
	}
	sku.Id = int32(combineSkuEntity.Id)
	return sku, nil

}

// Update implements domain.CombineSkuItemRepo.
func (up *combineSkuItemRepo) Update(ctx context.Context, sku *domain.CombineSkuItem) (int32, error) {

	var combineSkuEntity CombineSkuItemEntity

	_ = util.CopyProperties(sku, &combineSkuEntity, util.IgnoreNotMatchedProperty())

	result := up.repo.GetDBFromContext(ctx).Model(&CombineSkuItemEntity{Id: sku.Id}).Updates(&combineSkuEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewCombineSkuItemRepo .
func NewCombineSkuItemRepo(repo *BaseRepo, logger log.Logger) domain.CombineSkuItemRepo {
	return &combineSkuItemRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
