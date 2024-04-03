package util

import (
	"reflect"
)

// type TreeDataV1[I any] struct {
// 	Id       I
// 	Name     string
// 	ParentId I
// 	Children []*TreeData[I]
// }

// func BuildTreeV1[I comparable](elements []*TreeData[I], rootId I) []*TreeData[I] {

// 	var dataMap = make(map[I]*TreeData[I])
// 	for _, data := range elements {

// 		if current, ok := dataMap[data.Id]; ok {
// 			current.Name = data.Name
// 			current.ParentId = data.ParentId
// 			// !important ,确保 data ,和 current 是一致的
// 			data = current
// 		} else {
// 			dataMap[data.Id] = data
// 		}

// 		if parent, ok := dataMap[data.ParentId]; ok {
// 			parent.Children = append(parent.Children, data)
// 		} else {
// 			parent := &TreeData[I]{
// 				Id: data.ParentId,
// 			}
// 			parent.Children = append(parent.Children, data)
// 			dataMap[data.ParentId] = parent
// 		}

// 	}
// 	if data, ok := dataMap[rootId]; ok {
// 		return data.Children
// 	}

// 	return nil
// }

type treeBuilder[T any, I comparable] struct {
	cfg      *TreeConfig
	elements []*T
}

func (b *treeBuilder[T, I]) getId(t *T) I {
	return getValueByFieldName[T, I](t, b.cfg.idName)
}

func (b *treeBuilder[T, I]) setId(t *T, v I) {
	setValueByFieldName[T, I](t, b.cfg.idName, v)
}

func (b *treeBuilder[T, I]) getPid(t *T) I {
	return getValueByFieldName[T, I](t, b.cfg.pidName)
}

func (b *treeBuilder[T, I]) getChildren(t *T) []*T {

	return getValueByFieldName[T, []*T](t, b.cfg.childrenName)
}

func (b *treeBuilder[T, I]) setChildren(t *T, v []*T) {

	setValueByFieldName[T, []*T](t, b.cfg.childrenName, v)
}

type TreeConfig struct {
	idName       string
	pidName      string
	childrenName string
}

type TreeOption func(cfg *TreeConfig)

func WithIdName(idName string) TreeOption {
	return func(cfg *TreeConfig) {
		cfg.idName = idName
	}
}

func WithPidName(pidName string) TreeOption {
	return func(cfg *TreeConfig) {
		cfg.pidName = pidName
	}
}

func WithChildrenName(childRenName string) TreeOption {
	return func(cfg *TreeConfig) {
		cfg.childrenName = childRenName
	}
}

func NewTreeBuilder[T any, I comparable](options ...TreeOption) *treeBuilder[T, I] {

	var cfg = &TreeConfig{
		idName:       "Id",
		pidName:      "ParentId",
		childrenName: "Children",
	}
	for _, option := range options {
		option(cfg)
	}

	return &treeBuilder[T, I]{
		cfg: cfg,
	}
}

func (b *treeBuilder[T, I]) Build(elements []*T, rootId I) []*T {

	var dataMap = make(map[I]*T)
	for _, data := range elements {

		var id I = b.getId(data)
		var pid I = b.getPid(data)
		if current, ok := dataMap[id]; ok {
			children := b.getChildren(current)
			CopyProperties(data, current, IgnoreNotMatchedProperty())

			b.setChildren(current, children)
			// !important ,确保 data ,和 current 是一致的
			data = current
		} else {
			dataMap[id] = data
		}

		if parent, ok := dataMap[pid]; ok {
			children := b.getChildren(parent)
			children = append(children, data)
			b.setChildren(parent, children)
		} else {
			parent := new(T)
			b.setId(parent, pid)

			var children = []*T{data}
			dataMap[pid] = parent
			b.setChildren(parent, children)
		}

	}
	if data, ok := dataMap[rootId]; ok {
		return b.getChildren(data)
	}

	return []*T{}

}

func getValueByFieldName[T any, I any](t *T, name string) I {

	vf := reflect.ValueOf(t)
	for vf.Kind() == reflect.Pointer {
		vf = reflect.Indirect(vf)
	}
	fieldVal := vf.FieldByName(name)

	// fmt.Println(fieldVal.IsValid())
	if val, ok := fieldVal.Interface().(I); ok {
		return val
	}
	// haha
	var val I
	return val
}

func setValueByFieldName[T any, I any](t *T, name string, v I) {

	vf := reflect.ValueOf(t)
	for vf.Kind() == reflect.Pointer {
		vf = reflect.Indirect(vf)
	}
	fieldVal := vf.FieldByName(name)
	if _, ok := fieldVal.Interface().(I); ok {
		fieldVal.Set(reflect.ValueOf(v))
	}
}
