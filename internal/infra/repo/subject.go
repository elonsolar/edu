package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/domain/enum"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type SubjectEntity struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string
	Category    int32
	Description string
	Status      int32

	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32 `gorm:"default:1"`
}

func (t *SubjectEntity) TableName() string {
	return "subject"

}

type subjectRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// FindByID implements domain.SubjectRepo.
func (up *subjectRepo) FindByID(ctx context.Context, id int32) (*domain.Subject, error) {

	var subject SubjectEntity

	err := up.repo.GetDBFromContext(ctx).First(&subject, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.Subject{
		Id:          int32(subject.ID),
		Name:        subject.Name,
		Category:    int32(subject.Category),
		Description: subject.Description,
		Status:      subject.Status,
		Version:     subject.Version,
	}, nil
}

// List implements domain.SubjectRepo.
func (up *subjectRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.Subject, error) {

	var subjectList []*SubjectEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&subjectList).Error
	if err != nil {
		return nil, err
	}

	return Map(subjectList, func(subject *SubjectEntity) *domain.Subject {
		return &domain.Subject{
			Id:          int32(subject.ID),
			Name:        subject.Name,
			Category:    subject.Category,
			Status:      subject.Status,
			Description: subject.Description,
			Version:     subject.Version,
			UpdatedAt:   subject.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}), nil

}

// Count implements domain.SubjectRepo.
func (up *subjectRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&SubjectEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.SubjectRepo.
func (up *subjectRepo) Save(ctx context.Context, subject *domain.Subject) (*domain.Subject, error) {
	var subjectEntity = &SubjectEntity{
		Name:        subject.Name,
		Category:    subject.Category,
		Description: subject.Description,
		Status:      int32(enum.EnableStatusEnabled),
	}
	err := up.repo.GetDBFromContext(ctx).Save(subjectEntity).Error
	if err != nil {
		return nil, err
	}
	subject.Id = int32(subjectEntity.ID)
	return subject, nil

}

// Update implements domain.SubjectRepo.
func (up *subjectRepo) Update(ctx context.Context, subject *domain.Subject) (int32, error) {

	var subjectEntity = &SubjectEntity{
		Name:        subject.Name,
		Category:    subject.Category,
		Description: subject.Description,
		Status:      subject.Status,
		Version:     subject.Version,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&SubjectEntity{ID: uint(subject.Id)}).Updates(subjectEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewSubjectRepo .
func NewSubjectRepo(repo *BaseRepo, logger log.Logger) domain.SubjectRepo {
	return &subjectRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
