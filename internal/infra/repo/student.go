package repo

import (
	"context"
	"edu/internal/domain"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type StudentEntity struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string
	Mobile      string
	Birthday    time.Time
	Description string

	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32 `gorm:"default:1"`
}

func (t *StudentEntity) TableName() string {
	return "student"

}

type studentRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// FindByID implements domain.StudentRepo.
func (up *studentRepo) FindByID(ctx context.Context, id int32) (*domain.Student, error) {

	var student StudentEntity

	err := up.repo.GetDBFromContext(ctx).First(&student, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.Student{
		Id:          int32(student.ID),
		Name:        student.Name,
		Mobile:      student.Mobile,
		Description: student.Description,
		Birthday:    student.Birthday,
	}, nil
}

// List implements domain.StudentRepo.
func (up *studentRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.Student, error) {

	var studentList []*StudentEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&studentList).Error
	if err != nil {
		return nil, err
	}

	return Map(studentList, func(student *StudentEntity) *domain.Student {
		return &domain.Student{
			Id:          int32(student.ID),
			Name:        student.Name,
			Mobile:      student.Mobile,
			Description: student.Description,
			Birthday:    student.Birthday,
			Version:     student.Version,
			UpdatedAt:   student.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}), nil

}

// Count implements domain.StudentRepo.
func (up *studentRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&StudentEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.StudentRepo.
func (up *studentRepo) Save(ctx context.Context, student *domain.Student) (*domain.Student, error) {
	var studentEntity = &StudentEntity{
		Name:        student.Name,
		Mobile:      student.Mobile,
		Description: student.Description,
		Birthday:    student.Birthday,
	}
	err := up.repo.GetDBFromContext(ctx).Save(studentEntity).Error
	if err != nil {
		return nil, err
	}
	student.Id = int32(studentEntity.ID)
	return student, nil

}

// Update implements domain.StudentRepo.
func (up *studentRepo) Update(ctx context.Context, student *domain.Student) (int32, error) {

	var studentEntity = &StudentEntity{
		Name:        student.Name,
		Mobile:      student.Mobile,
		Description: student.Description,
		Birthday:    student.Birthday,
		Version:     student.Version,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&StudentEntity{ID: uint(student.Id)}).Updates(studentEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewStudentRepo .
func NewStudentRepo(repo *BaseRepo, logger log.Logger) domain.StudentRepo {
	return &studentRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
