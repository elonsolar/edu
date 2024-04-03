package domain

import (
	"github.com/go-kratos/kratos/v2/log"
)

type Sku struct {
	Id               int32
	Code             string
	Name             string
	Category         int32  // 课时，教材，其他
	Specifications   string //规格
	Price            float64
	OccupiedQuantity int32
	Quantity         int32  // 库存
	Unit             string //单位

	Description string
	Status      int32 //已下架，在售
	UpdatedAt   string
	Version     int32
}

type SkuRepo interface {
	BaseRepo[Sku]
}

type SkuService struct {
	BaseService[Sku]
}

func NewSkuService(log log.Logger, repo SkuRepo) *SkuService {

	return &SkuService{
		BaseService: BaseService[Sku]{
			repo: repo,
		},
	}
}
