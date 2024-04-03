package domain

import (
	"github.com/go-kratos/kratos/v2/log"
)

type CombineSkuItem struct {
	Id          int32
	CombineId   int32 //组合商品ID
	SkuId       int32
	SkuCode     string
	SkuName     string
	SkuPrice    float64
	SkuQuantity int32 // 数量
	UpdatedAt   string
	Version     int32
}

type CombineSkuItemRepo interface {
	BaseRepo[CombineSkuItem]
}

type CombineSkuItemService struct {
	BaseService[CombineSkuItem]
}

func NewCombineSkuItemService(log log.Logger, repo CombineSkuItemRepo) *CombineSkuItemService {

	return &CombineSkuItemService{
		BaseService: BaseService[CombineSkuItem]{
			repo: repo,
		},
	}
}
