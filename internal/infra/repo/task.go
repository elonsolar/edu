package repo

import (
	"context"
	"edu/internal/domain"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type TaskEntity struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string
	Status      int32
	Description string
	ParentId    int32
	ParentName  string
	Code        string // 层级查找用

	Sequence int32
	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32
}

func (t *TaskEntity) TableName() string {
	return "task"

}

type TaskRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// DeleteSubTask implements domain.TaskRepo.
func (tp *TaskRepo) DeleteSubTask(ctx context.Context, code string) (int, error) {
	result := tp.repo.GetDBFromContext(ctx).Where("code like ?", "%"+code+"%").Delete(&TaskEntity{})
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

// ListByParentCode implements domain.TaskRepo.
func (tp *TaskRepo) ListByParentCode(ctx context.Context, code string) ([]*domain.Task, error) {
	var tasks = make([]*TaskEntity, 10)
	err := tp.repo.GetDBFromContext(ctx).Where("code like ?", "%"+code+"%").Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return Map(tasks, func(task *TaskEntity) *domain.Task {
		return &domain.Task{
			Id:          int32(task.ID),
			Name:        task.Name,
			Description: task.Description,
			Code:        task.Code,
			Sequence:    task.Sequence,
			ParentId:    task.ParentId,
			ParentName:  task.ParentName,
			Status:      task.Status,
			Version:     task.Version,
			UpdatedAt:   task.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}), nil
}

// ListByParentID implements domain.TaskRepo.
func (tp *TaskRepo) ListByParentID(ctx context.Context, parentId int32) ([]*domain.Task, error) {

	panic("unimplemented")
}

// FindByID implements domain.TaskRepo.
func (up *TaskRepo) FindByID(ctx context.Context, id int32) (*domain.Task, error) {

	var task TaskEntity

	err := up.repo.GetDBFromContext(ctx).First(&task, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	return &domain.Task{
		Id:          int32(task.ID),
		Name:        task.Name,
		Description: task.Description,
		Code:        task.Code,
		Sequence:    task.Sequence,
		ParentId:    task.ParentId,
		Status:      task.Status,
		Version:     task.Version,
	}, nil
}

// List implements domain.TaskRepo.
func (up *TaskRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.Task, error) {

	var taskList []*TaskEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where("parent_id = 1").Where(query2Clause(query)).Order(" created_at desc").Find(&taskList).Error
	if err != nil {
		return nil, err
	}

	return Map(taskList, func(task *TaskEntity) *domain.Task {
		return &domain.Task{
			Id:          int32(task.ID),
			Name:        task.Name,
			Description: task.Description,
			Code:        task.Code,
			Sequence:    task.Sequence,
			ParentId:    task.ParentId,
			Status:      task.Status,
			Version:     task.Version,
			UpdatedAt:   task.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}), nil

}

// Count implements domain.TaskRepo.
func (up *TaskRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&TaskEntity{}).Where("parent_id =1 ").Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.TaskRepo.
func (up *TaskRepo) Save(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	var taskEntity = &TaskEntity{
		ID:          uint(task.Id),
		Name:        task.Name,
		Description: task.Description,
		Code:        task.Code,
		Sequence:    task.Sequence,
		ParentId:    task.ParentId,
	}
	err := up.repo.GetDBFromContext(ctx).Save(taskEntity).Error
	if err != nil {
		return nil, err
	}
	task.Id = int32(taskEntity.ID)
	return task, nil

}

// Update implements domain.TaskRepo.
func (up *TaskRepo) Update(ctx context.Context, task *domain.Task) (int32, error) {

	var taskEntity = &TaskEntity{
		ID:          uint(task.Id),
		Name:        task.Name,
		Description: task.Description,
		Code:        task.Code,
		Sequence:    task.Sequence,
		ParentId:    task.ParentId,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&TaskEntity{ID: uint(task.Id)}).Updates(taskEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewTaskRepo .
func NewTaskRepo(repo *BaseRepo, logger log.Logger) domain.TaskRepo {
	return &TaskRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
