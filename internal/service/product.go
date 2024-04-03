package service

import (
	"context"
	"encoding/json"

	pb "edu/api/product/v1"
	"edu/internal/domain"
	"edu/internal/domain/enum"
	"edu/internal/util"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type ProductService struct {
	pb.UnimplementedProductServer

	log                   *log.Helper
	tx                    domain.Tx
	skuService            *domain.SkuService
	combineSkuService     *domain.CombineSkuService
	combineSkuItemService *domain.CombineSkuItemService
}

func NewProductService(logger log.Logger, tx domain.Tx, skuService *domain.SkuService, combineSkuService *domain.CombineSkuService, combineSkuItemService *domain.CombineSkuItemService) *ProductService {
	return &ProductService{
		log:                   log.NewHelper(logger),
		tx:                    tx,
		skuService:            skuService,
		combineSkuService:     combineSkuService,
		combineSkuItemService: combineSkuItemService,
	}
}

// Sku
func (s *ProductService) CreateSku(ctx context.Context, req *pb.CreateSkuRequest) (*pb.CreateSkuReply, error) {

	var item domain.Sku
	_ = util.CopyProperties(req, &item, util.IgnoreNotMatchedProperty())

	_, err := s.skuService.Create(ctx, &item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.CreateSkuReply{}, nil
}
func (s *ProductService) UpdateSku(ctx context.Context, req *pb.UpdateSkuRequest) (*pb.UpdateSkuReply, error) {

	var item domain.Sku

	_ = util.CopyProperties(req, &item, util.IgnoreNotMatchedProperty())

	err := s.skuService.UpdateConcurrency(ctx, &item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.UpdateSkuReply{}, nil
}

func (s *ProductService) DeleteSku(ctx context.Context, req *pb.DeleteSkuRequest) (*pb.DeleteSkuReply, error) {

	var item domain.Sku

	_ = util.CopyProperties(req, &item, util.IgnoreNotMatchedProperty())

	_, err := s.skuService.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.DeleteSkuReply{}, nil
}

func (s *ProductService) GetSku(ctx context.Context, req *pb.GetSkuRequest) (*pb.GetSkuReply, error) {

	item, err := s.skuService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var reply pb.GetSkuReply
	_ = util.CopyProperties(item, &reply, util.IgnoreNotMatchedProperty())
	return &reply, nil
}

func (s *ProductService) ListSku(ctx context.Context, req *pb.ListSkuRequest) (*pb.ListSkuReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.skuService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var data = make([]*pb.ListSkuReply_Data, len(list))
	_ = util.CopyProperties(&list, &data, util.IgnoreNotMatchedProperty())

	return &pb.ListSkuReply{
		Data:  data,
		Total: int32(count),
	}, nil
}

func (s *ProductService) TakeDownSku(ctx context.Context, req *pb.TakeDownSkuRequest) (*pb.TakeDownSkuReply, error) {

	sku, err := s.skuService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if sku.Status != int32(enum.SkuStatusType_AVAILABLE) {
		return nil, status.Errorf(codes.Internal, "非在售状态商品无需下架")
	}

	err = s.skuService.UpdateConcurrency(ctx, &domain.Sku{Id: req.Id, Status: int32(enum.SkuStatusType_DISCONTINUED), Version: sku.Version})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.TakeDownSkuReply{}, nil
}

func (s *ProductService) PlaceUpSku(ctx context.Context, req *pb.PlaceUpSkuRequest) (*pb.PlaceUpSkuReply, error) {

	sku, err := s.skuService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if sku.Status != int32(enum.SkuStatusType_DISCONTINUED) {
		return nil, status.Errorf(codes.Internal, "非下架状态商品无需上架")
	}

	if sku.OccupiedQuantity != 0 {
		return nil, status.Errorf(codes.Internal, "存在未完结的销售订单,不可下架")
	}

	err = s.skuService.UpdateConcurrency(ctx, &domain.Sku{Id: req.Id, Status: int32(enum.SkuStatusType_AVAILABLE), Version: sku.Version})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.PlaceUpSkuReply{}, nil
}

// combineSku
func (s *ProductService) CreateCombineSku(ctx context.Context, req *pb.CreateCombineSkuRequest) (*pb.CreateCombineSkuReply, error) {

	var combineSku domain.CombineSku
	_ = util.CopyProperties(req, &combineSku, util.IgnoreNotMatchedProperty())

	var itemList = make([]*domain.CombineSkuItem, len(req.ItemList))

	_ = util.CopyProperties(&req.ItemList, &itemList, util.IgnoreNotMatchedProperty())

	err := s.tx.Transaction(ctx, func(txctx context.Context) error {
		newCombineSku, err := s.combineSkuService.Create(txctx, &combineSku)
		if err != nil {
			return err
		}
		for _, item := range itemList {
			item.CombineId = newCombineSku.Id
		}
		err = s.combineSkuItemService.BatchCreate(txctx, itemList)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.CreateCombineSkuReply{}, nil
}
func (s *ProductService) UpdateCombineSku(ctx context.Context, req *pb.UpdateCombineSkuRequest) (*pb.UpdateCombineSkuReply, error) {

	var combineSku domain.CombineSku

	_ = util.CopyProperties(req, &combineSku, util.IgnoreNotMatchedProperty())

	var itemList = make([]*domain.CombineSkuItem, len(req.ItemList))

	_ = util.CopyProperties(&req.ItemList, &itemList, util.IgnoreNotMatchedProperty())

	err := s.tx.Transaction(ctx, func(txctx context.Context) error {

		err := s.combineSkuService.UpdateConcurrency(ctx, &combineSku)
		if err != nil {
			return err
		}
		for _, item := range itemList {
			item.CombineId = combineSku.Id
		}
		err = s.combineSkuItemService.BatchCreate(txctx, itemList)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.UpdateCombineSkuReply{}, nil
}

func (s *ProductService) DeleteCombineSku(ctx context.Context, req *pb.DeleteCombineSkuRequest) (*pb.DeleteCombineSkuReply, error) {

	var item domain.CombineSku

	_ = util.CopyProperties(req, &item, util.IgnoreNotMatchedProperty())

	_, err := s.combineSkuService.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.DeleteCombineSkuReply{}, nil
}

func (s *ProductService) GetCombineSku(ctx context.Context, req *pb.GetCombineSkuRequest) (*pb.GetCombineSkuReply, error) {

	item, err := s.combineSkuService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var reply pb.GetCombineSkuReply
	_ = util.CopyProperties(item, &reply, util.IgnoreNotMatchedProperty())
	return &reply, nil
}

func (s *ProductService) ListCombineSku(ctx context.Context, req *pb.ListCombineSkuRequest) (*pb.ListCombineSkuReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.combineSkuService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var data = make([]*pb.ListCombineSkuReply_Data, len(list))
	_ = util.CopyProperties(&list, &data, util.IgnoreNotMatchedProperty())

	return &pb.ListCombineSkuReply{
		Data:  data,
		Total: int32(count),
	}, nil
}

func (s *ProductService) TakeDownCombineSku(ctx context.Context, req *pb.TakeDownCombineSkuRequest) (*pb.TakeDownCombineSkuReply, error) {

	sku, err := s.combineSkuService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if sku.Status != int32(enum.SkuStatusType_AVAILABLE) {
		return nil, status.Errorf(codes.Internal, "非在售状态商品无需下架")
	}
	if sku.OccupiedQuantity != 0 {
		return nil, status.Errorf(codes.Internal, "存在未完结的销售订单,不可下架")
	}

	err = s.combineSkuService.UpdateConcurrency(ctx, &domain.CombineSku{Id: req.Id, Status: int32(enum.SkuStatusType_DISCONTINUED), Version: sku.Version})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.TakeDownCombineSkuReply{}, nil
}

func (s *ProductService) PlaceUpCombineSku(ctx context.Context, req *pb.PlaceUpCombineSkuRequest) (*pb.PlaceUpCombineSkuReply, error) {

	sku, err := s.combineSkuService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if sku.Status != int32(enum.SkuStatusType_DISCONTINUED) {
		return nil, status.Errorf(codes.Internal, "非下架状态商品无需上架")
	}

	err = s.combineSkuService.UpdateConcurrency(ctx, &domain.CombineSku{Id: req.Id, Status: int32(enum.SkuStatusType_AVAILABLE), Version: sku.Version})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.PlaceUpCombineSkuReply{}, nil
}

// combineSkuItem
func (s *ProductService) CreateCombineSkuItem(ctx context.Context, req *pb.CreateCombineSkuItemRequest) (*pb.CreateCombineSkuItemReply, error) {

	var item domain.CombineSkuItem
	_ = util.CopyProperties(req, &item, util.IgnoreNotMatchedProperty())

	_, err := s.combineSkuItemService.Create(ctx, &item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.CreateCombineSkuItemReply{}, nil
}
func (s *ProductService) UpdateCombineSkuItem(ctx context.Context, req *pb.UpdateCombineSkuItemRequest) (*pb.UpdateCombineSkuItemReply, error) {

	var item domain.CombineSkuItem

	_ = util.CopyProperties(req, &item, util.IgnoreNotMatchedProperty())

	err := s.combineSkuItemService.UpdateConcurrency(ctx, &item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.UpdateCombineSkuItemReply{}, nil
}

func (s *ProductService) DeleteCombineSkuItem(ctx context.Context, req *pb.DeleteCombineSkuItemRequest) (*pb.DeleteCombineSkuItemReply, error) {

	var item domain.CombineSkuItem

	_ = util.CopyProperties(req, &item, util.IgnoreNotMatchedProperty())

	_, err := s.combineSkuItemService.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.DeleteCombineSkuItemReply{}, nil
}

func (s *ProductService) GetCombineSkuItem(ctx context.Context, req *pb.GetCombineSkuItemRequest) (*pb.GetCombineSkuItemReply, error) {

	item, err := s.combineSkuItemService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var reply pb.GetCombineSkuItemReply
	_ = util.CopyProperties(item, &reply, util.IgnoreNotMatchedProperty())
	return &reply, nil
}

func (s *ProductService) ListCombineSkuItem(ctx context.Context, req *pb.ListCombineSkuItemRequest) (*pb.ListCombineSkuItemReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.combineSkuItemService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var data = make([]*pb.ListCombineSkuItemReply_Data, len(list))
	_ = util.CopyProperties(&list, &data, util.IgnoreNotMatchedProperty())

	return &pb.ListCombineSkuItemReply{
		Data:  data,
		Total: int32(count),
	}, nil
}
