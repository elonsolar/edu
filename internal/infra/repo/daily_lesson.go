package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/domain/enum"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type DailyLessonEntity struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	CourseCode   string
	PlanId       int32
	PlanName     string
	PlanDetailId int32

	DateOfDay time.Time

	StartTime         string
	EndTime           string
	RoomId            int32
	RoomName          string
	TeacherId         int32
	TeacherName       string
	ActaulTeacherId   int32 // 代课老师
	ActaulTeacherName string

	SubjectId   int32
	SubjectName string
	GradeId     int32
	GradeName   string
	LessonNum   int32 // 消耗课时
	ActualNum   int32 //报名人数

	SignNum int32 // 签到人数
	Status  int32

	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32 `gorm:"default:1"`
}

func (t *DailyLessonEntity) TableName() string {
	return "daily_lesson"

}

type dailyLessonRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// BatchSave implements domain.DailyLessonRepo.
func (*dailyLessonRepo) BatchSave(context.Context, []*domain.DailyLesson) error {
	panic("unimplemented")
}

// Delete implements domain.DailyLessonRepo.
func (*dailyLessonRepo) Delete(context.Context, int32) (int32, error) {
	panic("unimplemented")
}

// ListAll implements domain.DailyLessonRepo.
func (*dailyLessonRepo) ListAll(context.Context, *domain.Expression) ([]*domain.DailyLesson, error) {
	panic("unimplemented")
}

// ListByMap implements domain.DailyLessonRepo.
func (tp *dailyLessonRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.DailyLesson, error) {

	var dailyLessonEntityList = make([]*DailyLessonEntity, 0)
	err := tp.repo.GetDBFromContext(ctx).Find(&dailyLessonEntityList, params).Error
	if err != nil {
		return nil, err
	}

	result := Map(dailyLessonEntityList, func(lesson *DailyLessonEntity) *domain.DailyLesson {
		return &domain.DailyLesson{
			Id:                int32(lesson.ID),
			PlanId:            lesson.PlanId,
			PlanName:          lesson.PlanName,
			CourseCode:        lesson.CourseCode,
			PlanDetailId:      lesson.PlanDetailId,
			DateOfDay:         lesson.DateOfDay,
			StartTime:         lesson.StartTime,
			EndTime:           lesson.EndTime,
			RoomId:            lesson.RoomId,
			RoomName:          lesson.RoomName,
			TeacherId:         lesson.TeacherId,
			TeacherName:       lesson.TeacherName,
			ActaulTeacherId:   lesson.ActaulTeacherId,
			ActaulTeacherName: lesson.ActaulTeacherName,
			SubjectId:         lesson.SubjectId,
			SubjectName:       lesson.SubjectName,
			GradeId:           lesson.GradeId,
			GradeName:         lesson.GradeName,
			LessonNum:         lesson.LessonNum,
			ActualNum:         lesson.ActualNum,
			SignNum:           lesson.SignNum,
			UpdatedAt:         lesson.UpdatedAt.Format("2006-01-02 15:04:05"),
			Status:            lesson.Status,
			Version:           lesson.Version,
		}
	})
	return result, nil
}

// FindByID implements domain.DailyLessonRepo.
func (up *dailyLessonRepo) FindByID(ctx context.Context, id int32) (*domain.DailyLesson, error) {

	var lesson DailyLessonEntity

	err := up.repo.GetDBFromContext(ctx).First(&lesson, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.DailyLesson{
		Id:                int32(lesson.ID),
		PlanId:            lesson.PlanId,
		PlanName:          lesson.PlanName,
		PlanDetailId:      lesson.PlanDetailId,
		DateOfDay:         lesson.DateOfDay,
		StartTime:         lesson.StartTime,
		EndTime:           lesson.EndTime,
		RoomId:            lesson.RoomId,
		RoomName:          lesson.RoomName,
		TeacherId:         lesson.TeacherId,
		TeacherName:       lesson.TeacherName,
		ActaulTeacherId:   lesson.ActaulTeacherId,
		ActaulTeacherName: lesson.ActaulTeacherName,
		SubjectId:         lesson.SubjectId,
		SubjectName:       lesson.SubjectName,
		GradeId:           lesson.GradeId,
		GradeName:         lesson.GradeName,
		LessonNum:         lesson.LessonNum,
		ActualNum:         lesson.ActualNum,
		SignNum:           lesson.SignNum,
		UpdatedAt:         lesson.UpdatedAt.Format("2006-01-02 15:04:05"),
		Status:            lesson.Status,
		Version:           lesson.Version,
	}, nil
}

// List implements domain.DailyLessonRepo.
func (up *dailyLessonRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.DailyLesson, error) {

	var dailyLessonList []*DailyLessonEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&dailyLessonList).Error
	if err != nil {
		return nil, err
	}

	return Map(dailyLessonList, func(lesson *DailyLessonEntity) *domain.DailyLesson {
		return &domain.DailyLesson{
			Id:                int32(lesson.ID),
			CourseCode:        lesson.CourseCode,
			PlanId:            lesson.PlanId,
			PlanName:          lesson.PlanName,
			PlanDetailId:      lesson.PlanDetailId,
			DateOfDay:         lesson.DateOfDay,
			StartTime:         lesson.StartTime,
			EndTime:           lesson.EndTime,
			RoomId:            lesson.RoomId,
			RoomName:          lesson.RoomName,
			TeacherId:         lesson.TeacherId,
			TeacherName:       lesson.TeacherName,
			ActaulTeacherId:   lesson.ActaulTeacherId,
			ActaulTeacherName: lesson.ActaulTeacherName,
			SubjectId:         lesson.SubjectId,
			SubjectName:       lesson.SubjectName,
			GradeId:           lesson.GradeId,
			GradeName:         lesson.GradeName,
			LessonNum:         lesson.LessonNum,
			ActualNum:         lesson.ActualNum,
			SignNum:           lesson.SignNum,
			UpdatedAt:         lesson.UpdatedAt.Format("2006-01-02 15:04:05"),
			Status:            lesson.Status,
			Version:           lesson.Version,
		}
	}), nil

}

// Count implements domain.DailyLessonRepo.
func (up *dailyLessonRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&DailyLessonEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.DailyLessonRepo.
func (up *dailyLessonRepo) Save(ctx context.Context, lesson *domain.DailyLesson) (*domain.DailyLesson, error) {
	var dailyLessonEntity = &DailyLessonEntity{
		CourseCode:        lesson.CourseCode,
		PlanId:            lesson.PlanId,
		PlanName:          lesson.PlanName,
		PlanDetailId:      lesson.PlanDetailId,
		DateOfDay:         lesson.DateOfDay,
		StartTime:         lesson.StartTime,
		EndTime:           lesson.EndTime,
		RoomId:            lesson.RoomId,
		RoomName:          lesson.RoomName,
		TeacherId:         lesson.TeacherId,
		TeacherName:       lesson.TeacherName,
		ActaulTeacherId:   lesson.ActaulTeacherId,
		ActaulTeacherName: lesson.ActaulTeacherName,
		SubjectId:         lesson.SubjectId,
		SubjectName:       lesson.SubjectName,
		GradeId:           lesson.GradeId,
		GradeName:         lesson.GradeName,
		LessonNum:         lesson.LessonNum,
		ActualNum:         lesson.ActualNum,
		SignNum:           lesson.SignNum,
		Status:            int32(enum.DailyLessonStatusType_PENDING),
	}
	err := up.repo.GetDBFromContext(ctx).Save(dailyLessonEntity).Error
	if err != nil {
		return nil, err
	}
	lesson.Id = int32(dailyLessonEntity.ID)
	return lesson, nil
}

func (up *dailyLessonRepo) Cancel(ctx context.Context, planDetailID int32) (int32, error) {

	result := up.repo.GetDBFromContext(ctx).
		Where(&DailyLessonEntity{Status: int32(enum.DailyLessonStatusType_PENDING), PlanDetailId: planDetailID}).
		Updates(&DailyLessonEntity{Status: int32(enum.DailyLessonStatusType_CANCELED)})
	if result.Error != nil {
		return 0, nil
	}
	return int32(result.RowsAffected), nil
}

// Update implements domain.DailyLessonRepo.
func (up *dailyLessonRepo) Update(ctx context.Context, lesson *domain.DailyLesson) (int32, error) {

	var dailyLessonEntity = &DailyLessonEntity{
		PlanId:            lesson.PlanId,
		PlanName:          lesson.PlanName,
		PlanDetailId:      lesson.PlanDetailId,
		DateOfDay:         lesson.DateOfDay,
		StartTime:         lesson.StartTime,
		EndTime:           lesson.EndTime,
		RoomId:            lesson.RoomId,
		RoomName:          lesson.RoomName,
		TeacherId:         lesson.TeacherId,
		TeacherName:       lesson.TeacherName,
		ActaulTeacherId:   lesson.ActaulTeacherId,
		ActaulTeacherName: lesson.ActaulTeacherName,
		SubjectId:         lesson.SubjectId,
		SubjectName:       lesson.SubjectName,
		GradeId:           lesson.GradeId,
		GradeName:         lesson.GradeName,
		LessonNum:         lesson.LessonNum,
		ActualNum:         lesson.ActualNum,
		SignNum:           lesson.SignNum,
		Status:            lesson.Status,
		Version:           lesson.Version,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&DailyLessonEntity{ID: uint(lesson.Id)}).Updates(dailyLessonEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewDailyLessonRepo .
func NewDailyLessonRepo(repo *BaseRepo, logger log.Logger) domain.DailyLessonRepo {
	return &dailyLessonRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
