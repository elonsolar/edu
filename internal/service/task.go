package service

import (
	"context"
	pb "edu/api/assistant/v1"
	"edu/internal/domain"
	"encoding/json"
	"fmt"
)

type TaskService struct {
	pb.UnimplementedTaskServer

	uc *domain.TaskService
	tx domain.Tx
}

func NewTaskService(uc *domain.TaskService, tx domain.Tx) *TaskService {
	return &TaskService{uc: uc, tx: tx}
}

func (s *TaskService) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskReply, error) {

	err := s.tx.Transaction(ctx, func(txctx context.Context) error {

		_, err := s.uc.Create(txctx, &domain.Task{
			Name:        req.Name,
			Description: req.Description,
			ParentId:    req.ParentId,
			ParentName:  req.ParentName,
			Status:      req.Status,
		})
		return err
	})

	if err != nil {
		return nil, err
	}
	return &pb.CreateTaskReply{}, nil
}
func (s *TaskService) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskReply, error) {
	count, err := s.uc.Update(ctx, &domain.Task{
		Id:          int32(req.Id),
		Name:        req.Name,
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
	return &pb.UpdateTaskReply{}, nil
}
func (s *TaskService) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskReply, error) {
	err := s.tx.Transaction(ctx, func(txctx context.Context) error {
		return s.uc.Delete(ctx, req.Ids)
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTaskReply{}, nil
}
func (s *TaskService) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskReply, error) {

	user, err := s.uc.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetTaskReply{
		Id:          int32(user.Id),
		Name:        user.Name,
		Description: user.Description,
		Status:      int32(user.Status),
		Version:     user.Version,
	}, nil
}
func (s *TaskService) ListTask(ctx context.Context, req *pb.ListTaskRequest) (*pb.ListTaskReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.uc.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	return &pb.ListTaskReply{
		Data:  convertTask(list),
		Total: int32(count),
	}, nil
}

func convertTask(dTaskList []*domain.Task) []*pb.ListTaskReply_Data {

	var pTaskList = []*pb.ListTaskReply_Data{}

	for _, dTask := range dTaskList {

		pTask := &pb.ListTaskReply_Data{
			Id:          dTask.Id,
			Name:        dTask.Name,
			Description: dTask.Description,
			ParentId:    dTask.ParentId,
			ParentName:  dTask.ParentName,
			Status:      dTask.Status,
			Version:     dTask.Version,
		}
		if len(dTask.Children) != 0 {
			pTask.Children = convertTask(dTask.Children)
		}
		pTaskList = append(pTaskList, pTask)
	}
	return pTaskList
}
