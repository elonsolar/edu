package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/domain/enum"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type DailyLessonStudentEntity struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	PlanId        int32
	PlanName      string
	PlanDetailId  int32
	LessonId      int32
	CustomerPhone string
	StudentId     int32
	StudentName   string
	Status        int32

	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32 `gorm:"default:1"`
}

func (t *DailyLessonStudentEntity) TableName() string {
	return "daily_lesson_student"

}

type dailyLessonStudentRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.DailyLessonStudentRepo.
func (up *dailyLessonStudentRepo) BatchSave(ctx context.Context, dailyLessonStudentList []*domain.DailyLessonStudent) error {

	if len(dailyLessonStudentList) == 0 {
		up.log.Debug()
		return nil
	}

	detailList := Map(dailyLessonStudentList, func(student *domain.DailyLessonStudent) *DailyLessonStudentEntity {

		return &DailyLessonStudentEntity{
			ID:            uint(student.Id),
			PlanId:        student.PlanId,
			PlanName:      student.PlanName,
			PlanDetailId:  student.PlanDetailId,
			LessonId:      student.LessonId,
			CustomerPhone: student.CustomerPhone,
			StudentId:     student.StudentId,
			StudentName:   student.StudentName,
			Status:        int32(enum.DailyLessonStudentStatusType_UNSIGNED),
			UserId:        0,
			TenantId:      0,
		}
	})

	err := up.repo.GetDBFromContext(ctx).Save(detailList).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete implements domain.DailyLessonStudentRepo.
func (*dailyLessonStudentRepo) Delete(context.Context, int32) (int32, error) {
	panic("unimplemented")
}

// ListAll implements domain.DailyLessonStudentRepo.
func (*dailyLessonStudentRepo) ListAll(context.Context, *domain.Expression) ([]*domain.DailyLessonStudent, error) {
	panic("unimplemented")
}

// ListByMap implements domain.DailyLessonStudentRepo.
func (tp *dailyLessonStudentRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.DailyLessonStudent, error) {

	var dailyLessonEntityList = make([]*DailyLessonStudentEntity, 0)
	err := tp.repo.GetDBFromContext(ctx).Find(&dailyLessonEntityList, params).Error
	if err != nil {
		return nil, err
	}

	result := Map(dailyLessonEntityList, func(lesson *DailyLessonStudentEntity) *domain.DailyLessonStudent {
		return &domain.DailyLessonStudent{
			Id:            int32(lesson.ID),
			PlanId:        lesson.PlanId,
			PlanName:      lesson.PlanName,
			PlanDetailId:  lesson.PlanDetailId,
			LessonId:      lesson.LessonId,
			CustomerPhone: lesson.CustomerPhone,
			StudentId:     lesson.StudentId,
			StudentName:   lesson.StudentName,
			UpdatedAt:     lesson.UpdatedAt.Format("2006-01-02 15:04:05"),
			Status:        lesson.Status,
			Version:       lesson.Version,
		}
	})
	return result, nil
}

// FindByID implements domain.DailyLessonStudentRepo.
func (up *dailyLessonStudentRepo) FindByID(ctx context.Context, id int32) (*domain.DailyLessonStudent, error) {

	var lesson DailyLessonStudentEntity

	err := up.repo.GetDBFromContext(ctx).First(&lesson, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.DailyLessonStudent{
		Id:            int32(lesson.ID),
		PlanId:        lesson.PlanId,
		PlanName:      lesson.PlanName,
		PlanDetailId:  lesson.PlanDetailId,
		LessonId:      lesson.LessonId,
		CustomerPhone: lesson.CustomerPhone,
		StudentId:     lesson.StudentId,
		StudentName:   lesson.StudentName,
		UpdatedAt:     lesson.UpdatedAt.Format("2006-01-02 15:04:05"),
		Status:        lesson.Status,
		Version:       lesson.Version,
	}, nil
}

// List implements domain.DailyLessonStudentRepo.
func (up *dailyLessonStudentRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.DailyLessonStudent, error) {

	var dailyLessonList []*DailyLessonStudentEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&dailyLessonList).Error
	if err != nil {
		return nil, err
	}

	return Map(dailyLessonList, func(lesson *DailyLessonStudentEntity) *domain.DailyLessonStudent {
		return &domain.DailyLessonStudent{
			Id:            int32(lesson.ID),
			PlanId:        lesson.PlanId,
			PlanName:      lesson.PlanName,
			PlanDetailId:  lesson.PlanDetailId,
			LessonId:      lesson.LessonId,
			CustomerPhone: lesson.CustomerPhone,
			StudentId:     lesson.StudentId,
			StudentName:   lesson.StudentName,
			UpdatedAt:     lesson.UpdatedAt.Format("2006-01-02 15:04:05"),
			Status:        lesson.Status,
			Version:       lesson.Version,
		}
	}), nil

}

// Count implements domain.DailyLessonStudentRepo.
func (up *dailyLessonStudentRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&DailyLessonStudentEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.DailyLessonStudentRepo.
func (up *dailyLessonStudentRepo) Save(ctx context.Context, lesson *domain.DailyLessonStudent) (*domain.DailyLessonStudent, error) {
	var dailyLessonEntity = &DailyLessonStudentEntity{
		PlanId:        lesson.PlanId,
		PlanName:      lesson.PlanName,
		PlanDetailId:  lesson.PlanDetailId,
		LessonId:      lesson.LessonId,
		CustomerPhone: lesson.CustomerPhone,
		StudentId:     lesson.StudentId,
		StudentName:   lesson.StudentName,
		Status:        int32(enum.DailyLessonStudentStatusType_UNSIGNED),
		UserId:        0,
		TenantId:      0,
	}
	err := up.repo.GetDBFromContext(ctx).Save(dailyLessonEntity).Error
	if err != nil {
		return nil, err
	}
	lesson.Id = int32(dailyLessonEntity.ID)
	return lesson, nil

}

// Update implements domain.DailyLessonStudentRepo.
func (up *dailyLessonStudentRepo) Update(ctx context.Context, lesson *domain.DailyLessonStudent) (int32, error) {

	var dailyLessonEntity = &DailyLessonStudentEntity{
		ID:            0,
		PlanId:        lesson.PlanId,
		PlanName:      lesson.PlanName,
		PlanDetailId:  lesson.PlanDetailId,
		LessonId:      lesson.LessonId,
		CustomerPhone: lesson.CustomerPhone,
		StudentId:     lesson.StudentId,
		StudentName:   lesson.StudentName,
		Status:        lesson.Status,
		Version:       lesson.Version,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&DailyLessonStudentEntity{ID: uint(lesson.Id)}).Updates(dailyLessonEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewDailyLessonStudentRepo .
func NewDailyLessonStudentRepo(repo *BaseRepo, logger log.Logger) domain.DailyLessonStudentRepo {
	return &dailyLessonStudentRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
