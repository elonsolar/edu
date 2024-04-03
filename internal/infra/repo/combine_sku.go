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

type CombineSkuEntity struct {
	Id          int32 `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Code        string
	Name        string
	Price       float64
	Description string
	Status      int32 //已下架，在售

	OccupiedQuantity int32
	UserId           int // 和登录账号绑定
	TenantId         int
	Version          int32 `gorm:"default:1"`
}

func (t *CombineSkuEntity) TableName() string {
	return "combine_sku"

}

type combineSkuRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.CombineSkuRepo.
func (*combineSkuRepo) BatchSave(context.Context, []*domain.CombineSku) error {
	panic("unimplemented")
}

// Delete implements domain.CombineSkuRepo.
func (*combineSkuRepo) Delete(context.Context, int32) (int32, error) {
	panic("unimplemented")
}

// ListAll implements domain.CombineSkuRepo.
func (*combineSkuRepo) ListAll(context.Context, *domain.Expression) ([]*domain.CombineSku, error) {
	panic("unimplemented")
}

// ListByMap implements domain.CombineSkuRepo.
func (pr *combineSkuRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.CombineSku, error) {
	var skuList = make([]*CombineSkuEntity, 0)
	err := pr.repo.GetDBFromContext(ctx).Find(&skuList, params).Error
	if err != nil {
		return nil, err
	}
	var result = make([]*domain.CombineSku, 0, len(skuList))
	_ = util.CopyProperties(&skuList, &result, util.IgnoreNotMatchedProperty())
	return result, nil
}

// FindByID implements domain.CombineSkuRepo.
func (up *combineSkuRepo) FindByID(ctx context.Context, id int32) (*domain.CombineSku, error) {

	var sku CombineSkuEntity

	err := up.repo.GetDBFromContext(ctx).First(&sku, id).Error

	if err != nil {
		return nil, err
	}
	var result domain.CombineSku
	_ = util.CopyProperties(&sku, &result, util.IgnoreNotMatchedProperty())

	return &result, nil
}

// List implements domain.CombineSkuRepo.
func (up *combineSkuRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.CombineSku, error) {

	var skuList []*CombineSkuEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&skuList).Error
	if err != nil {
		return nil, err
	}

	var result = make([]*domain.CombineSku, 0, len(skuList))
	_ = util.CopyProperties(&skuList, &result, util.IgnoreNotMatchedProperty())

	return result, nil

}

// Count implements domain.CombineSkuRepo.
func (up *combineSkuRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&CombineSkuEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.CombineSkuReppo.
func (up *combineSkuRepo) Save(ctx context.Context, sku *domain.CombineSku) (*domain.CombineSku, error) {
	var combineSkuEntity CombineSkuEntity

	_ = util.CopyProperties(sku, &combineSkuEntity, util.IgnoreNotMatchedProperty())

	combineSkuEntity.Status = int32(enum.SkuStatusType_AVAILABLE)
	err := up.repo.GetDBFromContext(ctx).Create(&combineSkuEntity).Error

	if err != nil {
		return nil, err
	}
	sku.Id = int32(combineSkuEntity.Id)
	return sku, nil

}

// Update implements domain.CombineSkuRepo.
func (up *combineSkuRepo) Update(ctx context.Context, sku *domain.CombineSku) (int32, error) {

	var combineSkuEntity CombineSkuEntity

	_ = util.CopyProperties(sku, &combineSkuEntity, util.IgnoreNotMatchedProperty())

	result := up.repo.GetDBFromContext(ctx).Model(&CombineSkuEntity{Id: sku.Id}).Updates(&combineSkuEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewCombineSkuRepo .
func NewCombineSkuRepo(repo *BaseRepo, logger log.Logger) domain.CombineSkuRepo {
	return &combineSkuRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
