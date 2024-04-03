package domain

import (
	"context"
	"edu/internal/util"

	"github.com/go-kratos/kratos/v2/log"
)

type TenantRepo interface {
	BaseRepo[Tenant]
}

type Tenant struct {
	Id          int32
	Name        string
	Mobile      string
	Description string

	Status    int32
	UpdatedAt string
	Version   int32
}

type TenantPermissionRepo interface {
	BaseRepo[TenantPermission]
	BatchDelete(context.Context, int32, []int32) error
}

type TenantPermission struct {
	Id           int32
	TenantId     int32
	PermissionId int32
}

// TenantService is a Greeter usecase.
type TenantService struct {
	BaseService[Tenant]
	log                  *log.Helper
	tenantPermissionRepo TenantPermissionRepo
}

func (ps *TenantService) GetPermissions(ctx context.Context, tenantId int32) ([]int32, error) {

	permissionList, err := ps.tenantPermissionRepo.ListByMap(ctx, map[string]interface{}{"tenant_id": tenantId})
	if err != nil {
		return nil, err
	}

	permissionIds := util.MapPtr[TenantPermission, int32](permissionList, func(rp *TenantPermission) int32 { return rp.PermissionId })
	return permissionIds, nil
}

func (ps *TenantService) SavePermissions(ctx context.Context, tenantId int32, newPermissions []int32) ([]int32, []int32, error) {

	permissionList, err := ps.tenantPermissionRepo.ListByMap(ctx, map[string]interface{}{"tenant_id": tenantId})
	if err != nil {
		return nil, nil, err
	}

	oldPermissionIds := util.MapPtr[TenantPermission, int32](permissionList, func(rp *TenantPermission) int32 { return rp.PermissionId })

	newItemList, oldItemList := util.Difference[int32](newPermissions, oldPermissionIds)

	if len(newItemList) != 0 {
		newTenantPermissionList := util.Map[int32, *TenantPermission](newItemList, func(i int32) *TenantPermission { return &TenantPermission{TenantId: tenantId, PermissionId: i} })

		// insert
		err = ps.tenantPermissionRepo.BatchSave(ctx, newTenantPermissionList)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(oldItemList) != 0 {

		// delete
		err = ps.tenantPermissionRepo.BatchDelete(ctx, tenantId, oldItemList)
		if err != nil {
			return nil, nil, err
		}
	}

	return newItemList, oldItemList, nil

}

// NewTenantService new a Tenant usecase.
func NewTenantService(repo TenantRepo, tenantPermissionRepo TenantPermissionRepo, logger log.Logger) *TenantService {
	return &TenantService{
		BaseService: BaseService[Tenant]{repo: repo}, log: log.NewHelper(logger),
		tenantPermissionRepo: tenantPermissionRepo,
	}
}
