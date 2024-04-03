package domain

import (
	"context"
	"errors"
	"fmt"
)

type TaskRepo interface {
	Save(context.Context, *Task) (*Task, error)
	Update(context.Context, *Task) (int32, error)
	FindByID(context.Context, int32) (*Task, error)
	List(context.Context, *Expression, *Page) ([]*Task, error)
	Count(context.Context, *Expression) (int32, error)
	ListByParentID(context.Context, int32) ([]*Task, error)
	ListByParentCode(context.Context, string) ([]*Task, error)
	DeleteSubTask(context.Context, string) (int, error)
}

type Task struct {
	Id          int32
	Name        string
	Status      int32
	Description string
	ParentId    int32
	ParentName  string
	Code        string // 层级查找用
	Children    []*Task

	Sequence  int32
	UpdatedAt string
	Version   int32
}

type TaskService struct {
	repo TaskRepo
}

func NewTaskService(repo TaskRepo) *TaskService {

	return &TaskService{repo: repo}
}

func (t *TaskService) Create(ctx context.Context, task *Task) (*Task, error) {

	if task.ParentId == 0 {
		task.ParentId = 1
	}

	parent, err := t.repo.FindByID(ctx, task.ParentId)
	if err != nil && !errors.Is(err, ErrRecordNotFound) {
		return nil, err
	}

	//
	if errors.Is(err, ErrRecordNotFound) {
		parent, err = t.repo.Save(ctx, &Task{Code: "", Sequence: 0, Id: 1})
		if err != nil {
			return nil, err
		}
	}

	parentCode := ""
	var sequence int32 = 0

	if parent != nil {

		parent.Sequence += 1

		parentCode = parent.Code
		sequence = parent.Sequence

		count, err := t.repo.Update(ctx, parent)
		if err != nil {
			return nil, err
		}
		if count != 1 {
			return nil, ErrCurrencyUpdate
		}
	}

	task.Code = fmt.Sprintf("%s%02d", parentCode, sequence+1)

	return t.repo.Save(ctx, task)
}

func (t *TaskService) Update(ctx context.Context, task *Task) (int32, error) {

	return t.repo.Update(ctx, task)
}

func (t *TaskService) FindByID(ctx context.Context, id int32) (*Task, error) {

	return t.repo.FindByID(ctx, id)
}

func (t *TaskService) List(ctx context.Context, query *Expression, page *Page) ([]*Task, int32, error) {

	count, err := t.repo.Count(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	if count == 0 {
		return nil, 0, nil
	}

	list, err := t.repo.List(ctx, query, page)
	if err != nil {
		return nil, 0, err
	}
	for _, task := range list {

		children, err := t.repo.ListByParentCode(ctx, task.Code)
		if err != nil {
			return nil, 0, err
		}

		taskMap := make(map[int32]*Task, len(children))
		taskMap[task.Id] = task

		for _, child := range children {
			if parent, ok := taskMap[child.ParentId]; ok {
				parent.Children = append(parent.Children, child)

			} else {
				parent := &Task{
					Id: child.ParentId,
				}
				parent.Children = append(parent.Children, child)
				taskMap[child.ParentId] = parent
			}

			if self, exist := taskMap[child.Id]; exist {
				self.Name = child.Name
				self.Code = child.Code
				self.Description = child.Description
			} else {
				taskMap[child.Id] = child
			}

		}
	}

	return list, count, nil
}

// Delete delete task and sub tasks
func (t *TaskService) Delete(ctx context.Context, ids []int32) error {

	for _, id := range ids {
		task, err := t.repo.FindByID(ctx, id)
		if err != nil && !errors.Is(err, ErrRecordNotFound) {
			return err
		}

		_, err = t.repo.DeleteSubTask(ctx, task.Code)
		if err != nil {
			return err
		}
	}

	return nil
}

// type Sink[I, O any] struct {
// 	opType int
// 	fn     func(I) O
// }

// func NewSink[I, O any]() *Sink[I, O] {
// 	return &Sink[I, O]{}
// }

// func Map[I, O any](fn func(src I) (dst O)) *Sink[I, O] {

// 	return &Sink[I, O]{fn: fn}
// }

// func Filter[I any, O bool](fn func(src I) O) *Sink[I, O] {

// 	return &Sink[I, O]{fn: fn}
// }

// func collectToList(sinks Sink[I,O any]) any {

// 	for _, val := range s.srcVal {
// 		var nv = val

// 		for i := len(s.ops) - 1; i > 0; i-- {
// 			// if
// 			nv = s.ops[i].fn(nv)
// 		}
// 		clt(nv)
// 	}

// 	return nil
// }
