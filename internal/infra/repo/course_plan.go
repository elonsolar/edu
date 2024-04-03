package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/domain/enum"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type CoursePlanEntity struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	CycleType   int32
	ExcludeRule []struct {
		ExcludeType int32
		ExcludeDate []string // [周六], [2023-01-02,2023-03-02]
	} `gorm:"serializer:json"`
	Status int32

	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32 `gorm:"default:1"`
}

func (t *CoursePlanEntity) TableName() string {
	return "course_plan"

}

type coursePlanRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// FindByID implements domain.CoursePlanRepo.
func (up *coursePlanRepo) FindByID(ctx context.Context, id int32) (*domain.CoursePlan, error) {

	var coursePlan CoursePlanEntity

	err := up.repo.GetDBFromContext(ctx).First(&coursePlan, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.CoursePlan{
		Id:          int32(coursePlan.ID),
		Name:        coursePlan.Name,
		Description: coursePlan.Description,
		StartTime:   coursePlan.StartTime,
		EndTime:     coursePlan.EndTime,
		CycleType:   coursePlan.CycleType,
		ExcludeRule: coursePlan.ExcludeRule,
		Status:      coursePlan.Status,
		Version:     coursePlan.Version,
	}, nil
}

// List implements domain.CoursePlanRepo.
func (up *coursePlanRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.CoursePlan, error) {

	var coursePlanList []*CoursePlanEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&coursePlanList).Error
	if err != nil {
		return nil, err
	}

	return Map(coursePlanList, func(coursePlan *CoursePlanEntity) *domain.CoursePlan {
		return &domain.CoursePlan{
			Id:          int32(coursePlan.ID),
			Name:        coursePlan.Name,
			Description: coursePlan.Description,
			StartTime:   coursePlan.StartTime,
			EndTime:     coursePlan.EndTime,
			CycleType:   coursePlan.CycleType,
			ExcludeRule: coursePlan.ExcludeRule,
			Status:      coursePlan.Status,
			Version:     coursePlan.Version,
			UpdatedAt:   coursePlan.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}), nil

}

// Count implements domain.CoursePlanRepo.
func (up *coursePlanRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&CoursePlanEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.CoursePlanRepo.
func (up *coursePlanRepo) Save(ctx context.Context, coursePlan *domain.CoursePlan) (*domain.CoursePlan, error) {
	var coursePlanEntity = &CoursePlanEntity{
		Name:        coursePlan.Name,
		Description: coursePlan.Description,
		StartTime:   coursePlan.StartTime,
		EndTime:     coursePlan.EndTime,
		CycleType:   coursePlan.CycleType,
		ExcludeRule: coursePlan.ExcludeRule,
		Status:      int32(enum.CoursePlanStatusType_NEW),
	}
	err := up.repo.GetDBFromContext(ctx).Save(coursePlanEntity).Error
	if err != nil {
		return nil, err
	}
	coursePlan.Id = int32(coursePlanEntity.ID)
	return coursePlan, nil

}

// Update implements domain.CoursePlanRepo.
func (up *coursePlanRepo) Update(ctx context.Context, coursePlan *domain.CoursePlan) (int32, error) {

	var coursePlanEntity = &CoursePlanEntity{
		Name:        coursePlan.Name,
		Description: coursePlan.Description,
		StartTime:   coursePlan.StartTime,
		EndTime:     coursePlan.EndTime,
		CycleType:   coursePlan.CycleType,
		ExcludeRule: coursePlan.ExcludeRule,
		Status:      coursePlan.Status,
		Version:     coursePlan.Version,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&CoursePlanEntity{ID: uint(coursePlan.Id)}).Updates(coursePlanEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewCoursePlanRepo .
func NewCoursePlanRepo(repo *BaseRepo, logger log.Logger) domain.CoursePlanRepo {
	return &coursePlanRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
