package repo

import (
	"context"
	"edu/internal/domain"
	"edu/internal/domain/enum"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type RoomEntity struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Code        string
	Subjects    []int32 `gorm:"serializer:json"`
	Description string
	Status      int32

	UserId   int // 和登录账号绑定
	TenantId int
	Version  int32 `gorm:"default:1"`
}

func (t *RoomEntity) TableName() string {
	return "room"

}

type roomRepo struct {
	repo *BaseRepo
	log  *log.Helper
}

// FindByID implements domain.RoomRepo.
func (up *roomRepo) FindByID(ctx context.Context, id int32) (*domain.Room, error) {

	var room RoomEntity

	err := up.repo.GetDBFromContext(ctx).First(&room, id).Error

	if err != nil {
		return nil, err
	}

	return &domain.Room{
		Id:          int32(room.ID),
		Code:        room.Code,
		Subjects:    room.Subjects,
		Description: room.Description,
		Status:      room.Status,
		Version:     room.Version,
	}, nil
}

// List implements domain.RoomRepo.
func (up *roomRepo) List(ctx context.Context, query *domain.Expression, page *domain.Page) ([]*domain.Room, error) {

	var roomList []*RoomEntity
	err := up.repo.GetDBFromContext(ctx).Scopes(Paginate(page)).Where(query2Clause(query)).Find(&roomList).Error
	if err != nil {
		return nil, err
	}

	return Map(roomList, func(room *RoomEntity) *domain.Room {
		return &domain.Room{
			Id:          int32(room.ID),
			Code:        room.Code,
			Subjects:    room.Subjects,
			Status:      room.Status,
			Description: room.Description,
			Version:     room.Version,
			UpdatedAt:   room.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}), nil

}

// Count implements domain.RoomRepo.
func (up *roomRepo) Count(ctx context.Context, query *domain.Expression) (int32, error) {
	var count int64
	err := up.repo.GetDBFromContext(ctx).Model(&RoomEntity{}).Where(query2Clause(query)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// Save implements domain.RoomReppo.
func (up *roomRepo) Save(ctx context.Context, room *domain.Room) (*domain.Room, error) {
	var roomEntity = &RoomEntity{
		Code:        room.Code,
		Subjects:    room.Subjects,
		Description: room.Description,
		Status:      int32(enum.EnableStatusEnabled),
	}
	err := up.repo.GetDBFromContext(ctx).Create(roomEntity).Error

	if err != nil {
		return nil, err
	}
	room.Id = int32(roomEntity.ID)
	return room, nil

}

// Update implements domain.RoomRepo.
func (up *roomRepo) Update(ctx context.Context, room *domain.Room) (int32, error) {

	var roomEntity = &RoomEntity{
		Code:        room.Code,
		Subjects:    room.Subjects,
		Status:      room.Status,
		Description: room.Description,
		Version:     room.Version,
	}

	result := up.repo.GetDBFromContext(ctx).Model(&RoomEntity{ID: uint(room.Id)}).Updates(roomEntity)
	if result.Error != nil {
		return 0, result.Error
	}
	up.log.Error("updateCount", result.RowsAffected)
	return int32(result.RowsAffected), nil
}

// NewRoomRepo .
func NewRoomRepo(repo *BaseRepo, logger log.Logger) domain.RoomRepo {
	return &roomRepo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
