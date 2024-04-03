package service

import (
	"context"
	pb "edu/api/admin/v1"
	"edu/internal/domain"
	"edu/internal/domain/enum"
	"edu/internal/domain/model"
	"edu/internal/util"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwt4 "github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	pb.UnimplementedAuthServer
	tx                domain.Tx
	userService       *domain.UserService
	permissionService *domain.PermissionService
	roleService       *domain.RoleService
	tenantService     *domain.TenantService
}

func NewAuthService(tx domain.Tx, uc *domain.UserService, permissionService *domain.PermissionService, roleService *domain.RoleService, tenantService *domain.TenantService) *AuthService {
	return &AuthService{
		tx:                tx,
		userService:       uc,
		permissionService: permissionService,
		roleService:       roleService,
		tenantService:     tenantService,
	}
}

func (auth *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {

	// todo 加密,盐值
	user, err := auth.userService.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	var reply pb.LoginReply
	_ = util.CopyProperties(user, &reply, util.IgnoreNotMatchedProperty())
	mySigningKey := []byte("edu")

	claims := &model.Claims{
		Id:           user.Id,
		Username:     user.Username,
		TenantId:     user.TenantId,
		IsSuperAdmin: user.IsSuperAdmin,
		RegisteredClaims: jwt4.RegisteredClaims{
			Issuer:    "cxq",
			ExpiresAt: jwt4.NewNumericDate(time.Unix(1702309131, 0))},
	}

	token := jwt4.NewWithClaims(jwt4.SigningMethodHS256, claims)
	sign, err := token.SignedString(mySigningKey)
	if err != nil {
		return nil, err
	}

	reply.Token = sign

	return &reply, nil
}

func (auth *AuthService) GetUserInfo(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
	claims, ok := jwt.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("未授权")
	}
	c, ok := (claims).(*model.Claims)

	if !ok {
		return nil, fmt.Errorf("claims 类型错误")
	}

	// getUser
	user, err := auth.userService.FindByID(ctx, c.Id)
	if err != nil {
		return nil, fmt.Errorf("用户查找错误, %w", err)
	}

	// getPermission
	permissionIds, err := auth.roleService.GetPermissions(ctx, user.RoleId)
	if err != nil {
		return nil, fmt.Errorf("查找用户权限错误 %w", err)
	}
	permissions, err := auth.permissionService.ListByMap(ctx, map[string]interface{}{"id": permissionIds})
	if err != nil {
		return nil, err
	}

	var reply pb.UserInfoReply
	reply.Id = user.Id
	reply.Username = user.Username
	reply.Role = &pb.UserInfoReply_Role{
		Id:          user.RoleId,
		Name:        user.RoleName,
		Permissions: []*pb.UserInfoReply_Role_Permission{},
	}
	reply.Role.Permissions = []*pb.UserInfoReply_Role_Permission{}
	for _, permission := range permissions {
		if permission.PermissionType != int32(enum.PermissionType_MENU) {
			continue
		}
		pm := pb.UserInfoReply_Role_Permission{
			PermissionId:    permission.Code,
			PermissionName:  permission.Name,
			ActionEntitySet: []*pb.UserInfoReply_Role_Permission_Action{},
		}
		for _, sub := range permissions {
			if sub.ParentId == permission.Id && sub.PermissionType == int32(enum.PermissionType_ACTION) {
				pm.ActionEntitySet = append(pm.ActionEntitySet, &pb.UserInfoReply_Role_Permission_Action{
					Action:       sub.Code,
					Describe:     sub.Name,
					DefaultCheck: false,
				})

			}
		}
		reply.Role.Permissions = append(reply.Role.Permissions, &pm)

	}

	return &reply, nil

}

func (s *AuthService) CreatePermission(ctx context.Context, req *pb.CreatePermissionRequest) (*pb.CreatePermissionReply, error) {

	var permission domain.Permission
	_ = util.CopyProperties(req, &permission, util.IgnoreNotMatchedProperty())

	_, err := s.permissionService.Create(ctx, &permission)

	if err != nil {
		return nil, err
	}
	return &pb.CreatePermissionReply{}, nil
}
func (s *AuthService) UpdatePermission(ctx context.Context, req *pb.UpdatePermissionRequest) (*pb.UpdatePermissionReply, error) {

	var permission domain.Permission
	_ = util.CopyProperties(req, &permission, util.IgnoreNotMatchedProperty())

	err := s.permissionService.UpdateConcurrency(ctx, &permission)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePermissionReply{}, nil
}
func (s *AuthService) DeletePermission(ctx context.Context, req *pb.DeletePermissionRequest) (*pb.DeletePermissionReply, error) {
	return &pb.DeletePermissionReply{}, nil
}
func (s *AuthService) GetPermission(ctx context.Context, req *pb.GetPermissionRequest) (*pb.GetPermissionReply, error) {

	room, err := s.permissionService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var reply pb.GetPermissionReply

	_ = util.CopyProperties(room, &reply, util.IgnoreNotMatchedProperty())
	return &reply, nil
}
func (s *AuthService) ListPermission(ctx context.Context, req *pb.ListPermissionRequest) (*pb.ListPermissionReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.permissionService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var permissionList = []*pb.ListPermissionReply_Data{}
	for _, permission := range list {

		var replyData pb.ListPermissionReply_Data
		_ = util.CopyProperties(permission, &replyData, util.IgnoreNotMatchedProperty())
		permissionList = append(permissionList, &replyData)
	}
	return &pb.ListPermissionReply{
		Data:  permissionList,
		Total: int32(count),
	}, nil
}

func (s *AuthService) GetPermissionTree(ctx context.Context, req *pb.GetPermissionTreeRequest) (*pb.GetPermissionTreeReply, error) {

	list, err := s.permissionService.GetTree(ctx, req.ParentId, req.PermissionType)
	if err != nil {
		return nil, err
	}

	var reply pb.GetPermissionTreeReply
	_ = util.CopyProperties(&list, &(reply.Data), util.IgnoreNotMatchedProperty())

	return &reply, nil
}

// role
func (s *AuthService) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleReply, error) {

	var role domain.Role
	_ = util.CopyProperties(req, &role, util.IgnoreNotMatchedProperty())

	_, err := s.roleService.Create(ctx, &role)

	if err != nil {
		return nil, err
	}
	return &pb.CreateRoleReply{}, nil
}
func (s *AuthService) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.UpdateRoleReply, error) {

	var role domain.Role
	_ = util.CopyProperties(req, &role, util.IgnoreNotMatchedProperty())

	err := s.roleService.UpdateConcurrency(ctx, &role)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateRoleReply{}, nil
}
func (s *AuthService) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*pb.DeleteRoleReply, error) {
	return &pb.DeleteRoleReply{}, nil
}
func (s *AuthService) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleReply, error) {

	room, err := s.roleService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var reply pb.GetRoleReply

	_ = util.CopyProperties(room, &reply, util.IgnoreNotMatchedProperty())
	return &reply, nil
}

func (s *AuthService) ListRole(ctx context.Context, req *pb.ListRoleRequest) (*pb.ListRoleReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.roleService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var roleList = []*pb.ListRoleReply_Data{}
	for _, role := range list {

		var replyData pb.ListRoleReply_Data
		_ = util.CopyProperties(role, &replyData, util.IgnoreNotMatchedProperty())
		roleList = append(roleList, &replyData)
	}
	return &pb.ListRoleReply{
		Data:  roleList,
		Total: int32(count),
	}, nil
}

func (s *AuthService) GetRolePermission(ctx context.Context, req *pb.GetRolePermissionRequest) (*pb.GetRolePermissionReply, error) {

	list, err := s.roleService.GetPermissions(ctx, req.RoleId)
	if err != nil {
		return nil, err
	}

	var reply pb.GetRolePermissionReply
	_ = util.CopyProperties(&list, &(reply.Permissions), util.IgnoreNotMatchedProperty())

	return &reply, nil
}

func (s *AuthService) SaveRolePermission(ctx context.Context, req *pb.SaveRolePermissionRequest) (*pb.SaveRolePermissionReply, error) {

	if len(req.Permissions) == 0 {
		return &pb.SaveRolePermissionReply{}, nil
	}
	err := s.tx.Transaction(ctx, func(txctx context.Context) error {

		err := s.roleService.SavePermissions(ctx, req.RoleId, req.Permissions)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &pb.SaveRolePermissionReply{}, nil
}
