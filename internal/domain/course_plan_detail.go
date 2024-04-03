package domain

import (
	"context"
	"edu/internal/domain/enum"
	"fmt"
)

type CoursePlanDetailRepo interface {
	BaseRepo[CoursePlanDetail]
}

type CoursePlanDetail struct {
	Id           int32
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
	UpdatedAt    string
	// Status    int32
	Version int32
}

type CoursePlanDetailService struct {
	BaseService[CoursePlanDetail]
}

func (cs *CoursePlanDetailService) Schedule(ctx context.Context, planDetailId int32) error {

	planDetail, err := cs.FindByID(ctx, planDetailId)
	if err != nil {
		return err
	}

	if planDetail.Status != int32(enum.CoursePlanDetailStatusType_NEW) {
		return nil
	}

	err = cs.UpdateConcurrency(ctx, &CoursePlanDetail{
		Id:      planDetailId,
		Status:  int32(enum.CoursePlanDetailStatusType_SCHEDULED),
		Version: planDetail.Version,
	})
	if err != nil {
		return fmt.Errorf("更新课程状态错误 %w", err)
	}

	return nil
}

func NewCoursePlanDetailService(repo CoursePlanDetailRepo) *CoursePlanDetailService {

	return &CoursePlanDetailService{
		BaseService: BaseService[CoursePlanDetail]{
			repo: repo,
		},
	}
}
