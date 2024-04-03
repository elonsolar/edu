package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/domain/enum"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type TeacherEntity struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string
	Mobile      string
	Description string
	Status      int32

	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32 `gorm:"default:1"`
}

func (t *TeacherEntity) TableName() string {
	return "teacher"

}

type teacherRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// FindByID implements domain.TeacherRepo.
func (up *teacherRepo) FindByID(ctx context.Context, id int32) (*domain.Teacher, error) {

	var teacher TeacherEntity

	err := up.repo.GetDBFromContext(ctx).First(&teacher, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.Teacher{
		Id:      int32(teacher.ID),
		Name:    teacher.Name,
		Mobile:  teacher.Mobile,
		Status:  teacher.Status,
		Version: teacher.Version,
	}, nil
}

// List implements domain.TeacherRepo.
func (up *teacherRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.Teacher, error) {

	var teacherList []*TeacherEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&teacherList).Error
	if err != nil {
		return nil, err
	}

	return Map(teacherList, func(teacher *TeacherEntity) *domain.Teacher {
		return &domain.Teacher{
			Id:          int32(teacher.ID),
			Name:        teacher.Name,
			Mobile:      teacher.Mobile,
			Status:      teacher.Status,
			Description: teacher.Description,
			Version:     teacher.Version,
			UpdatedAt:   teacher.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}), nil

}

// Count implements domain.TeacherRepo.
func (up *teacherRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&TeacherEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.TeacherRepo.
func (up *teacherRepo) Save(ctx context.Context, teacher *domain.Teacher) (*domain.Teacher, error) {
	var teacherEntity = &TeacherEntity{
		Name:        teacher.Name,
		Mobile:      teacher.Mobile,
		Description: teacher.Description,
		Status:      int32(enum.EnableStatusEnabled),
	}
	err := up.repo.GetDBFromContext(ctx).Create(teacherEntity).Error
	if err != nil {
		return nil, err
	}
	teacher.Id = int32(teacherEntity.ID)
	return teacher, nil

}

// Update implements domain.TeacherRepo.
func (up *teacherRepo) Update(ctx context.Context, teacher *domain.Teacher) (int32, error) {

	var teacherEntity = &TeacherEntity{
		Name:        teacher.Name,
		Mobile:      teacher.Mobile,
		Description: teacher.Description,
		Version:     teacher.Version,
		Status:      teacher.Status,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&TeacherEntity{ID: uint(teacher.Id)}).Updates(teacherEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewTeacherRepo .
func NewTeacherRepo(repo *BaseRepo, logger log.Logger) domain.TeacherRepo {
	return &teacherRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
