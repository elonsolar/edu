package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/domain/enum"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type CoursePlanStudentEntity struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CustomerPhone string
	StudentId     int32
	StudentName   string

	PlanId   int32
	PlanName string

	PlanDetailId int32

	Status   int32
	UserId   int32 // 和登录账号绑定
	TenantId int32
	Version  int32 `gorm:"default:1"`
}

func (t *CoursePlanStudentEntity) TableName() string {
	return "course_plan_student"

}

type coursePlanStudentRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// ListByMap implements domain.CoursePlanStudentRepo.
func (tp *coursePlanStudentRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.CoursePlanStudent, error) {
	var coursePlanStudentList = make([]*CoursePlanStudentEntity, 0)
	err := tp.repo.GetDBFromContext(ctx).Model(&CoursePlanStudentEntity{}).Find(&coursePlanStudentList, params).Error
	if err != nil {
		return nil, err
	}

	result := Map(coursePlanStudentList, func(student *CoursePlanStudentEntity) *domain.CoursePlanStudent {
		return &domain.CoursePlanStudent{
			Id:            int32(student.ID),
			CustomerPhone: student.CustomerPhone,
			StudentId:     student.StudentId,
			StudentName:   student.StudentName,
			PlanId:        student.PlanId,
			PlanName:      student.PlanName,
			PlanDetailId:  student.PlanDetailId,
			UpdatedAt:     student.UpdatedAt.Format("2006-01-02 15:04:05"),
			Status:        student.Status,
			Version:       student.Version,
		}
	})
	return result, nil
}

// BatchSave implements domain.CoursePlanStudentRepo.
func (*coursePlanStudentRepo) BatchSave(context.Context, []*domain.CoursePlanStudent) error {
	panic("unimplemented")
}

// Delete implements domain.CoursePlanStudentRepo.
func (*coursePlanStudentRepo) Delete(context.Context, int32) (int32, error) {
	panic("unimplemented")
}

// ListAll implements domain.CoursePlanStudentRepo.
func (*coursePlanStudentRepo) ListAll(context.Context, *domain.Expression) ([]*domain.CoursePlanStudent, error) {
	panic("unimplemented")
}

// FindByID implements domain.CoursePlanStudentRepo.
func (up *coursePlanStudentRepo) FindByID(ctx context.Context, id int32) (*domain.CoursePlanStudent, error) {

	var coursePlanStudent CoursePlanStudentEntity

	err := up.repo.GetDBFromContext(ctx).First(&coursePlanStudent, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.CoursePlanStudent{
		Id:            int32(coursePlanStudent.ID),
		CustomerPhone: coursePlanStudent.CustomerPhone,
		StudentId:     coursePlanStudent.StudentId,
		StudentName:   coursePlanStudent.StudentName,
		PlanId:        coursePlanStudent.PlanId,
		PlanName:      coursePlanStudent.PlanName,
		PlanDetailId:  coursePlanStudent.PlanDetailId,
		UpdatedAt:     coursePlanStudent.UpdatedAt.Format("2006-01-02 15:04:05"),
		Status:        coursePlanStudent.Status,
		Version:       coursePlanStudent.Version,
	}, nil
}

// List implements domain.CoursePlanStudentRepo.
func (up *coursePlanStudentRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.CoursePlanStudent, error) {

	var coursePlanStudentList []*CoursePlanStudentEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&coursePlanStudentList).Error
	if err != nil {
		return nil, err
	}

	return Map(coursePlanStudentList, func(coursePlanStudent *CoursePlanStudentEntity) *domain.CoursePlanStudent {
		return &domain.CoursePlanStudent{
			Id:            int32(coursePlanStudent.ID),
			CustomerPhone: coursePlanStudent.CustomerPhone,
			StudentId:     coursePlanStudent.StudentId,
			StudentName:   coursePlanStudent.StudentName,
			PlanId:        coursePlanStudent.PlanId,
			PlanName:      coursePlanStudent.PlanName,
			PlanDetailId:  coursePlanStudent.PlanDetailId,
			UpdatedAt:     coursePlanStudent.UpdatedAt.Format("2006-01-02 15:04:05"),
			Status:        coursePlanStudent.Status,
			Version:       coursePlanStudent.Version,
		}
	}), nil

}

// Count implements domain.CoursePlanStudentRepo.
func (up *coursePlanStudentRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&CoursePlanStudentEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.CoursePlanStudentRepo.
func (up *coursePlanStudentRepo) Save(ctx context.Context, coursePlanStudent *domain.CoursePlanStudent) (*domain.CoursePlanStudent, error) {
	var coursePlanStudentEntity = &CoursePlanStudentEntity{
		CustomerPhone: coursePlanStudent.CustomerPhone,
		StudentId:     coursePlanStudent.StudentId,
		StudentName:   coursePlanStudent.StudentName,
		PlanId:        coursePlanStudent.PlanId,
		PlanName:      coursePlanStudent.PlanName,
		PlanDetailId:  coursePlanStudent.PlanDetailId,
		Status:        int32(enum.CoursePlanStudentStatusType_NEW),
	}
	err := up.repo.GetDBFromContext(ctx).Save(coursePlanStudentEntity).Error
	if err != nil {
		return nil, err
	}
	coursePlanStudent.Id = int32(coursePlanStudentEntity.ID)
	return coursePlanStudent, nil

}

// Update implements domain.CoursePlanStudentRepo.
func (up *coursePlanStudentRepo) Update(ctx context.Context, coursePlanStudent *domain.CoursePlanStudent) (int32, error) {
	var coursePlanStudentEntity = &CoursePlanStudentEntity{
		CustomerPhone: coursePlanStudent.CustomerPhone,
		StudentId:     coursePlanStudent.StudentId,
		StudentName:   coursePlanStudent.StudentName,
		PlanId:        coursePlanStudent.PlanId,
		PlanName:      coursePlanStudent.PlanName,
		PlanDetailId:  coursePlanStudent.PlanDetailId,
		Status:        coursePlanStudent.Status,
		Version:       coursePlanStudent.Version,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&CoursePlanStudentEntity{ID: uint(coursePlanStudent.Id)}).Updates(coursePlanStudentEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewCoursePlanStudentRepo .
func NewCoursePlanStudentRepo(repo *BaseRepo, logger log.Logger) domain.CoursePlanStudentRepo {
	return &coursePlanStudentRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
