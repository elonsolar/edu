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

type SaleOrderEntity struct {
	Id            int32 `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	SerialNumber  string         // 订单流水号
	CustomerId    string         // 客户
	CustomerName  string
	CustomerPhone string

	SalespersonId    string // 销售人员
	SalespersonName  string
	SalespersonPhone string

	TotalPrice    float64 // 总价
	DiscountPrice float64 // 折扣价
	Description   string
	Status        int32

	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32 `gorm:"default:1"`
}

func (t *SaleOrderEntity) TableName() string {
	return "combine_sku"

}

type saleOrderRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.SaleOrderRepo.
func (*saleOrderRepo) BatchSave(context.Context, []*domain.SaleOrder) error {
	panic("unimplemented")
}

// Delete implements domain.SaleOrderRepo.
func (*saleOrderRepo) Delete(context.Context, int32) (int32, error) {
	panic("unimplemented")
}

// ListAll implements domain.SaleOrderRepo.
func (*saleOrderRepo) ListAll(context.Context, *domain.Expression) ([]*domain.SaleOrder, error) {
	panic("unimplemented")
}

// ListByMap implements domain.SaleOrderRepo.
func (pr *saleOrderRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.SaleOrder, error) {
	var skuList = make([]*SaleOrderEntity, 0)
	err := pr.repo.GetDBFromContext(ctx).Find(&skuList, params).Error
	if err != nil {
		return nil, err
	}
	var result = make([]*domain.SaleOrder, 0, len(skuList))
	_ = util.CopyProperties(&skuList, &result, util.IgnoreNotMatchedProperty())
	return result, nil
}

// FindByID implements domain.SaleOrderRepo.
func (up *saleOrderRepo) FindByID(ctx context.Context, id int32) (*domain.SaleOrder, error) {

	var sku SaleOrderEntity

	err := up.repo.GetDBFromContext(ctx).First(&sku, id).Error

	if err != nil {
		return nil, err
	}
	var result domain.SaleOrder
	_ = util.CopyProperties(&sku, &result, util.IgnoreNotMatchedProperty())

	return &result, nil
}

// List implements domain.SaleOrderRepo.
func (up *saleOrderRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.SaleOrder, error) {

	var skuList []*SaleOrderEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&skuList).Error
	if err != nil {
		return nil, err
	}

	var result = make([]*domain.SaleOrder, 0, len(skuList))
	_ = util.CopyProperties(&skuList, &result, util.IgnoreNotMatchedProperty())

	return result, nil

}

// Count implements domain.SaleOrderRepo.
func (up *saleOrderRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&SaleOrderEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.SaleOrderReppo.
func (up *saleOrderRepo) Save(ctx context.Context, sku *domain.SaleOrder) (*domain.SaleOrder, error) {
	var saleOrderEntity SaleOrderEntity

	_ = util.CopyProperties(sku, &saleOrderEntity, util.IgnoreNotMatchedProperty())

	saleOrderEntity.Status = int32(enum.SkuStatusType_AVAILABLE)
	err := up.repo.GetDBFromContext(ctx).Create(&saleOrderEntity).Error

	if err != nil {
		return nil, err
	}
	sku.Id = int32(saleOrderEntity.Id)
	return sku, nil

}

// Update implements domain.SaleOrderRepo.
func (up *saleOrderRepo) Update(ctx context.Context, sku *domain.SaleOrder) (int32, error) {

	var saleOrderEntity SaleOrderEntity

	_ = util.CopyProperties(sku, &saleOrderEntity, util.IgnoreNotMatchedProperty())

	result := up.repo.GetDBFromContext(ctx).Model(&SaleOrderEntity{Id: sku.Id}).Updates(&saleOrderEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewSaleOrderRepo .
func NewSaleOrderRepo(repo *BaseRepo, logger log.Logger) domain.SaleOrderRepo {
	return &saleOrderRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
