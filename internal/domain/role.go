package domain

import (
	"context"
	"edu/internal/util"
)

type RoleRepo interface {
	BaseRepo[Role]
}

type RolePermissionRepo interface {
	BaseRepo[RolePermission]
	DeleteRolePermissions(context.Context, int32, []int32) error
	DeleteTenantPermissions(context.Context, int32, []int32) error
}

type Role struct {
	Id          int32
	Code        string
	Name        string
	Description string

	UpdatedAt string
	Status    int32
	Version   int32
}

type RolePermission struct {
	Id           int32
	RoleId       int32
	PermissionId int32
}

type RoleService struct {
	BaseService[Role]
	roleRepo           RoleRepo
	rolePermissionRepo RolePermissionRepo
}

func (ps *RoleService) GetPermissions(ctx context.Context, roldId int32) ([]int32, error) {

	permissionList, err := ps.rolePermissionRepo.ListByMap(ctx, map[string]interface{}{"role_id": roldId})
	if err != nil {
		return nil, err
	}

	permissionIds := util.MapPtr[RolePermission, int32](permissionList, func(rp *RolePermission) int32 { return rp.PermissionId })
	return permissionIds, nil
}

func (ps *RoleService) SavePermissions(ctx context.Context, roldId int32, newPermissions []int32) error {

	permissionList, err := ps.rolePermissionRepo.ListByMap(ctx, map[string]interface{}{"role_id": roldId})
	if err != nil {
		return err
	}

	oldPermissionIds := util.MapPtr[RolePermission, int32](permissionList, func(rp *RolePermission) int32 { return rp.PermissionId })

	newItemList, oldItemList := util.Difference[int32](newPermissions, oldPermissionIds)

	if len(newItemList) != 0 {
		newRolePermissionList := util.Map[int32, *RolePermission](newItemList, func(i int32) *RolePermission { return &RolePermission{RoleId: roldId, PermissionId: i} })

		// insert
		err = ps.rolePermissionRepo.BatchSave(ctx, newRolePermissionList)
		if err != nil {
			return err
		}
	}

	if len(oldItemList) != 0 {

		// delete
		err = ps.rolePermissionRepo.DeleteRolePermissions(ctx, roldId, oldItemList)
		if err != nil {
			return err
		}
	}

	return nil

}

func (rs *RoleService) DeleteTenantPermissions(ctx context.Context, tenantId int32, permissions []int32) error {
	return rs.rolePermissionRepo.DeleteTenantPermissions(ctx, tenantId, permissions)
}

func NewRoleService(roleRepo RoleRepo, rolePermissionRepo RolePermissionRepo) *RoleService {

	return &RoleService{
		BaseService: BaseService[Role]{
			repo: roleRepo,
		},
		roleRepo:           roleRepo,
		rolePermissionRepo: rolePermissionRepo,
	}
}
