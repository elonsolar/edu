package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/util"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type SaleOrderItemEntity struct {
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

func (t *SaleOrderItemEntity) TableName() string {
	return "combine_sku_item"

}

type saleOrderItemRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.SaleOrderItemRepo.
func (cr *saleOrderItemRepo) BatchSave(ctx context.Context, itemList []*domain.SaleOrderItem) error {

	var itemEntityList = make([]*SaleOrderItemEntity, len(itemList))
	_ = util.CopyProperties(&itemList, &itemEntityList, util.IgnoreNotMatchedProperty())
	err := cr.repo.GetDBFromContext(ctx).Save(&itemEntityList).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete implements domain.SaleOrderItemRepo.
func (cr *saleOrderItemRepo) Delete(ctx context.Context, id int32) (int32, error) {
	result := cr.repo.GetDBFromContext(ctx).Delete(&SaleOrderItemEntity{}, id)
	if result.Error != nil {
		return 0, result.Error
	}
	return int32(result.RowsAffected), nil
}

// ListAll implements domain.SaleOrderItemRepo.
func (*saleOrderItemRepo) ListAll(context.Context, *domain.Expression) ([]*domain.SaleOrderItem, error) {
	panic("unimplemented")
}

// ListByMap implements domain.SaleOrderItemRepo.
func (pr *saleOrderItemRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.SaleOrderItem, error) {
	var skuList = make([]*SaleOrderItemEntity, 0)
	err := pr.repo.GetDBFromContext(ctx).Find(&skuList, params).Error
	if err != nil {
		return nil, err
	}
	var result = make([]*domain.SaleOrderItem, 0, len(skuList))
	_ = util.CopyProperties(&skuList, &result, util.IgnoreNotMatchedProperty())
	return result, nil
}

// FindByID implements domain.SaleOrderItemRepo.
func (up *saleOrderItemRepo) FindByID(ctx context.Context, id int32) (*domain.SaleOrderItem, error) {

	var sku SaleOrderItemEntity

	err := up.repo.GetDBFromContext(ctx).First(&sku, id).Error

	if err != nil {
		return nil, err
	}
	var result domain.SaleOrderItem
	_ = util.CopyProperties(&sku, &result, util.IgnoreNotMatchedProperty())

	return &result, nil
}

// List implements domain.SaleOrderItemRepo.
func (up *saleOrderItemRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.SaleOrderItem, error) {

	var skuList []*SaleOrderItemEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&skuList).Error
	if err != nil {
		return nil, err
	}

	var result = make([]*domain.SaleOrderItem, 0, len(skuList))
	_ = util.CopyProperties(&skuList, &result, util.IgnoreNotMatchedProperty())

	return result, nil

}

// Count implements domain.SaleOrderItemRepo.
func (up *saleOrderItemRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&SaleOrderItemEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.SaleOrderItemReppo.
func (up *saleOrderItemRepo) Save(ctx context.Context, sku *domain.SaleOrderItem) (*domain.SaleOrderItem, error) {
	var combineSkuEntity SaleOrderItemEntity

	_ = util.CopyProperties(sku, &combineSkuEntity, util.IgnoreNotMatchedProperty())

	err := up.repo.GetDBFromContext(ctx).Create(&combineSkuEntity).Error

	if err != nil {
		return nil, err
	}
	sku.Id = int32(combineSkuEntity.Id)
	return sku, nil

}

// Update implements domain.SaleOrderItemRepo.
func (up *saleOrderItemRepo) Update(ctx context.Context, sku *domain.SaleOrderItem) (int32, error) {

	var combineSkuEntity SaleOrderItemEntity

	_ = util.CopyProperties(sku, &combineSkuEntity, util.IgnoreNotMatchedProperty())

	result := up.repo.GetDBFromContext(ctx).Model(&SaleOrderItemEntity{Id: sku.Id}).Updates(&combineSkuEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewSaleOrderItemRepo .
func NewSaleOrderItemRepo(repo *BaseRepo, logger log.Logger) domain.SaleOrderItemRepo {
	return &saleOrderItemRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
