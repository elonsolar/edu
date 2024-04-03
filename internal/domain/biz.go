package domain

import (
	"context"
	"errors"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewUserService, NewTenantService, NewTeacherService, NewStudentService, NewRoomService, NewSubjectService, NewCustomerService, NewLessonHistoryService, NewCoursePlanService, NewCoursePlanDetailService, NewCoursePlanStudentService, NewDailyLessonService, NewDailyLessonStudentService, NewMetaService, NewPermissionService, NewRoleService, NewTaskService, NewSkuService, NewCombineSkuService, NewCombineSkuItemService)

type Tx interface {
	Transaction(ctx context.Context, fn func(txctx context.Context) error) error
}

var (
	ErrUnknown        = errors.New("失败")
	ErrRecordNotFound = errors.New("记录不存在")
	ErrCurrencyUpdate = errors.New("更新并发")
)

type BaseRepo[T any] interface {
	Save(context.Context, *T) (*T, error)
	BatchSave(context.Context, []*T) error
	Update(context.Context, *T) (int32, error)
	Delete(context.Context, int32) (int32, error)
	FindByID(context.Context, int32) (*T, error)
	List(context.Context, *Expression, *Page) ([]*T, error)
	Count(context.Context, *Expression) (int32, error)
	ListAll(context.Context, *Expression) ([]*T, error)
	ListByMap(context.Context, map[string]interface{}) ([]*T, error)
}

type BaseService[T any] struct {
	repo BaseRepo[T]
}

func (t *BaseService[T]) Create(ctx context.Context, model *T) (*T, error) {

	return t.repo.Save(ctx, model)
}

func (t *BaseService[T]) BatchCreate(ctx context.Context, model []*T) error {

	return t.repo.BatchSave(ctx, model)
}

func (t *BaseService[T]) Update(ctx context.Context, model *T) (int32, error) {

	return t.repo.Update(ctx, model)
}

// 这个方法会判断更新数量，如果为0 会抛出错误， 不用上层方法继续判断
func (t *BaseService[T]) UpdateConcurrency(ctx context.Context, model *T) error {

	num, err := t.repo.Update(ctx, model)
	if err != nil {
		return err
	}
	if num == 0 {
		return ErrCurrencyUpdate
	}
	return nil
}
func (t *BaseService[T]) FindByID(ctx context.Context, id int32) (*T, error) {

	return t.repo.FindByID(ctx, id)
}

func (t *BaseService[T]) Delete(ctx context.Context, id int32) (int32, error) {
	count, err := t.repo.Delete(ctx, id)
	if err != nil {
		return count, err
	}
	if count == 0 {
		return 0, ErrUnknown
	}
	return count, err
}

func (t *BaseService[T]) List(ctx context.Context, query *Expression, page *Page) ([]*T, int32, error) {

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

func (t *BaseService[T]) ListAll(ctx context.Context, query *Expression) ([]*T, error) {

	list, err := t.repo.ListAll(ctx, query)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (t *BaseService[T]) ListByMap(ctx context.Context, params map[string]interface{}) ([]*T, error) {

	list, err := t.repo.ListByMap(ctx, params)
	if err != nil {
		return nil, err
	}
	return list, nil
}
