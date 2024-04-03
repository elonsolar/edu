package domain

import (
	"context"
	"edu/internal/domain/enum"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type CoursePlanRepo interface {
	Save(context.Context, *CoursePlan) (*CoursePlan, error)
	Update(context.Context, *CoursePlan) (int32, error)
	FindByID(context.Context, int32) (*CoursePlan, error)
	List(context.Context, *Expression, *Page) ([]*CoursePlan, error)
	Count(context.Context, *Expression) (int32, error)
}

type CoursePlan struct {
	Id          int32
	Name        string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	CycleType   int32

	ExcludeRule []struct {
		ExcludeType int32
		ExcludeDate []string // [周六], [2023-01-02,2023-03-02]
	}

	UpdatedAt string
	Status    int32
	Version   int32
}

type CoursePlanService struct {
	repo CoursePlanRepo
	log  *log.Helper
}

func NewCoursePlanService(logger log.Logger, repo CoursePlanRepo) *CoursePlanService {

	return &CoursePlanService{
		log:  log.NewHelper(logger),
		repo: repo,
	}
}

func (t *CoursePlanService) Create(ctx context.Context, teacher *CoursePlan) (*CoursePlan, error) {

	return t.repo.Save(ctx, teacher)
}

func (t *CoursePlanService) Update(ctx context.Context, teacher *CoursePlan) (int32, error) {

	return t.repo.Update(ctx, teacher)
}

func (t *CoursePlanService) FindByID(ctx context.Context, id int32) (*CoursePlan, error) {

	return t.repo.FindByID(ctx, id)
}

func (t *CoursePlanService) List(ctx context.Context, query *Expression, page *Page) ([]*CoursePlan, int32, error) {

	count, err := t.repo.Count(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	if count == 0 {
		return nil, 0, nil
	}

	list, err := t.repo.List(ctx, query, page)
	if err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (t *CoursePlanService) Release(ctx context.Context, idList []int32) (int32, int32, string) {

	var releaseOne = func(c context.Context, id int32) (int32, error) {
		plan, err := t.FindByID(c, id)
		if err != nil {
			return 0, err
		}
		if plan.Status != int32(enum.CoursePlanStatusType_NEW) {
			return 0, errors.New("计划状态不是新纪录，无法发布!")
		}

		updateParams := &CoursePlan{Id: id, Version: plan.Version, Status: int32(enum.CoursePlanStatusType_RELEASED)}
		count, err := t.repo.Update(ctx, updateParams)
		if err != nil {
			return 0, err
		}

		return count, nil
	}

	var (
		successNum int32
		failureNum int32
		buf        = strings.Builder{}
	)
	for _, id := range idList {

		num, err := releaseOne(ctx, id)
		if err != nil {
			buf.WriteString(err.Error())
			failureNum++
			continue
		}

		if num == 0 {
			buf.WriteString(fmt.Sprintf("更新并发,id:%d", id))
			failureNum++
		} else {
			successNum++
		}
	}

	return successNum, failureNum, buf.String()
}

// ScheduleDate
func (t *CoursePlanService) ScheduleDate(ctx context.Context, plan *CoursePlan, detailList []*CoursePlanDetail) ([]*DailyLesson, error) {

	var dailyLessonList = make([]*DailyLesson, 0, 20)

	if plan.CycleType == 0 { //日
		for today, endDate := plan.StartTime, plan.EndTime; today.Compare(endDate) <= 0; today = today.Add(time.Hour * 24) {

			if t.isDateShouldSkip(ctx, plan, today) {
				continue
			}

			// 所有课程，都复制到当天
			for _, detail := range detailList {
				dailyLessonList = append(dailyLessonList, &DailyLesson{
					PlanId:            plan.Id,
					PlanName:          plan.Name,
					CourseCode:        detail.Code,
					PlanDetailId:      detail.Id,
					DateOfDay:         today,
					StartTime:         detail.StartTime,
					EndTime:           detail.EndTime,
					RoomId:            detail.RoomId,
					RoomName:          detail.RoomName,
					TeacherId:         detail.TeacherId,
					TeacherName:       detail.TeacherName,
					ActaulTeacherId:   0,
					ActaulTeacherName: "",
					SubjectId:         detail.SubjectId,
					SubjectName:       detail.SubjectName,
					GradeId:           detail.GradeId,
					GradeName:         detail.GradeName,
					LessonNum:         detail.LessonNum,
					ActualNum:         detail.LessonNum,
					SignNum:           0,
					Status:            int32(enum.DailyLessonStatusType_PENDING),
				})

			}
		}
	} else if plan.CycleType == 1 { // 周
		for today, endDate := plan.StartTime, plan.EndTime; today.Compare(endDate) <= 0; today = today.Add(time.Hour * 24) {
			if t.isDateShouldSkip(ctx, plan, today) {
				continue
			}

			// 周一课，复制到周一

			for _, detail := range detailList {

				var weekDay = today.Weekday() // 1,2,3,4,5,6,0 -> 1,2,3,4,5,6,7 (周末weekday 是0)
				if detail.DayIndex == int32(weekDay) || (detail.DayIndex == 7 && weekDay == 0) {

					dailyLessonList = append(dailyLessonList, &DailyLesson{
						CourseCode:        detail.Code,
						PlanId:            plan.Id,
						PlanName:          plan.Name,
						PlanDetailId:      detail.Id,
						DateOfDay:         today,
						StartTime:         detail.StartTime,
						EndTime:           detail.EndTime,
						RoomId:            detail.RoomId,
						RoomName:          detail.RoomName,
						TeacherId:         detail.TeacherId,
						TeacherName:       detail.TeacherName,
						ActaulTeacherId:   0,
						ActaulTeacherName: "",
						SubjectId:         detail.SubjectId,
						SubjectName:       detail.SubjectName,
						GradeId:           detail.GradeId,
						GradeName:         detail.GradeName,
						LessonNum:         detail.LessonNum,
						ActualNum:         detail.LessonNum,
						SignNum:           0,
						Status:            int32(enum.DailyLessonStatusType_PENDING),
					})

				}

			}

		}

	}

	return dailyLessonList, nil
}

func (t *CoursePlanService) isDateShouldSkip(ctx context.Context, plan *CoursePlan, target time.Time) bool {

	for _, rule := range plan.ExcludeRule {

		if rule.ExcludeType == int32(enum.CoursePlanExcludeDateType_WeekDay) {

			var targetWeekDay = target.Weekday()
			for _, weekDayOfRule := range rule.ExcludeDate {
				if weekDayOfRule == "周一" && targetWeekDay == 1 ||
					weekDayOfRule == "周二" && targetWeekDay == 2 ||
					weekDayOfRule == "周三" && targetWeekDay == 3 ||
					weekDayOfRule == "周四" && targetWeekDay == 4 ||
					weekDayOfRule == "周五" && targetWeekDay == 5 ||
					weekDayOfRule == "周六" && targetWeekDay == 6 ||
					weekDayOfRule == "周日" && targetWeekDay == 0 {

					t.log.WithContext(ctx).Infof("日期 %s ,匹配排除规则限定 %s", target.Format("2006-01-02"), weekDayOfRule)
					return true
				}
			}
		} else {
			begin, end := rule.ExcludeDate[0], rule.ExcludeDate[1]
			targetDateStr := target.Format("2006-01-02")
			if strings.Compare(targetDateStr, begin) >= 0 && strings.Compare(end, targetDateStr) >= 0 {
				t.log.WithContext(ctx).Infof("日期 %s ,匹配排除规则限定 %s ~ %s", target.Format("2006-01-02"), begin, end)
				return true
			}
		}
	}
	return false
}

func (t *CoursePlanService) ScheduleDateDaily() {

}

func (t *CoursePlanService) ScheduleDateWeekly() {}

// 自动排课
// 找出当前上课的各种科目人数，按照时间自动生成课程
// todo 需要统计各种科目和人数

// 手动排课
// 手动录入时间和科目，生成课程
// 排课周期，决定排课的天数
// 周期为1天，只需要排1天
// 周期为1周，排课需要排1周
// 考虑实际上课时间可能不从周一开始， 可以设置生成课程时候，排除的日期，也满足了节假日不上课的要求，但是调休的这种，可以后续用调课功能
// 考虑 周期小于1周（1天） , 周六周日不上课，不需要排课，可以设置排出的兴起
