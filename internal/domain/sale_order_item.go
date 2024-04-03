package domain

import (
	"github.com/go-kratos/kratos/v2/log"
)

type SaleOrderItem struct {
	Id          int32
	OrderId     int32 //组合商品ID
	SkuType     int32 // 商品类型，单品和组合商品
	SkuId       int32
	SkuCode     string
	SkuName     string
	SkuPrice    float64
	SkuQuantity int32 // 数量
	UpdatedAt   string
	Version     int32
}

type SaleOrderItemRepo interface {
	BaseRepo[SaleOrderItem]
}

type SaleOrderItemService struct {
	BaseService[SaleOrderItem]
}

func NewSaleOrderItemService(log log.Logger, repo SaleOrderItemRepo) *SaleOrderItemService {

	return &SaleOrderItemService{
		BaseService: BaseService[SaleOrderItem]{
			repo: repo,
		},
	}
}
