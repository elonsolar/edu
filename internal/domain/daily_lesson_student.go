package domain

import (
	"context"
	"edu/internal/domain/enum"
)

type DailyLessonStudentRepo interface {
	BaseRepo[DailyLessonStudent]
}

type DailyLessonStudent struct {
	Id int32

	PlanId        int32
	PlanName      string
	PlanDetailId  int32
	LessonId      int32
	CustomerPhone string
	StudentId     int32
	StudentName   string

	UpdatedAt string
	Status    int32
	Version   int32
}

type DailyLessonStudentService struct {
	BaseService[DailyLessonStudent]
}

func (s *DailyLessonStudentService) CreateByPlanStudent(ctx context.Context, lesson *DailyLesson, planStudentList []*CoursePlanStudent) error {

	lessonStudentList := make([]*DailyLessonStudent, 0, len(planStudentList))
	for _, planStudent := range planStudentList {
		lessonStudentList = append(lessonStudentList, &DailyLessonStudent{
			PlanId:        planStudent.PlanId,
			PlanName:      planStudent.PlanName,
			PlanDetailId:  planStudent.PlanDetailId,
			LessonId:      lesson.Id,
			CustomerPhone: planStudent.CustomerPhone,
			StudentId:     planStudent.StudentId,
			StudentName:   planStudent.StudentName,
			Status:        int32(enum.DailyLessonStudentStatusType_UNSIGNED),
		})
	}

	err := s.BatchCreate(ctx, lessonStudentList)
	if err != nil {
		return err
	}
	return nil
}

func NewDailyLessonStudentService(repo DailyLessonStudentRepo) *DailyLessonStudentService {

	return &DailyLessonStudentService{
		BaseService: BaseService[DailyLessonStudent]{
			repo: repo,
		},
	}
}
