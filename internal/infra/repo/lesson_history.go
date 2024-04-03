package repo

import (
	"context"
	"edu/internal/domain"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type LessonHistoryEntity struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Mobile      string
	SourceType  int32
	OriginNum   int32
	NumChange   int32
	DetailId    int32
	Description string

	UserId   int // 和登录账号绑定
	TenantId int
}

func (t *LessonHistoryEntity) TableName() string {
	return "lesson_history"

}

type lessonHistoryRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// FindByID implements domain.LessonHistoryRepo.
func (up *lessonHistoryRepo) FindByID(ctx context.Context, id int32) (*domain.LessonHistory, error) {

	var lessonHistory LessonHistoryEntity

	err := up.repo.GetDBFromContext(ctx).First(&lessonHistory, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.LessonHistory{
		Id:          int32(lessonHistory.ID),
		Mobile:      lessonHistory.Mobile,
		NumChange:   lessonHistory.NumChange,
		SourceType:  lessonHistory.SourceType,
		DetailId:    lessonHistory.DetailId,
		Description: lessonHistory.Description,
	}, nil
}

// List implements domain.LessonHistoryRepo.
func (up *lessonHistoryRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.LessonHistory, error) {

	var lessonHistoryList []*LessonHistoryEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Order("updated_at desc").Find(&lessonHistoryList).Error
	if err != nil {
		return nil, err
	}

	return Map(lessonHistoryList, func(lessonHistory *LessonHistoryEntity) *domain.LessonHistory {
		return &domain.LessonHistory{
			Id:          int32(lessonHistory.ID),
			Mobile:      lessonHistory.Mobile,
			OriginNum:   lessonHistory.OriginNum,
			NumChange:   lessonHistory.NumChange,
			SourceType:  lessonHistory.SourceType,
			DetailId:    lessonHistory.DetailId,
			Description: lessonHistory.Description,
			UpdatedAt:   lessonHistory.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}), nil

}

// Count implements domain.LessonHistoryRepo.
func (up *lessonHistoryRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&LessonHistoryEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.LessonHistoryRepo.
func (up *lessonHistoryRepo) Save(ctx context.Context, lessonHistory *domain.LessonHistory) (*domain.LessonHistory, error) {
	var lessonHistoryEntity = &LessonHistoryEntity{
		Mobile:      lessonHistory.Mobile,
		OriginNum:   lessonHistory.OriginNum,
		NumChange:   lessonHistory.NumChange,
		SourceType:  lessonHistory.SourceType,
		DetailId:    lessonHistory.DetailId,
		Description: lessonHistory.Description,
	}
	err := up.repo.GetDBFromContext(ctx).Save(lessonHistoryEntity).Error
	if err != nil {
		return nil, err
	}
	lessonHistory.Id = int32(lessonHistoryEntity.ID)
	return lessonHistory, nil

}

// Update implements domain.LessonHistoryRepo.
func (up *lessonHistoryRepo) Update(ctx context.Context, lessonHistory *domain.LessonHistory) (int32, error) {
	var lessonHistoryEntity = &LessonHistoryEntity{
		Mobile:      lessonHistory.Mobile,
		NumChange:   lessonHistory.NumChange,
		SourceType:  lessonHistory.SourceType,
		DetailId:    lessonHistory.DetailId,
		Description: lessonHistory.Description,
		// LessonNumber: lessonHistory.LessonNumber,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&RoomEntity{ID: uint(lessonHistoryEntity.ID)}).Updates(lessonHistoryEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewLessonHistoryRepo .
func NewLessonHistoryRepo(repo *BaseRepo, logger log.Logger) domain.LessonHistoryRepo {
	return &lessonHistoryRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
