package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "edu/api/customer/v1"
	"edu/internal/domain"
	"edu/internal/domain/enum"
)

type CustomerService struct {
	pb.UnimplementedCustomerServer

	tx                   domain.Tx
	customerService      *domain.CustomerService
	studentService       *domain.StudentService
	lessonHistoryService *domain.LessonHistoryService
}

func NewCustomerService(tx domain.Tx, customerService *domain.CustomerService, studentService *domain.StudentService, lessonHistoryService *domain.LessonHistoryService) *CustomerService {
	return &CustomerService{
		tx:                   tx,
		studentService:       studentService,
		customerService:      customerService,
		lessonHistoryService: lessonHistoryService,
	}
}

func (s *CustomerService) CreateStudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.CreateStudentReply, error) {

	birthDay, err := time.ParseInLocation("2006-01-02 15:04:05", req.Birthday, time.Local)

	// birthDay, err := time.Parse("2006-01-02 15:04:05", req.Birthday)
	if err != nil {
		return nil, err
	}

	_, err = s.studentService.Create(ctx, &domain.Student{
		Name:        req.Name,
		Mobile:      req.Mobile,
		Birthday:    birthDay,
		Description: req.Description,
	})

	if err != nil {
		return nil, err
	}
	return &pb.CreateStudentReply{}, nil
}
func (s *CustomerService) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.UpdateStudentReply, error) {
	// birthDay, err := time.Parse("2006-01-02 15:04:05", req.Birthday)

	birthDay, err := time.ParseInLocation("2006-01-02 15:04:05", req.Birthday, time.Local)
	if err != nil {
		return nil, err
	}
	count, err := s.studentService.Update(ctx, &domain.Student{
		Id:          int32(req.Id),
		Name:        req.Name,
		Mobile:      req.Mobile,
		Description: req.Description,
		Birthday:    birthDay,
		Version:     req.Version,
	})
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("更新并发稍后重试")

	}
	return &pb.UpdateStudentReply{}, nil
}
func (s *CustomerService) DeleteStudent(ctx context.Context, req *pb.DeleteStudentRequest) (*pb.DeleteStudentReply, error) {
	return &pb.DeleteStudentReply{}, nil
}
func (s *CustomerService) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.GetStudentReply, error) {

	student, err := s.studentService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetStudentReply{
		Id:          int32(student.Id),
		Name:        student.Name,
		Mobile:      student.Mobile,
		Birthday:    student.Birthday.Format("2006-01-02 15:04:05"),
		Description: student.Description,
		Version:     student.Version,
	}, nil
}
func (s *CustomerService) ListStudent(ctx context.Context, req *pb.ListStudentRequest) (*pb.ListStudentReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.studentService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var studentList = []*pb.ListStudentReply_Data{}
	for _, student := range list {
		studentList = append(studentList, &pb.ListStudentReply_Data{
			Id:          student.Id,
			Name:        student.Name,
			Mobile:      student.Mobile,
			Birthday:    student.Birthday.Format("2006-01-02 15:04:05"),
			Description: student.Description,
			Version:     student.Version,
			UpdatedAt:   student.UpdatedAt,
		})
	}
	return &pb.ListStudentReply{
		Data:  studentList,
		Total: int32(count),
	}, nil
}

// customer
func (s *CustomerService) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerReply, error) {

	_, err := s.customerService.Create(ctx, &domain.Customer{
		Name:        req.Name,
		Mobile:      req.Mobile,
		Description: req.Description,
		Community:   req.Community,
	})

	if err != nil {
		return nil, err
	}
	return &pb.CreateCustomerReply{}, nil
}
func (s *CustomerService) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.UpdateCustomerReply, error) {

	count, err := s.customerService.Update(ctx, &domain.Customer{
		Id:          int32(req.Id),
		Name:        req.Name,
		Mobile:      req.Mobile,
		Description: req.Description,
		Community:   req.Community,
		Version:     req.Version,
	})
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("更新并发稍后重试")

	}
	return &pb.UpdateCustomerReply{}, nil
}
func (s *CustomerService) DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.DeleteCustomerReply, error) {
	return &pb.DeleteCustomerReply{}, nil
}
func (s *CustomerService) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerReply, error) {

	customer, err := s.customerService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetCustomerReply{
		Id:          customer.Id,
		Name:        customer.Name,
		Mobile:      customer.Mobile,
		Description: customer.Description,
		Community:   customer.Community,
		Version:     customer.Version,
	}, nil
}
func (s *CustomerService) ListCustomer(ctx context.Context, req *pb.ListCustomerRequest) (*pb.ListCustomerReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.customerService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var customerList = []*pb.ListCustomerReply_Data{}
	for _, customer := range list {
		customerList = append(customerList, &pb.ListCustomerReply_Data{
			Id:           customer.Id,
			Name:         customer.Name,
			Mobile:       customer.Mobile,
			LessonNumber: customer.LessonNumber,
			Description:  customer.Description,
			Community:    customer.Community,
			UpdatedAt:    customer.UpdatedAt,
			Version:      customer.Version,
		})
	}
	return &pb.ListCustomerReply{
		Data:  customerList,
		Total: int32(count),
	}, nil
}

func (s *CustomerService) AdjustLessonNumber(ctx context.Context, req *pb.AdjustLessonNumberRequest) (*pb.AdjustLessonNumberReply, error) {

	err := s.tx.Transaction(ctx, func(txctx context.Context) error {

		// 更新客户课程
		err := s.customerService.ChangeLessonNum(txctx, req.Id, req.NumChange, req.Version)
		if err != nil {
			return err
		}

		// 增加历史记录
		customer, _ := s.customerService.FindByID(txctx, req.Id)

		_, err = s.lessonHistoryService.Create(txctx, &domain.LessonHistory{
			Mobile:      customer.Mobile,
			SourceType:  int32(enum.ManualAdjust),
			OriginNum:   customer.LessonNumber - req.NumChange,
			NumChange:   req.NumChange,
			Description: req.ChangeDescription,
		})
		return err
	})

	if err != nil {
		return &pb.AdjustLessonNumberReply{}, nil
	}
	return nil, err
}

func (s *CustomerService) ListLessonHistory(ctx context.Context, req *pb.ListLessonHistoryRequest) (*pb.ListLessonHistoryReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.lessonHistoryService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var customerList = []*pb.ListLessonHistoryReply_Data{}
	for _, customer := range list {
		customerList = append(customerList, &pb.ListLessonHistoryReply_Data{
			Id:          customer.Id,
			Mobile:      customer.Mobile,
			OriginNum:   customer.OriginNum,
			NumChange:   customer.NumChange,
			SourceType:  customer.SourceType,
			Description: customer.Description,
			UpdatedAt:   customer.UpdatedAt,
		})
	}
	return &pb.ListLessonHistoryReply{
		Data:  customerList,
		Total: int32(count),
	}, nil
}
