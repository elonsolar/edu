package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"

	pb "edu/api/course/v1"
	"edu/internal/domain"
	"edu/internal/domain/enum"

	"github.com/go-kratos/kratos/v2/log"
)

type CourseService struct {
	pb.UnimplementedCourseServer

	tx                        domain.Tx
	coursePlanService         *domain.CoursePlanService
	coursePlanDetailService   *domain.CoursePlanDetailService
	coursePlanStudentService  *domain.CoursePlanStudentService
	dailyLessonService        *domain.DailyLessonService
	dailyLessonStudentService *domain.DailyLessonStudentService
	metaService               *domain.MetaService
	log                       *log.Helper
}

func NewCourseService(logger log.Logger, tx domain.Tx, coursePlanService *domain.CoursePlanService, coursePlanDetailService *domain.CoursePlanDetailService, coursePlanStudentService *domain.CoursePlanStudentService, dailyLessonService *domain.DailyLessonService, dailyLessonStudentService *domain.DailyLessonStudentService, metaService *domain.MetaService) *CourseService {
	return &CourseService{
		log:                       log.NewHelper(logger),
		tx:                        tx,
		coursePlanService:         coursePlanService,
		coursePlanDetailService:   coursePlanDetailService,
		coursePlanStudentService:  coursePlanStudentService,
		dailyLessonService:        dailyLessonService,
		dailyLessonStudentService: dailyLessonStudentService,
		metaService:               metaService,
	}
}

func (s *CourseService) CreateCoursePlan(ctx context.Context, req *pb.CreateCoursePlanRequest) (*pb.CreateCoursePlanReply, error) {

	var coursePlan = &domain.CoursePlan{
		Name: req.Name,
		// StartTime:   req.StartTime,
		// EndTime:     req.EndTime,
		CycleType: req.CycleType,
		// CourseTime:  req.CourseTime,
		// ExcludeDay:  req.ExcludeDay,
		// ExcludeDate: req.ExcludeDate,
		Description: req.Description,
	}

	startTime, err := time.ParseInLocation("2006-01-02", req.StartTime, time.Local)
	if err != nil {
		return nil, fmt.Errorf("startTime parse err:%w", err)
	}

	endTime, err := time.ParseInLocation("2006-01-02", req.EndTime, time.Local)
	if err != nil {
		return nil, fmt.Errorf("endTime parse err:%w", err)
	}

	coursePlan.StartTime = startTime
	coursePlan.EndTime = endTime

	for _, rule := range req.ExcludeRule {
		coursePlan.ExcludeRule = append(coursePlan.ExcludeRule, struct {
			ExcludeType int32
			ExcludeDate []string
		}{
			ExcludeType: rule.ExcludeType,
			ExcludeDate: rule.ExcludeDate,
		})
	}

	_, err = s.coursePlanService.Create(ctx, coursePlan)

	if err != nil {
		return nil, err
	}
	return &pb.CreateCoursePlanReply{}, nil
}
func (s *CourseService) UpdateCoursePlan(ctx context.Context, req *pb.UpdateCoursePlanRequest) (*pb.UpdateCoursePlanReply, error) {
	var coursePlan = &domain.CoursePlan{
		Id:   req.Id,
		Name: req.Name,
		// StartTime:   req.StartTime,
		// EndTime:     req.EndTime,
		CycleType: req.CycleType,
		// CourseTime:  req.CourseTime,
		// ExcludeDay:  req.ExcludeDay,
		// ExcludeDate: req.ExcludeDate,
		Description: req.Description,
		Version:     req.Version,
	}

	startTime, err := time.ParseInLocation("2006-01-02", req.StartTime, time.Local)
	if err != nil {
		return nil, fmt.Errorf("startTime parse err:%w", err)
	}

	endTime, err := time.ParseInLocation("2006-01-02", req.EndTime, time.Local)
	if err != nil {
		return nil, fmt.Errorf("endTime parse err:%w", err)
	}

	coursePlan.StartTime = startTime
	coursePlan.EndTime = endTime

	for _, rule := range req.ExcludeRule {
		coursePlan.ExcludeRule = append(coursePlan.ExcludeRule, struct {
			ExcludeType int32
			ExcludeDate []string
		}{
			ExcludeType: rule.ExcludeType,
			ExcludeDate: rule.ExcludeDate,
		})
	}
	count, err := s.coursePlanService.Update(ctx, coursePlan)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("更新并发稍后重试")

	}
	return &pb.UpdateCoursePlanReply{}, nil
}

func (s *CourseService) DeleteCoursePlan(ctx context.Context, req *pb.DeleteCoursePlanRequest) (*pb.DeleteCoursePlanReply, error) {
	return &pb.DeleteCoursePlanReply{}, nil
}

func (s *CourseService) GetCoursePlan(ctx context.Context, req *pb.GetCoursePlanRequest) (*pb.GetCoursePlanReply, error) {

	coursePlan, err := s.coursePlanService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	reply := &pb.GetCoursePlanReply{
		Id:        int32(coursePlan.Id),
		Name:      coursePlan.Name,
		StartTime: coursePlan.StartTime.Format("2006-01-02"),
		EndTime:   coursePlan.EndTime.Format("2006-01-02"),
		CycleType: coursePlan.CycleType,
		// CourseTime:  coursePlan.CourseTime.Format("2006-01-02 15:04:05"),
		// ExcludeDay:  coursePlan.ExcludeDay,
		// ExcludeDate: coursePlan.ExcludeDate,
		Description: coursePlan.Description,
		Version:     coursePlan.Version,
	}

	return reply, nil
}

func (s *CourseService) ListCoursePlan(ctx context.Context, req *pb.ListCoursePlanRequest) (*pb.ListCoursePlanReply, error) {

	ctx, span := tracer.Start(ctx, "CourseService.ListCoursePlan")
	defer span.End()
	span.SetAttributes(attribute.String("query", req.Expr))
	// s.log.WithContext(ctx).Info("hello")

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.coursePlanService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var coursePlanList = []*pb.ListCoursePlanReply_Data{}
	for _, coursePlan := range list {
		reply := &pb.ListCoursePlanReply_Data{
			Id:        int32(coursePlan.Id),
			Name:      coursePlan.Name,
			StartTime: coursePlan.StartTime.Format("2006-01-02"),
			EndTime:   coursePlan.EndTime.Format("2006-01-02"),
			CycleType: coursePlan.CycleType,
			// CourseTime:  coursePlan.CourseTime.Format("2006-01-02 15:04:05"),
			// ExcludeDay:  coursePlan.ExcludeDay,
			// ExcludeDate: coursePlan.ExcludeDate,
			Status:      coursePlan.Status,
			Description: coursePlan.Description,
			Version:     coursePlan.Version,
			UpdatedAt:   coursePlan.UpdatedAt,
		}

		for _, rule := range coursePlan.ExcludeRule {
			reply.ExcludeRule = append(reply.ExcludeRule, &pb.ListCoursePlanReply_Data_ExcludeRule{
				ExcludeType: rule.ExcludeType,
				ExcludeDate: rule.ExcludeDate,
			})
		}

		coursePlanList = append(coursePlanList, reply)
	}
	return &pb.ListCoursePlanReply{
		Data:  coursePlanList,
		Total: int32(count),
	}, nil
}

func (s *CourseService) ReleaseCoursePlan(ctx context.Context, req *pb.ReleaseCoursePlanRequest) (*pb.ReleaseCoursePlanReply, error) {

	success, failure, msg := s.coursePlanService.Release(ctx, req.IdList)
	return &pb.ReleaseCoursePlanReply{
		SuccessNum: success,
		FailureNum: failure,
		Message:    msg,
	}, nil
}

// 2023-10-19 08:08
// 已发布状态，才可以排期
// 当查询计划的当前状态为已发布，但是在接下来的程序执行过程中，计划的状态可能被其他并发操作改变。判断的依据不再可靠
// 如何解决？ 那让 所有涉及状态的操作互斥，
// 1.可以加进程外，或进程内显示的排他锁，
// 2.  利用数据库的行锁互斥（也算是一种进程外，但不是显示的锁）（比如查询出当前的状态合法，就立刻更新该记录（要用下面的乐观锁，并开启事物）
// 3. 利用cas，用版本号或者时间戳做乐观锁（开启事务，多个操作可以回滚）
//
// 2依赖于3
// 2 利用乐观锁，和数据库的行锁，模拟悲观锁， 特点是获取锁的时间长（对数据库资源占用大），但是可以保证接下来的操作一定可以完成（除了外部异常）
// 可以考虑在操作完成后，再执行cas操作，但是可能有并发的操作，多次执行，但是最终只有一个失败的情况 ，但是依然可以保持数据一致性
func (s *CourseService) ScheduleDateForCoursePlan(ctx context.Context, req *pb.ScheduleDateForCoursePlanRequest) (*pb.ScheduleDateForCoursePlanReply, error) {

	plan, err := s.coursePlanService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 已发布 才可以排期
	if plan.Status != int32(enum.CoursePlanStatusType_RELEASED) {
		return nil, errors.New("已发布的计划才可以排期")
	}

	var detailParams = map[string]interface{}{"plan_id": req.Id, "status": enum.CoursePlanDetailStatusType_NEW}
	if len(req.DetailIds) != 0 {
		detailParams["id"] = req.DetailIds
	}

	// 计划课程信息
	detailList, err := s.coursePlanDetailService.ListByMap(ctx, detailParams)
	if err != nil {
		return nil, fmt.Errorf("查询课程明细错误: %w", err)
	}

	var schedulePlanDetail = func(detail *domain.CoursePlanDetail) error {

		studentList, err := s.coursePlanStudentService.ListByMap(ctx, map[string]interface{}{"plan_detail_id": detail.Id})
		if err != nil {
			return err
		}

		dailyLessonList, err := s.coursePlanService.ScheduleDate(ctx, plan, []*domain.CoursePlanDetail{detail})
		if err != nil {
			return err
		}

		err = s.tx.Transaction(ctx, func(txctx context.Context) error {

			for _, dailyLesson := range dailyLessonList {
				newLesson, err := s.dailyLessonService.Create(txctx, dailyLesson)
				if err != nil {
					return err
				}

				err = s.dailyLessonStudentService.CreateByPlanStudent(txctx, newLesson, studentList)
				if err != nil {
					return err
				}
			}

			// 更新课程排课状态
			err = s.coursePlanDetailService.UpdateConcurrency(ctx, &domain.CoursePlanDetail{
				Id:      detail.Id,
				Status:  int32(enum.CoursePlanDetailStatusType_SCHEDULED),
				Version: detail.Version,
			})
			if err != nil {
				return err
			}

			for _, student := range studentList {
				err := s.coursePlanStudentService.UpdateConcurrency(ctx, &domain.CoursePlanStudent{Id: student.Id, Status: int32(enum.CoursePlanStudentStatusType_SCHEDULED), Version: student.Version})
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	}

	// 对每节课 排期
	for _, detail := range detailList {
		err = schedulePlanDetail(detail)
		if err != nil {
			return nil, err
		}
	}

	return &pb.ScheduleDateForCoursePlanReply{}, nil
}

// ScheduleDateForCoursePlanDetail  schedule single plan Detail
func (s *CourseService) ScheduleDateForCoursePlanDetail(ctx context.Context, req *pb.ScheduleDateForCoursePlanDetailRequest) (*pb.ScheduleDateForCoursePlanDetailReply, error) {

	detail, err := s.coursePlanDetailService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("查找课程错误 %w", err)
	}

	_, err = s.ScheduleDateForCoursePlan(ctx, &pb.ScheduleDateForCoursePlanRequest{
		Id:        detail.PlanId,
		DetailIds: []int32{req.Id},
	})
	if err != nil {
		return nil, err
	}

	return &pb.ScheduleDateForCoursePlanDetailReply{}, nil
}

// coursePlanDetail
func (s *CourseService) CreateCoursePlanDetail(ctx context.Context, req *pb.CreateCoursePlanDetailRequest) (*pb.CreateCoursePlanDetailReply, error) {

	code, err := s.metaService.GetLessonSequenceCode(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.coursePlanDetailService.Create(ctx, &domain.CoursePlanDetail{
		Id:           req.Id,
		Code:         code,
		PlanId:       req.PlanId,
		PlanName:     req.PlanName,
		DayIndex:     req.DayIndex,
		DayIndexName: req.DayIndexName,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		RoomId:       req.RoomId,
		RoomName:     req.RoomName,
		TeacherId:    req.TeacherId,
		TeacherName:  req.TeacherName,
		SubjectId:    req.SubjectId,
		SubjectName:  req.SubjectName,
		GradeId:      req.GradeId,
		GradeName:    req.GradeName,
		LessonNum:    req.LessonNum,
		PlanNum:      req.PlanNum,
		ActualNum:    req.ActualNum,
		Status:       req.Status,
	})

	if err != nil {
		return nil, err
	}

	return &pb.CreateCoursePlanDetailReply{}, nil
}
func (s *CourseService) UpdateCoursePlanDetail(ctx context.Context, req *pb.UpdateCoursePlanDetailRequest) (*pb.UpdateCoursePlanDetailReply, error) {

	count, err := s.coursePlanDetailService.Update(ctx, &domain.CoursePlanDetail{
		Id:          int32(req.Id),
		Code:        req.Code,
		PlanId:      req.PlanId,
		PlanName:    req.PlanName,
		DayIndex:    req.DayIndex,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		RoomId:      req.RoomId,
		RoomName:    req.RoomName,
		TeacherId:   req.TeacherId,
		TeacherName: req.TeacherName,
		SubjectId:   req.SubjectId,
		SubjectName: req.SubjectName,
		GradeId:     req.GradeId,
		GradeName:   req.GradeName,
		LessonNum:   req.LessonNum,
		PlanNum:     req.PlanNum,
		Version:     req.Version,
	})
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("更新并发稍后重试")

	}
	return &pb.UpdateCoursePlanDetailReply{}, nil
}
func (s *CourseService) DeleteCoursePlanDetail(ctx context.Context, req *pb.DeleteCoursePlanDetailRequest) (*pb.DeleteCoursePlanDetailReply, error) {

	planDetail, err := s.coursePlanDetailService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if planDetail.Status != int32(enum.CoursePlanDetailStatusType_NEW) {
		return nil, errors.New("不是新纪录不得删除")
	}

	if planDetail.ActualNum != 0 {
		return nil, errors.New("已经有学生报名，不得删除")
	}

	_, err = s.coursePlanDetailService.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteCoursePlanDetailReply{}, nil
}
func (s *CourseService) GetCoursePlanDetail(ctx context.Context, req *pb.GetCoursePlanDetailRequest) (*pb.GetCoursePlanDetailReply, error) {

	planDetail, err := s.coursePlanDetailService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetCoursePlanDetailReply{
		Id:          planDetail.Id,
		PlanId:      planDetail.PlanId,
		DayIndex:    planDetail.DayIndex,
		StartTime:   planDetail.StartTime,
		EndTime:     planDetail.EndTime,
		RoomId:      planDetail.RoomId,
		RoomName:    planDetail.RoomName,
		TeacherId:   planDetail.TeacherId,
		TeacherName: planDetail.TeacherName,
		SubjectId:   planDetail.SubjectId,
		SubjectName: planDetail.SubjectName,
		GradeId:     planDetail.GradeId,
		GradeName:   planDetail.GradeName,
		Version:     planDetail.Version,
	}, nil
}
func (s *CourseService) ListCoursePlanDetail(ctx context.Context, req *pb.ListCoursePlanDetailRequest) (*pb.ListCoursePlanDetailReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.coursePlanDetailService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var planDetailList = []*pb.ListCoursePlanDetailReply_Data{}
	for _, detail := range list {
		planDetailList = append(planDetailList, &pb.ListCoursePlanDetailReply_Data{
			Id:           detail.Id,
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
			LessonNum:    detail.LessonNum,
			PlanNum:      detail.PlanNum,
			ActualNum:    detail.ActualNum,
			Status:       detail.Status,
			Version:      detail.Version,
		})
	}
	return &pb.ListCoursePlanDetailReply{
		Data:  planDetailList,
		Total: int32(count),
	}, nil
}
func (s *CourseService) ListAllCoursePlanDetail(ctx context.Context, req *pb.ListAllCoursePlanDetailRequest) (*pb.ListAllCoursePlanDetailReply, error) {
	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}
	list, err := s.coursePlanDetailService.ListAll(ctx, &query)
	if err != nil {
		return nil, err
	}

	var detailList = make([]*pb.ListAllCoursePlanDetailReply_Data, 0, len(list))
	for _, detail := range list {

		detailList = append(detailList, &pb.ListAllCoursePlanDetailReply_Data{
			Id:           detail.Id,
			PlanId:       detail.PlanId,
			Code:         detail.Code,
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
			LessonNum:    detail.LessonNum,
			PlanNum:      detail.PlanNum,
			ActualNum:    detail.ActualNum,
			Status:       detail.Status,
			Version:      detail.Version,
		})

	}

	return &pb.ListAllCoursePlanDetailReply{
		Data: detailList,
	}, nil
}

func (s *CourseService) BatchAddCoursePlanDetail(ctx context.Context, req *pb.BatchAddCoursePlanDetailRequest) (*pb.BatchAddCoursePlanDetailReply, error) {

	var details = make([]*domain.CoursePlanDetail, 0, len(req.Data))

	for _, detail := range req.Data {

		var newDetail = &domain.CoursePlanDetail{
			Id:           detail.Id,
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
			LessonNum:    detail.LessonNum,
			PlanNum:      detail.PlanNum,
			ActualNum:    detail.ActualNum,
			Status:       detail.Status,
		}

		if len(newDetail.Code) == 0 {
			code, err := s.metaService.GetLessonSequenceCode(ctx)
			if err != nil {
				return nil, err
			}
			newDetail.Code = code
		}
		details = append(details, newDetail)
	}

	err := s.coursePlanDetailService.BatchCreate(ctx, details)
	if err != nil {
		return nil, err
	}

	return &pb.BatchAddCoursePlanDetailReply{}, nil
}

// 停课,取消所有未上课的课
func (s *CourseService) StopCoursePlanDetail(ctx context.Context, req *pb.StopCoursePlanDetailRequest) (*pb.StopCoursePlanDetailReply, error) {

	planDetail, err := s.coursePlanDetailService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("查询课程错误 err:%w", err)
	}
	if planDetail.Status != int32(enum.CoursePlanDetailStatusType_SCHEDULED) {
		return nil, fmt.Errorf("课程状态不为已排期,无需取消")
	}

	err = s.tx.Transaction(ctx, func(txctx context.Context) error {
		// 取消所有待上课
		err = s.dailyLessonService.Cancel(txctx, req.Id)
		if err != nil {
			return fmt.Errorf("取消课程错误 err:%w", err)
		}

		// 更新 课程为停课
		err = s.coursePlanDetailService.UpdateConcurrency(txctx, &domain.CoursePlanDetail{Id: planDetail.Id, Version: planDetail.Version, Status: int32(enum.CoursePlanDetailStatusType_STOPPED)})
		if err != nil {
			return fmt.Errorf("修改课程状态错误 err:%w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &pb.StopCoursePlanDetailReply{}, nil
}

// CreateCoursePlanStudent ,课程新增学生，如果课程已经排期，需要将学生添加的日程表
func (s *CourseService) CreateCoursePlanStudent(ctx context.Context, req *pb.CreateCoursePlanStudentRequest) (*pb.CreateCoursePlanStudentReply, error) {

	studentList, err := s.coursePlanStudentService.ListByMap(ctx, map[string]interface{}{"student_id": req.StudentId, "plan_detail_id": req.PlanDetailId})
	if err != nil {
		return nil, err
	}
	if len(studentList) > 0 {
		return nil, fmt.Errorf("学生已经存在")
	}

	err = s.tx.Transaction(ctx, func(txctx context.Context) error {

		var coursePlanStudent = &domain.CoursePlanStudent{
			CustomerPhone: req.CustomerPhone,
			StudentId:     req.StudentId,
			StudentName:   req.StudentName,
			PlanId:        req.PlanId,
			PlanName:      req.PlanName,
			PlanDetailId:  req.PlanDetailId,
		}

		coursePlanStudent, err := s.coursePlanStudentService.Create(txctx, coursePlanStudent)

		if err != nil {
			return err
		}

		coursePlanDetail, err := s.coursePlanDetailService.FindByID(txctx, coursePlanStudent.PlanDetailId)
		if err != nil {
			return err
		}

		// 更新班级人数
		err = s.coursePlanDetailService.UpdateConcurrency(txctx, &domain.CoursePlanDetail{Id: coursePlanDetail.Id, Version: coursePlanDetail.Version, ActualNum: coursePlanDetail.ActualNum + 1})
		if err != nil {
			return err
		}

		// 已经排课，需要将学生信息同步到 日程课表
		if coursePlanDetail.Status == int32(enum.CoursePlanDetailStatusType_SCHEDULED) {

			pendingDailyLessonList, err := s.dailyLessonService.ListByMap(txctx, map[string]interface{}{"plan_detail_id": coursePlanDetail.Id, "status": int32(enum.DailyLessonStatusType_PENDING)})
			if err != nil {
				return err
			}

			var dailyLessonStudent = make([]*domain.DailyLessonStudent, 0, len(pendingDailyLessonList))
			for _, dailyLesson := range pendingDailyLessonList {
				dailyLessonStudent = append(dailyLessonStudent, &domain.DailyLessonStudent{
					PlanId:        coursePlanStudent.PlanId,
					PlanName:      coursePlanStudent.PlanName,
					PlanDetailId:  coursePlanDetail.Id,
					LessonId:      dailyLesson.Id,
					CustomerPhone: coursePlanStudent.CustomerPhone,
					StudentId:     coursePlanStudent.StudentId,
					StudentName:   coursePlanStudent.StudentName,
					Status:        int32(enum.DailyLessonStudentStatusType_UNSIGNED),
				})
			}

			err = s.dailyLessonStudentService.BatchCreate(txctx, dailyLessonStudent)
			if err != nil {
				return err
			}

			return s.coursePlanStudentService.UpdateConcurrency(txctx, &domain.CoursePlanStudent{Id: coursePlanStudent.Id, Status: int32(enum.CoursePlanStudentStatusType_SCHEDULED)})

		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &pb.CreateCoursePlanStudentReply{}, nil
}
func (s *CourseService) UpdateCoursePlanStudent(ctx context.Context, req *pb.UpdateCoursePlanStudentRequest) (*pb.UpdateCoursePlanStudentReply, error) {
	var coursePlanStudent = &domain.CoursePlanStudent{
		Id:            req.Id,
		CustomerPhone: req.CustomerPhone,
		StudentId:     req.StudentId,
		StudentName:   req.StudentName,
		PlanId:        req.PlanId,
		PlanName:      req.PlanName,
		PlanDetailId:  req.PlanDetailId,
		Version:       req.Version,
	}

	count, err := s.coursePlanStudentService.Update(ctx, coursePlanStudent)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("更新并发稍后重试")

	}
	return &pb.UpdateCoursePlanStudentReply{}, nil
}

func (s *CourseService) DeleteCoursePlanStudent(ctx context.Context, req *pb.DeleteCoursePlanStudentRequest) (*pb.DeleteCoursePlanStudentReply, error) {
	return &pb.DeleteCoursePlanStudentReply{}, nil
}

func (s *CourseService) GetCoursePlanStudent(ctx context.Context, req *pb.GetCoursePlanStudentRequest) (*pb.GetCoursePlanStudentReply, error) {

	coursePlanStudent, err := s.coursePlanStudentService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	reply := &pb.GetCoursePlanStudentReply{
		Id:            int32(coursePlanStudent.Id),
		CustomerPhone: coursePlanStudent.CustomerPhone,
		StudentId:     coursePlanStudent.StudentId,
		StudentName:   coursePlanStudent.PlanName,
		PlanId:        coursePlanStudent.PlanId,
		PlanName:      coursePlanStudent.PlanName,
		PlanDetailId:  coursePlanStudent.PlanDetailId,
		Status:        coursePlanStudent.Status,
		Version:       coursePlanStudent.Version,
	}

	return reply, nil
}

// student
func (s *CourseService) ListCoursePlanStudent(ctx context.Context, req *pb.ListCoursePlanStudentRequest) (*pb.ListCoursePlanStudentReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.coursePlanStudentService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var coursePlanList = []*pb.ListCoursePlanStudentReply_Data{}
	for _, coursePlanStudent := range list {
		reply := &pb.ListCoursePlanStudentReply_Data{
			Id:            int32(coursePlanStudent.Id),
			CustomerPhone: coursePlanStudent.CustomerPhone,
			StudentId:     coursePlanStudent.StudentId,
			StudentName:   coursePlanStudent.StudentName,
			PlanId:        coursePlanStudent.PlanId,
			PlanName:      coursePlanStudent.PlanName,
			PlanDetailId:  coursePlanStudent.PlanDetailId,
			Status:        coursePlanStudent.Status,
			Version:       coursePlanStudent.Version,
			UpdatedAt:     coursePlanStudent.UpdatedAt,
		}

		coursePlanList = append(coursePlanList, reply)
	}
	return &pb.ListCoursePlanStudentReply{
		Data:  coursePlanList,
		Total: int32(count),
	}, nil
}

func (s *CourseService) StopCoursePlanStudent(ctx context.Context, req *pb.StopCoursePlanStudentRequest) (*pb.StopCoursePlanStudentReply, error) {

	student, err := s.coursePlanStudentService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("查找课程学生信息错误 %w", err)
	}

	if student.Status == int32(enum.CoursePlanStudentStatusType_STOPPED) {
		return nil, fmt.Errorf("学生已经停课")
	}

	detail, err := s.coursePlanDetailService.FindByID(ctx, student.PlanDetailId)
	if err != nil {
		return nil, fmt.Errorf("查找课程信息错误 %w", err)
	}

	if detail.Status == int32(enum.CoursePlanDetailStatusType_STOPPED) {
		return nil, fmt.Errorf("课程已经停课,无需操作")
	} // 此时状态不用一致性校验，因为后面什么都不操作，意味着，不管现在查到的状态是什么都不会影响数据的一致性

	err = s.tx.Transaction(ctx, func(txctx context.Context) error {

		// 更新课程,版本号() --> cas
		// 确保 刚查到的 课程信息没有被其他 操作更新 ( 为什么？ 1. 如果被其他事物修改，则修改数为0， 2.  由于在事务内，事物没有提交，不会释放锁)
		err = s.coursePlanDetailService.UpdateConcurrency(txctx, &domain.CoursePlanDetail{Id: detail.Id, Version: detail.Version})
		if err != nil {
			return fmt.Errorf("更新课程信息错误 %w", err)
		}

		err := s.coursePlanStudentService.UpdateConcurrency(txctx, &domain.CoursePlanStudent{Id: student.Id, Version: student.Version, Status: int32(enum.CoursePlanStudentStatusType_STOPPED)})
		if err != nil {
			return fmt.Errorf("更新课程学生信息错误 %w", err)
		}

		if detail.Status == int32(enum.CoursePlanDetailStatusType_SCHEDULED) {
			s.log.WithContext(txctx).Info("课程是已排课,更新学生状态,和每日课程里的学生状态")

			// 查找待上课的 日程(课
			dailyLessons, err := s.dailyLessonService.ListByMap(txctx, map[string]interface{}{"plan_detail_id": detail.Id, "status": int32(enum.DailyLessonStatusType_PENDING)})
			if err != nil {
				return fmt.Errorf("查找每日课程错误 %w", err)
			}

			// lessonIds := util.MapPtr[domain.DailyLesson, int32](dailyLessons, func(dl *domain.DailyLesson) int32 { return dl.Id })
			// 不批量更新，因为要做一致性检查
			for _, lesson := range dailyLessons {

				lessonStudents, err := s.dailyLessonStudentService.ListByMap(txctx, map[string]interface{}{"lesson_id": lesson.Id, "student_id": student.StudentId})
				if err != nil {
					return fmt.Errorf("查找课程学生错误 %w", err)
				}

				if len(lessonStudents) != 1 {
					return fmt.Errorf("查找课程学生错误,未查询到唯一的课程学生 %w", err)
				}
				lessonStudent := lessonStudents[0]

				err = s.dailyLessonStudentService.UpdateConcurrency(txctx, &domain.DailyLessonStudent{Id: lessonStudent.Id, Version: lessonStudent.Version, Status: int32(enum.DailyLessonStudentStatusType_CANCELED)})
				if err != nil {
					return err
				}

				err = s.dailyLessonService.UpdateConcurrency(txctx, &domain.DailyLesson{Id: lesson.Id, Version: lesson.Version})
				if err != nil {
					return err
				}

			}

		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 更新停课状态
	err = s.coursePlanStudentService.UpdateConcurrency(ctx, &domain.CoursePlanStudent{Id: student.Id, Version: student.Version, Status: int32(enum.CoursePlanStudentStatusType_STOPPED)})
	if err != nil {
		return nil, err
	}

	if student.Status == int32(enum.CoursePlanStudentStatusType_NEW) {
		return &pb.StopCoursePlanStudentReply{}, nil
	}

	return &pb.StopCoursePlanStudentReply{}, nil
}

func (s *CourseService) ListDailyLesson(ctx context.Context, req *pb.ListDailyLessonRequest) (*pb.ListDailyLessonReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.dailyLessonService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var planDetailList = []*pb.ListDailyLessonReply_Data{}
	for _, detail := range list {
		planDetailList = append(planDetailList, &pb.ListDailyLessonReply_Data{
			Id:          detail.Id,
			PlanId:      detail.PlanId,
			PlanName:    detail.PlanName,
			CourseCode:  detail.CourseCode,
			DateOfDay:   detail.DateOfDay.Format("2006-01-02"),
			StartTime:   detail.StartTime,
			EndTime:     detail.EndTime,
			RoomId:      detail.RoomId,
			RoomName:    detail.RoomName,
			TeacherId:   detail.TeacherId,
			TeacherName: detail.TeacherName,
			SubjectId:   detail.SubjectId,
			SubjectName: detail.SubjectName,
			GradeId:     detail.GradeId,
			GradeName:   detail.GradeName,
			LessonNum:   detail.LessonNum,
			ActualNum:   detail.ActualNum,
			SignNum:     detail.SignNum,
			Status:      detail.Status,
			Version:     detail.Version,
		})
	}
	return &pb.ListDailyLessonReply{
		Data:  planDetailList,
		Total: int32(count),
	}, nil
}

func (s *CourseService) StartDailyLesson(ctx context.Context, req *pb.StartDailyLessonRequest) (*pb.StartDailyLessonReply, error) {

	err := s.dailyLessonService.Start(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.StartDailyLessonReply{}, nil
}

func (s *CourseService) CancelDailyLesson(ctx context.Context, req *pb.CancelDailyLessonRequest) (*pb.CancelDailyLessonReply, error) {

	err := s.dailyLessonService.Cancel(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.CancelDailyLessonReply{}, nil
}

func (s *CourseService) FinishDailyLesson(ctx context.Context, req *pb.FinishDailyLessonRequest) (*pb.FinishDailyLessonReply, error) {

	err := s.dailyLessonService.Finish(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.FinishDailyLessonReply{}, nil
}

func (s *CourseService) ListDailyLessonStudent(ctx context.Context, req *pb.ListDailyLessonStudentRequest) (*pb.ListDailyLessonStudentReply, error) {

	var query domain.Expression
	err := json.Unmarshal([]byte(req.Expr), &query)
	if err != nil {
		return nil, err
	}

	list, count, err := s.dailyLessonStudentService.List(ctx, &query, &domain.Page{PageNo: int(req.PageNo), PageSize: int(req.PageSize)})
	if err != nil {
		return nil, err
	}

	var studentList = []*pb.ListDailyLessonStudentReply_Data{}
	for _, detail := range list {
		studentList = append(studentList, &pb.ListDailyLessonStudentReply_Data{
			Id:            detail.Id,
			PlanId:        detail.PlanId,
			PlanName:      detail.PlanName,
			PlanDetailId:  detail.PlanDetailId,
			LessonId:      detail.LessonId,
			CustomerPhone: detail.CustomerPhone,
			StudentId:     detail.StudentId,
			StudentName:   detail.StudentName,
			Status:        detail.Status,
			Version:       detail.Version,
			UpdatedAt:     detail.UpdatedAt,
		})
	}

	return &pb.ListDailyLessonStudentReply{
		Data:  studentList,
		Total: count,
	}, nil
}

func (s *CourseService) SignDailyLessonStudent(ctx context.Context, req *pb.SignDailyLessonStudentRequest) (*pb.SignDailyLessonStudentReply, error) {

	dailyLessonStudent, err := s.dailyLessonStudentService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if dailyLessonStudent.Status != int32(enum.DailyLessonStudentStatusType_UNSIGNED) {

		return nil, fmt.Errorf("学生状态不是未签到,不可签到 %d", req.Id)
	}

	err = s.dailyLessonStudentService.UpdateConcurrency(ctx, &domain.DailyLessonStudent{Id: req.Id, Version: dailyLessonStudent.Version, Status: int32(enum.DailyLessonStudentStatusType_SIGNED)})
	if err != nil {
		return nil, fmt.Errorf("更新学生状态错误 %w", err)
	}

	return &pb.SignDailyLessonStudentReply{}, nil
}

func (s *CourseService) LeaveDailyLessonStudent(ctx context.Context, req *pb.LeaveDailyLessonStudentRequest) (*pb.LeaveDailyLessonStudentReply, error) {

	dailyLessonStudent, err := s.dailyLessonStudentService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if dailyLessonStudent.Status != int32(enum.DailyLessonStudentStatusType_UNSIGNED) {

		return nil, fmt.Errorf("学生状态不是未签到,不可请假 %d", req.Id)
	}
	err = s.dailyLessonStudentService.UpdateConcurrency(ctx, &domain.DailyLessonStudent{Id: req.Id, Version: dailyLessonStudent.Version, Status: int32(enum.DailyLessonStudentStatusType_ABSENT)})
	if err != nil {
		return nil, fmt.Errorf("更新学生状态错误 %w", err)
	}
	return &pb.LeaveDailyLessonStudentReply{}, nil
}
