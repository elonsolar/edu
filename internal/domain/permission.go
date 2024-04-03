package domain

import (
	"context"
	"edu/internal/domain/enum"
	"edu/internal/util"
)

type PermissionRepo interface {
	BaseRepo[Permission]
}

type Permission struct {
	Id   int32
	Code string
	Name string

	PermissionType int32

	description string
	ParentId    int32

	UpdatedAt string
	Status    int32
	Version   int32
	Children  []*Permission
}

type PermissionService struct {
	BaseService[Permission]
	permissionRepo PermissionRepo
}

func (ps *PermissionService) GetTree(ctx context.Context, parentId int32, permissionType int32) ([]*Permission, error) {

	var params = map[string]interface{}{}
	if permissionType == int32(enum.PermissionType_MENU) || permissionType == int32(enum.PermissionType_ACTION) {
		params = map[string]interface{}{"permission_type": permissionType}
	}

	permissionList, err := ps.ListByMap(ctx, params)
	if err != nil {
		return nil, err
	}

	return util.NewTreeBuilder[Permission, int32]().Build(permissionList, parentId), nil
}

func (ps *PermissionService) GetTenantPermissionTree(ctx context.Context, parentId int32, permissionType int32, permissionIds []int32) ([]*Permission, error) {

	var params = map[string]interface{}{}
	if permissionType == int32(enum.PermissionType_MENU) || permissionType == int32(enum.PermissionType_ACTION) {
		params = map[string]interface{}{"permission_type": permissionType}
	}
	params["id"] = permissionIds

	permissionList, err := ps.ListByMap(ctx, params)
	if err != nil {
		return nil, err
	}

	return util.NewTreeBuilder[Permission, int32]().Build(permissionList, parentId), nil
}

func NewPermissionService(repo PermissionRepo) *PermissionService {

	return &PermissionService{
		BaseService: BaseService[Permission]{
			repo: repo,
		},
		permissionRepo: repo,
	}
}
