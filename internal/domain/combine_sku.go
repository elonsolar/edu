package domain

import (
	"github.com/go-kratos/kratos/v2/log"
)

type CombineSku struct {
	Id    int32
	Code  string
	Name  string
	Price float64

	Description      string
	Status           int32 //已下架，在售
	OccupiedQuantity int32
	UpdatedAt        string
	Version          int32
}

type CombineSkuRepo interface {
	BaseRepo[CombineSku]
}

type CombineSkuService struct {
	BaseService[CombineSku]
}

func NewCombineSkuService(log log.Logger, repo CombineSkuRepo) *CombineSkuService {

	return &CombineSkuService{
		BaseService: BaseService[CombineSku]{
			repo: repo,
		},
	}
}
