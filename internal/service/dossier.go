package service

import (
	"context"
	"encoding/json"
	"fmt"

	pb "edu/api/admin/v1"
	"edu/internal/domain"
	"edu/internal/domain/model"
	"edu/internal/util"

	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
)

type DossierService struct {
	pb.UnimplementedDossierServer

	tx                domain.Tx
	tenantService     *domain.TenantService
	userService       *domain.UserService
	teacherService    *domain.TeacherService
	studentService    *domain.StudentService
	roomService       *domain.RoomService
	subjectService    *domain.SubjectService
	permissionService *domain.PermissionService
	roleService       *domain.RoleService
}

func NewDossierService(tx domain.Tx, userService *domain.UserService, tenantService *domain.TenantService, teacherService *domain.TeacherService, studentService *domain.StudentService, roomService *domain.RoomService, subjectService *domain.SubjectService, permissionService *domain.PermissionService, roleService *domain.RoleService) *DossierService {
	return &DossierService{
		tx:                tx,
		userService:       userService,
		teacherService:    teacherService,
		studentService:    studentService,
		roomService:       roomService,
		subjectService:    subjectService,
		tenantService:     tenantService,
		permissionService: permissionService,
		roleService:       roleService,
	}
}

// tenant
func (s *DossierService) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest) (*pb.CreateTenantReply, error) {

	var tenant domain.Tenant
	_ = util.CopyProperties(req, &tenant, util.IgnoreNotMatchedProperty())

	err := s.tx.Transaction(ctx, func(txctx context.Context) error {

		_, err := s.tenantService.Create(txctx, &tenant)
		if err != nil {
			return nil
		}
		_, err = s.userService.Create(txctx, &domain.User{
			Username: tenant.Mobile,
			Password: "11111111",
			Mobile:   tenant.Mobile,
			IsTenant: true,
		})

		if err != nil {
			return nil
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &pb.CreateTenantReply{}, nil
}
func (s *DossierService) UpdateTenant(ctx context.Context, req *pb.UpdateTenantRequest) (*pb.UpdateTenantReply, error) {

	var tenant domain.Tenant
	_ = util.CopyProperties(req, &tenant, util.IgnoreNotMatchedProperty())

	err := s.tenantService.UpdateConcurrency(ctx, &tenant)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateTenantReply{}, nil
}
func (s *DossierService) DeleteTenant(ctx context.Context, req *pb.DeleteTenantRequest) (*pb.DeleteTenantReply, error) {
	return &pb.DeleteTenantReply{}, nil
}
func (s *DossierService) GetTenant(ctx context.Context, req *pb.GetTenantRequest) (*pb.GetTenantReply, error) {

	room, err := s.tenantService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var reply pb.GetTenantReply

	_ = util.CopyProperties(room, &reply, util.IgnoreNotMatchedProperty())
	return &reply, nil
}

func (s *DossierService) ListTenant(ctx context.Context, req *pb.ListTenantRequest) (*pb.ListTenantReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.tenantService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var tenantList = []*pb.ListTenantReply_Data{}
	for _, tenant := range list {

		var replyData pb.ListTenantReply_Data
		_ = util.CopyProperties(tenant, &replyData, util.IgnoreNotMatchedProperty())
		tenantList = append(tenantList, &replyData)
	}
	return &pb.ListTenantReply{
		Data:  tenantList,
		Total: int32(count),
	}, nil
}

func (s *DossierService) GetTenantPermissionTree(ctx context.Context, req *pb.GetTenantPermissionTreeRequest) (*pb.GetTenantPermissionTreeReply, error) {
	claims, ok := jwt.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("用户未登陆，没有权限")
	}
	c, ok := (claims).(*model.Claims)
	if !ok {
		return nil, fmt.Errorf("获取身份信息错误")
	}

	permisionIds, err := s.tenantService.GetPermissions(ctx, c.TenantId)
	if err != nil {
		return nil, err
	}

	list, err := s.permissionService.GetTenantPermissionTree(ctx, req.ParentId, req.PermissionType, permisionIds)
	if err != nil {
		return nil, err
	}

	var reply pb.GetTenantPermissionTreeReply
	_ = util.CopyProperties(&list, &(reply.Data), util.IgnoreNotMatchedProperty())

	return &reply, nil
}

func (s *DossierService) GetTenantPermission(ctx context.Context, req *pb.GetTenantPermissionRequest) (*pb.GetTenantPermissionReply, error) {

	list, err := s.tenantService.GetPermissions(ctx, req.TenantId)
	if err != nil {
		return nil, err
	}

	var reply pb.GetTenantPermissionReply
	_ = util.CopyProperties(&list, &(reply.Permissions), util.IgnoreNotMatchedProperty())

	return &reply, nil
}

func (s *DossierService) SaveTenantPermission(ctx context.Context, req *pb.SaveTenantPermissionRequest) (*pb.SaveTenantPermissionReply, error) {

	if len(req.Permissions) == 0 {
		return &pb.SaveTenantPermissionReply{}, nil
	}
	err := s.tx.Transaction(ctx, func(txctx context.Context) error {

		_, oldPermissionId, err := s.tenantService.SavePermissions(txctx, req.TenantId, req.Permissions)
		if err != nil {
			return err
		}
		err = s.roleService.DeleteTenantPermissions(txctx, req.TenantId, oldPermissionId)

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &pb.SaveTenantPermissionReply{}, nil
}

func (s *DossierService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {

	var user = &domain.User{}
	_ = util.CopyProperties(req, user, util.IgnoreNotMatchedProperty())

	_, err := s.userService.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserReply{}, nil
}
func (s *DossierService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {

	err := s.userService.UpdateConcurrency(ctx, &domain.User{
		Id:          req.Id,
		Username:    req.Username,
		Mobile:      req.Mobile,
		Description: req.Description,
		Avatar:      req.Mobile,
		Status:      req.Status,
		RoleId:      req.RoleId,
		RoleName:    req.RoleName,
		Version:     req.Version,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserReply{}, nil
}
func (s *DossierService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {

	return &pb.DeleteUserReply{}, nil
}
func (s *DossierService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {

	user, err := s.userService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserReply{
		Username:    user.Username,
		Mobile:      user.Mobile,
		Description: user.Description,
		Status:      int32(user.Status),
		Avatar:      user.Avatar,
		RoleId:      user.RoleId,
		RoleName:    user.RoleName,
		Version:     user.Version,
	}, nil
}
func (s *DossierService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.userService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var userList = []*pb.ListUserReply_Data{}
	for _, user := range list {
		userList = append(userList, &pb.ListUserReply_Data{
			Id:          user.Id,
			Username:    user.Username,
			Mobile:      user.Mobile,
			Description: user.Description,
			Status:      user.Status,
			Avatar:      user.Avatar,
			RoleId:      user.RoleId,
			RoleName:    user.RoleName,
			UpdatedAt:   user.UpdatedAt,
			Version:     user.Version,
		})
	}
	return &pb.ListUserReply{
		Data:  userList,
		Total: int32(count),
	}, nil

}

func (s *DossierService) ChangPassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordReply, error) {

	s.userService.ChangePassword(ctx, req.Id, req.OldPassword, req.NewPassword)

	return &pb.ChangePasswordReply{}, nil
}

// teacher
func (s *DossierService) CreateTeacher(ctx context.Context, req *pb.CreateTeacherRequest) (*pb.CreateTeacherReply, error) {

	_, err := s.teacherService.Create(ctx, &domain.Teacher{
		Name:   req.Name,
		Mobile: req.Mobile,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateTeacherReply{}, nil
}
func (s *DossierService) UpdateTeacher(ctx context.Context, req *pb.UpdateTeacherRequest) (*pb.UpdateTeacherReply, error) {
	count, err := s.teacherService.Update(ctx, &domain.Teacher{
		Id:          int32(req.Id),
		Name:        req.Name,
		Mobile:      req.Mobile,
		Description: req.Description,
		Status:      req.Status,
		Version:     req.Version,
	})
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("更新并发稍后重试")

	}
	return &pb.UpdateTeacherReply{}, nil
}
func (s *DossierService) DeleteTeacher(ctx context.Context, req *pb.DeleteTeacherRequest) (*pb.DeleteTeacherReply, error) {
	return &pb.DeleteTeacherReply{}, nil
}
func (s *DossierService) GetTeacher(ctx context.Context, req *pb.GetTeacherRequest) (*pb.GetTeacherReply, error) {

	user, err := s.teacherService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetTeacherReply{
		Id:      int32(user.Id),
		Name:    user.Name,
		Mobile:  user.Mobile,
		Status:  int32(user.Status),
		Version: user.Version,
	}, nil
}
func (s *DossierService) ListTeacher(ctx context.Context, req *pb.ListTeacherRequest) (*pb.ListTeacherReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.teacherService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var teacherList = []*pb.ListTeacherReply_Data{}
	for _, teacher := range list {
		teacherList = append(teacherList, &pb.ListTeacherReply_Data{
			Id:          teacher.Id,
			Name:        teacher.Name,
			Mobile:      teacher.Mobile,
			Status:      int32(teacher.Status),
			Description: teacher.Description,
			Version:     teacher.Version,
			UpdatedAt:   teacher.UpdatedAt,
		})
	}
	return &pb.ListTeacherReply{
		Data:  teacherList,
		Total: int32(count),
	}, nil
}

// func struct2Query(obj interface{}) map[string]string {

// 	val := reflect.Indirect(reflect.ValueOf(obj))
// 	if val.Kind() != reflect.Struct {
// 		return map[string]string{}
// 	}
// 	var ret = map[string]string{}

// 	for i := 0; i < val.NumField(); i++ {
// 		if val.Type().Field(i).IsExported() {
// 			ret[val.Type().Field(i).Name] = fmt.Sprintf("%v", val.Field(i).Interface())
// 		}
// 	}

// 	return ret
// }

// room
func (s *DossierService) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomReply, error) {

	_, err := s.roomService.Create(ctx, &domain.Room{
		Code:        req.Code,
		Subjects:    req.Subjects,
		Status:      req.Status,
		Description: req.Description,
	})

	if err != nil {
		return nil, err
	}
	return &pb.CreateRoomReply{}, nil
}
func (s *DossierService) UpdateRoom(ctx context.Context, req *pb.UpdateRoomRequest) (*pb.UpdateRoomReply, error) {

	count, err := s.roomService.Update(ctx, &domain.Room{
		Id:          int32(req.Id),
		Code:        req.Code,
		Subjects:    req.Subjects,
		Description: req.Description,
		Status:      req.Status,
		Version:     req.Version,
	})
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("更新并发稍后重试")

	}
	return &pb.UpdateRoomReply{}, nil
}
func (s *DossierService) DeleteRoom(ctx context.Context, req *pb.DeleteRoomRequest) (*pb.DeleteRoomReply, error) {
	return &pb.DeleteRoomReply{}, nil
}
func (s *DossierService) GetRoom(ctx context.Context, req *pb.GetRoomRequest) (*pb.GetRoomReply, error) {

	room, err := s.roomService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetRoomReply{
		Id:          room.Id,
		Code:        room.Code,
		Subjects:    room.Subjects,
		Description: room.Description,
		Status:      int32(room.Status),
		Version:     room.Version,
	}, nil
}
func (s *DossierService) ListRoom(ctx context.Context, req *pb.ListRoomRequest) (*pb.ListRoomReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.roomService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var roomList = []*pb.ListRoomReply_Data{}
	for _, room := range list {
		roomList = append(roomList, &pb.ListRoomReply_Data{
			Id:          room.Id,
			Code:        room.Code,
			Subjects:    room.Subjects,
			Description: room.Description,
			Status:      int32(room.Status),
			UpdatedAt:   room.UpdatedAt,
			Version:     room.Version,
		})
	}
	return &pb.ListRoomReply{
		Data:  roomList,
		Total: int32(count),
	}, nil
}

// subject
func (s *DossierService) CreateSubject(ctx context.Context, req *pb.CreateSubjectRequest) (*pb.CreateSubjectReply, error) {

	_, err := s.subjectService.Create(ctx, &domain.Subject{
		Name:        req.Name,
		Category:    req.Category,
		Status:      req.Status,
		Description: req.Description,
	})

	if err != nil {
		return nil, err
	}
	return &pb.CreateSubjectReply{}, nil
}
func (s *DossierService) UpdateSubject(ctx context.Context, req *pb.UpdateSubjectRequest) (*pb.UpdateSubjectReply, error) {

	count, err := s.subjectService.Update(ctx, &domain.Subject{
		Id:          int32(req.Id),
		Name:        req.Name,
		Category:    req.Category,
		Description: req.Description,
		Status:      req.Status,
		Version:     req.Version,
	})
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("更新并发稍后重试")

	}
	return &pb.UpdateSubjectReply{}, nil
}
func (s *DossierService) DeleteSubject(ctx context.Context, req *pb.DeleteSubjectRequest) (*pb.DeleteSubjectReply, error) {
	return &pb.DeleteSubjectReply{}, nil
}
func (s *DossierService) GetSubject(ctx context.Context, req *pb.GetSubjectRequest) (*pb.GetSubjectReply, error) {

	room, err := s.subjectService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetSubjectReply{
		Id:          room.Id,
		Name:        room.Name,
		Category:    room.Category,
		Description: room.Description,
		Status:      int32(room.Status),
		Version:     room.Version,
	}, nil
}
func (s *DossierService) ListSubject(ctx context.Context, req *pb.ListSubjectRequest) (*pb.ListSubjectReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.subjectService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var subjectList = []*pb.ListSubjectReply_Data{}
	for _, subject := range list {
		subjectList = append(subjectList, &pb.ListSubjectReply_Data{
			Id:          subject.Id,
			Name:        subject.Name,
			Category:    subject.Category,
			Description: subject.Description,
			Status:      int32(subject.Status),
			UpdatedAt:   subject.UpdatedAt,
			Version:     subject.Version,
		})
	}
	return &pb.ListSubjectReply{
		Data:  subjectList,
		Total: int32(count),
	}, nil
}
