package domain

import (
	"github.com/go-kratos/kratos/v2/log"
)

type SaleOrder struct {
	Id            int32
	SerialNumber  string // 订单流水号
	CustomerId    string // 客户
	CustomerName  string
	CustomerPhone string

	SalespersonId    string // 销售人员
	SalespersonName  string
	SalespersonPhone string

	TotalPrice    float64 // 总价
	DiscountPrice float64 // 折扣价
	Description   string
	UpdatedAt     string
	Version       int32
}

type SaleOrderRepo interface {
	BaseRepo[SaleOrder]
}

type SaleOrderService struct {
	BaseService[SaleOrder]
}

func NewSaleOrderService(log log.Logger, repo SaleOrderRepo) *SaleOrderService {

	return &SaleOrderService{
		BaseService: BaseService[SaleOrder]{
			repo: repo,
		},
	}
}
