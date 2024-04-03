package domain

import (
	"context"
	"edu/internal/domain/enum"
	"fmt"
	"time"
)

type DailyLessonRepo interface {
	BaseRepo[DailyLesson]

	Cancel(ctx context.Context, planDetailID int32) (int32, error)
}

type DailyLesson struct {
	Id           int32
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

	UpdatedAt string
	Status    int32
	Version   int32
}

type DailyLessonService struct {
	BaseService[DailyLesson]
	lessonRepo DailyLessonRepo
}

func (ds *DailyLessonService) Start(ctx context.Context, id int32) error {

	lesson, err := ds.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if lesson.Status != int32(enum.DailyLessonStatusType_PENDING) {
		return fmt.Errorf("课程不是新纪录，无法开始 id: %d", id)
	}

	err = ds.UpdateConcurrency(ctx, &DailyLesson{
		Id:      lesson.Id,
		Status:  int32(enum.DailyLessonStatusType_START),
		Version: lesson.Version,
	})
	if err != nil {
		return fmt.Errorf("更新状态失败 %w", err)
	}
	return nil
}

func (ds *DailyLessonService) Finish(ctx context.Context, id int32) error {

	lesson, err := ds.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if lesson.Status != int32(enum.DailyLessonStatusType_START) {
		return fmt.Errorf("课程不是开始状态，无法结束 id: %d", id)
	}

	num, err := ds.Update(ctx, &DailyLesson{
		Id:      lesson.Id,
		Status:  int32(enum.DailyLessonStatusType_COMPLETED),
		Version: lesson.Version,
	})
	if err != nil {
		return fmt.Errorf("更新状态失败 %w", err)
	}
	if num == 0 {
		return ErrCurrencyUpdate
	}
	return nil
}

func (ds *DailyLessonService) Cancel(ctx context.Context, id int32) error {

	lesson, err := ds.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if lesson.Status != int32(enum.DailyLessonStatusType_PENDING) {
		return fmt.Errorf("非未开始上课的课，不可取消")

	}

	num, err := ds.Update(ctx, &DailyLesson{Id: id, Status: int32(enum.DailyLessonStatusType_CANCELED), Version: lesson.Version})
	if err != nil {
		return fmt.Errorf("更新状态失败 %w", err)
	}
	if num == 0 {
		return ErrCurrencyUpdate
	}
	return nil
}

func NewDailyLessonService(repo DailyLessonRepo) *DailyLessonService {

	return &DailyLessonService{
		BaseService: BaseService[DailyLesson]{
			repo: repo,
		},
		lessonRepo: repo,
	}
}
