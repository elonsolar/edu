package domain

import (
	"context"
	"edu/internal/domain/enum"
	"fmt"
	"strconv"
	"time"
)

type MetaRepo interface {
	BaseRepo[Meta]
}

type Meta struct {
	Id        int32
	DataType  int32
	Name      string
	Value     string
	UpdatedAt string
	Version   int32
}

type MetaService struct {
	BaseService[Meta]
	repo MetaRepo
}

func (m *MetaService) GetLessonSequenceCode(ctx context.Context) (string, error) {

	metaList, err := m.ListByMap(ctx, map[string]interface{}{"data_type": enum.CourseCodeSequence, "name": time.Now().Year()})
	if err != nil {
		return "", fmt.Errorf("查询课程编号错误 err:%w", err)
	}
	var meta *Meta
	if len(metaList) == 0 {
		meta, err = m.Create(ctx, &Meta{
			DataType: int32(enum.CourseCodeSequence),
			Name:     fmt.Sprintf("%d", time.Now().Year()),
			Value:    "1",
			Version:  1,
		})
		if err != nil {
			return "", fmt.Errorf("新增编号错误 err:%w", err)
		}
	} else {
		meta = metaList[0]
	}

	code, err := strconv.Atoi(meta.Value)
	if err != nil {
		return "", fmt.Errorf("课程版本号类型转换错误 err:%w", err)
	}

	num, err := m.Update(ctx, &Meta{Id: meta.Id, Version: meta.Version, Value: fmt.Sprintf("%d", code+1)})
	if err != nil {
		return "", fmt.Errorf("课程序列号更新失败 err:%w", err)
	}
	if num != 1 {
		return "", fmt.Errorf("课程序列号更新失败 err:%w", ErrCurrencyUpdate)
	}
	return fmt.Sprintf("%d%03d", time.Now().Year(), code), nil
}

func NewMetaService(repo MetaRepo) *MetaService {

	return &MetaService{
		BaseService: BaseService[Meta]{
			repo: repo,
		},
		repo: repo,
	}
}
