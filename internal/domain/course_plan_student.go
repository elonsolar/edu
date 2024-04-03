package domain

import (
	"context"
	"edu/internal/domain/enum"
)

type CoursePlanStudentRepo interface {
	BaseRepo[CoursePlanStudent]
}

type CoursePlanStudent struct {
	Id            int32
	CustomerPhone string
	StudentId     int32
	StudentName   string

	PlanId   int32
	PlanName string

	PlanDetailId int32
	UpdatedAt    string
	Status       int32
	Version      int32
}

type CoursePlanStudentService struct {
	BaseService[CoursePlanStudent]
}

func (cps *CoursePlanStudentService) Schedule(ctx context.Context, student *CoursePlanStudent) error {
	err := cps.UpdateConcurrency(ctx, &CoursePlanStudent{Id: student.Id, Status: int32(enum.CoursePlanStudentStatusType_SCHEDULED), Version: student.Version})
	if err != nil {
		return err
	}
	return nil
}

func NewCoursePlanStudentService(repo CoursePlanStudentRepo) *CoursePlanStudentService {

	return &CoursePlanStudentService{
		BaseService: BaseService[CoursePlanStudent]{repo: repo},
	}
}
