package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/domain/enum"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type CoursePlanDetailEntity struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Code         string
	PlanId       int32
	PlanName     string
	DayIndex     int32
	DayIndexName string
	StartTime    string
	EndTime      string
	RoomId       int32
	RoomName     string
	TeacherId    int32
	TeacherName  string
	SubjectId    int32
	SubjectName  string
	GradeId      int32
	GradeName    string
	LessonNum    int32 // 消耗课时
	PlanNum      int32 // 计划上课人数
	ActualNum    int32 // 报名人数
	Status       int32

	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32 `gorm:"default:1"`
}

func (t *CoursePlanDetailEntity) TableName() string {
	return "course_plan_detail"

}

type coursePlanDetailRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// ListByMap implements domain.CoursePlanDetailRepo.
func (tp *coursePlanDetailRepo) ListByMap(ctx context.Context, params map[string]interface{}) ([]*domain.CoursePlanDetail, error) {

	var coursePlanDetailList = make([]*CoursePlanDetailEntity, 0)
	err := tp.repo.GetDBFromContext(ctx).Model(&CoursePlanDetailEntity{}).Find(&coursePlanDetailList, params).Error
	if err != nil {
		return nil, err
	}

	result := Map(coursePlanDetailList, func(detail *CoursePlanDetailEntity) *domain.CoursePlanDetail {
		return &domain.CoursePlanDetail{
			Id:           int32(detail.ID),
			Code:         detail.Code,
			PlanId:       detail.PlanId,
			PlanName:     detail.PlanName,
			DayIndex:     detail.DayIndex,
			DayIndexName: detail.DayIndexName,
			StartTime:    detail.StartTime,
			EndTime:      detail.EndTime,
			RoomId:       detail.RoomId,
			RoomName:     detail.RoomName,
			TeacherId:    detail.TeacherId,
			TeacherName:  detail.TeacherName,
			SubjectId:    detail.SubjectId,
			SubjectName:  detail.SubjectName,
			GradeId:      detail.GradeId,
			GradeName:    detail.GradeName,
			UpdatedAt:    detail.UpdatedAt.Format("2006-01-02 15:04:05"),
			LessonNum:    detail.LessonNum,
			PlanNum:      detail.PlanNum,
			ActualNum:    detail.ActualNum,
			Status:       detail.Status,
			Version:      detail.Version,
		}
	})
	return result, nil
}

// Delete implements domain.CoursePlanDetailRepo.
func (tp *coursePlanDetailRepo) Delete(ctx context.Context, id int32) (int32, error) {
	result := tp.repo.GetDBFromContext(ctx).Delete(&CoursePlanDetailEntity{}, id)
	if result.Error != nil {
		return 0, result.Error
	}
	return int32(result.RowsAffected), nil
}

// ListAll implements domain.CoursePlanDetailRepo.
func (up *coursePlanDetailRepo) ListAll(ctx context.Context, query *domain.Expression) ([]*domain.CoursePlanDetail, error) {
	var coursePlanDetailList []*CoursePlanDetailEntity
	err := up.repo.GetDBFromContext(ctx).Where(query2Clause(query)).Find(&coursePlanDetailList).Error
	if err != nil {
		return nil, err
	}
	return Map(coursePlanDetailList, func(detail *CoursePlanDetailEntity) *domain.CoursePlanDetail {
		return &domain.CoursePlanDetail{
			Id:           int32(detail.ID),
			Code:         detail.Code,
			PlanId:       detail.PlanId,
			PlanName:     detail.PlanName,
			DayIndex:     detail.DayIndex,
			DayIndexName: detail.DayIndexName,
			StartTime:    detail.StartTime,
			EndTime:      detail.EndTime,
			RoomId:       detail.RoomId,
			RoomName:     detail.RoomName,
			TeacherId:    detail.TeacherId,
			TeacherName:  detail.TeacherName,
			SubjectId:    detail.SubjectId,
			SubjectName:  detail.SubjectName,
			GradeId:      detail.GradeId,
			GradeName:    detail.GradeName,
			UpdatedAt:    detail.UpdatedAt.Format("2006-01-02 15:04:05"),
			LessonNum:    detail.LessonNum,
			PlanNum:      detail.PlanNum,
			ActualNum:    detail.ActualNum,
			Status:       detail.Status,
			Version:      detail.Version,
		}
	}), nil
}

// FindByID implements domain.CoursePlanDetailRepo.
func (up *coursePlanDetailRepo) FindByID(ctx context.Context, id int32) (*domain.CoursePlanDetail, error) {

	var coursePlanDetail CoursePlanDetailEntity

	err := up.repo.GetDBFromContext(ctx).First(&coursePlanDetail, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.CoursePlanDetail{
		Id:           int32(coursePlanDetail.ID),
		PlanId:       coursePlanDetail.PlanId,
		DayIndex:     coursePlanDetail.DayIndex,
		DayIndexName: coursePlanDetail.DayIndexName,
		StartTime:    coursePlanDetail.StartTime,
		EndTime:      coursePlanDetail.EndTime,
		RoomId:       coursePlanDetail.RoomId,
		RoomName:     coursePlanDetail.RoomName,
		TeacherId:    coursePlanDetail.TeacherId,
		TeacherName:  coursePlanDetail.TeacherName,
		SubjectId:    coursePlanDetail.SubjectId,
		SubjectName:  coursePlanDetail.SubjectName,
		GradeId:      coursePlanDetail.GradeId,
		GradeName:    coursePlanDetail.GradeName,
		LessonNum:    coursePlanDetail.LessonNum,
		PlanName:     coursePlanDetail.PlanName,
		ActualNum:    coursePlanDetail.ActualNum,
		Status:       coursePlanDetail.Status,
		Version:      coursePlanDetail.Version,
	}, nil
}

// List implements domain.CoursePlanDetailRepo.
func (up *coursePlanDetailRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.CoursePlanDetail, error) {

	var coursePlanDetailList []*CoursePlanDetailEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&coursePlanDetailList).Error
	if err != nil {
		return nil, err
	}

	return Map(coursePlanDetailList, func(coursePlanDetail *CoursePlanDetailEntity) *domain.CoursePlanDetail {
		return &domain.CoursePlanDetail{
			Id:           int32(coursePlanDetail.ID),
			Code:         coursePlanDetail.Code,
			PlanId:       coursePlanDetail.PlanId,
			PlanName:     coursePlanDetail.PlanName,
			DayIndex:     coursePlanDetail.DayIndex,
			DayIndexName: coursePlanDetail.DayIndexName,
			StartTime:    coursePlanDetail.StartTime,
			EndTime:      coursePlanDetail.EndTime,
			RoomId:       coursePlanDetail.RoomId,
			RoomName:     coursePlanDetail.RoomName,
			TeacherId:    coursePlanDetail.TeacherId,
			TeacherName:  coursePlanDetail.TeacherName,
			SubjectId:    coursePlanDetail.SubjectId,
			SubjectName:  coursePlanDetail.SubjectName,
			GradeId:      coursePlanDetail.GradeId,
			GradeName:    coursePlanDetail.GradeName,
			LessonNum:    coursePlanDetail.LessonNum,
			PlanNum:      coursePlanDetail.PlanNum,
			ActualNum:    coursePlanDetail.ActualNum,
			Status:       coursePlanDetail.Status,
			Version:      coursePlanDetail.Version,
			UpdatedAt:    coursePlanDetail.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}), nil

}

// Count implements domain.CoursePlanDetailRepo.
func (up *coursePlanDetailRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&CoursePlanDetailEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.CoursePlanDetailRepo.
func (up *coursePlanDetailRepo) Save(ctx context.Context, coursePlanDetail *domain.CoursePlanDetail) (*domain.CoursePlanDetail, error) {
	var coursePlanDetailEntity = &CoursePlanDetailEntity{
		Code:         coursePlanDetail.Code,
		PlanId:       coursePlanDetail.PlanId,
		PlanName:     coursePlanDetail.PlanName,
		DayIndex:     coursePlanDetail.DayIndex,
		DayIndexName: coursePlanDetail.DayIndexName,
		StartTime:    coursePlanDetail.StartTime,
		EndTime:      coursePlanDetail.EndTime,
		RoomId:       coursePlanDetail.RoomId,
		RoomName:     coursePlanDetail.RoomName,
		TeacherId:    coursePlanDetail.TeacherId,
		TeacherName:  coursePlanDetail.TeacherName,
		SubjectId:    coursePlanDetail.SubjectId,
		SubjectName:  coursePlanDetail.SubjectName,
		GradeId:      coursePlanDetail.GradeId,
		GradeName:    coursePlanDetail.GradeName,
		LessonNum:    coursePlanDetail.LessonNum,
		PlanNum:      coursePlanDetail.PlanNum,
		Status:       int32(enum.CoursePlanDetailStatusType_NEW),
	}
	err := up.repo.GetDBFromContext(ctx).Save(coursePlanDetailEntity).Error
	if err != nil {
		return nil, err
	}
	coursePlanDetail.Id = int32(coursePlanDetailEntity.ID)
	return coursePlanDetail, nil

}

// BatchSave implements domain.CoursePlanDetailRepo.
func (up *coursePlanDetailRepo) BatchSave(ctx context.Context, coursePlanDetailList []*domain.CoursePlanDetail) error {

	detailList := Map(coursePlanDetailList, func(coursePlanDetail *domain.CoursePlanDetail) *CoursePlanDetailEntity {

		return &CoursePlanDetailEntity{
			ID:           uint(coursePlanDetail.Id),
			Code:         coursePlanDetail.Code,
			PlanId:       coursePlanDetail.PlanId,
			PlanName:     coursePlanDetail.PlanName,
			DayIndex:     coursePlanDetail.DayIndex,
			DayIndexName: coursePlanDetail.DayIndexName,
			StartTime:    coursePlanDetail.StartTime,
			EndTime:      coursePlanDetail.EndTime,
			RoomId:       coursePlanDetail.RoomId,
			RoomName:     coursePlanDetail.RoomName,
			TeacherId:    coursePlanDetail.TeacherId,
			TeacherName:  coursePlanDetail.TeacherName,
			SubjectId:    coursePlanDetail.SubjectId,
			SubjectName:  coursePlanDetail.SubjectName,
			GradeId:      coursePlanDetail.GradeId,
			GradeName:    coursePlanDetail.GradeName,
			LessonNum:    coursePlanDetail.LessonNum,
			PlanNum:      coursePlanDetail.PlanNum,
			Status:       int32(enum.CoursePlanDetailStatusType_NEW),
		}
	})

	err := up.repo.GetDBFromContext(ctx).Save(detailList).Error
	if err != nil {
		return err
	}
	return nil
}

// Update implements domain.CoursePlanDetailRepo.
func (up *coursePlanDetailRepo) Update(ctx context.Context, coursePlanDetail *domain.CoursePlanDetail) (int32, error) {

	var coursePlanDetailEntity = &CoursePlanDetailEntity{
		PlanId:       coursePlanDetail.PlanId,
		Code:         coursePlanDetail.Code,
		PlanName:     coursePlanDetail.PlanName,
		DayIndex:     coursePlanDetail.DayIndex,
		DayIndexName: coursePlanDetail.DayIndexName,
		StartTime:    coursePlanDetail.StartTime,
		EndTime:      coursePlanDetail.EndTime,
		RoomId:       coursePlanDetail.RoomId,
		RoomName:     coursePlanDetail.RoomName,
		TeacherId:    coursePlanDetail.TeacherId,
		TeacherName:  coursePlanDetail.TeacherName,
		SubjectId:    coursePlanDetail.SubjectId,
		SubjectName:  coursePlanDetail.SubjectName,
		GradeId:      coursePlanDetail.GradeId,
		GradeName:    coursePlanDetail.GradeName,
		LessonNum:    coursePlanDetail.LessonNum,
		PlanNum:      coursePlanDetail.PlanNum,
		ActualNum:    coursePlanDetail.ActualNum,
		Status:       coursePlanDetail.Status,
		Version:      coursePlanDetail.Version,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&CoursePlanDetailEntity{ID: uint(coursePlanDetail.Id)}).Updates(coursePlanDetailEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewCoursePlanDetailRepo .
func NewCoursePlanDetailRepo(repo *BaseRepo, logger log.Logger) domain.CoursePlanDetailRepo {
	return &coursePlanDetailRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
