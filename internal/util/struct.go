package util

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// 将一个结构体的属性，复制到另一个结构体中，支持嵌套

type copyer struct {
	config *Config
}
type Config struct {
	ignoreNotMatchedProperty bool
}

func IgnoreNotMatchedProperty() option {
	return func(c *Config) {
		c.ignoreNotMatchedProperty = true
	}
}

type option func(c *Config)

func (c *copyer) copy(src, dst reflect.Value) error {

	if !dst.IsValid() {
		if c.config.ignoreNotMatchedProperty {
			return nil
		} else {
			//TODO
			// return fmt.Errorf("not supported  src:%v dst:%v", src.Kind(), dst.Kind())
		}
	}

	switch src.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32:
		switch dst.Kind() {
		case reflect.Int, reflect.Int64, reflect.Int32:
			dst.Set(src)
		case reflect.String:
			dst.Set(reflect.ValueOf(fmt.Sprintf("%d", src.Interface())))
		default:
			if c.config.ignoreNotMatchedProperty {
				return nil
			}
			return fmt.Errorf("not supported type src:%v dst:%v", src.Kind(), dst.Kind())
		}

	case reflect.Float32:
		switch dst.Kind() {
		case reflect.Float32:
			dst.Set(src)
		case reflect.Float64:
			dst.Set(reflect.ValueOf((float64)(src.Interface().(float32))))
		default:
			if c.config.ignoreNotMatchedProperty {
				return nil
			}
			return fmt.Errorf("not supported type src:%v dst:%v", src.Kind(), dst.Kind())
		}
	case reflect.Float64:
		switch dst.Kind() {
		case reflect.Float64:
			dst.Set(src)
		case reflect.Float32:
			dst.Set(reflect.ValueOf((float32)(src.Interface().(float64))))
		default:
			if c.config.ignoreNotMatchedProperty {
				return nil
			}
			return fmt.Errorf("not supported type src:%v dst:%v", src.Kind(), dst.Kind())
		}

	case reflect.String:

		val, _ := src.Interface().(string)
		if val == "" { // ignore empty
			return nil
		}

		switch dst.Kind() {
		case reflect.String:
			dst.Set(src)
		case reflect.Int:
			srcVal, err := strconv.Atoi(src.Interface().(string))
			if err != nil {
				return err
			}
			dst.Set(reflect.ValueOf(srcVal))
		case reflect.Int32:
			srcVal, err := strconv.Atoi(src.Interface().(string))
			if err != nil {
				return err
			}
			dst.Set(reflect.ValueOf(int32(srcVal)))
		case reflect.Int64:
			srcVal, err := strconv.Atoi(src.Interface().(string))
			if err != nil {
				return err
			}
			dst.Set(reflect.ValueOf(int64(srcVal)))

		case reflect.Struct: // string ->time

			if _, ok := dst.Interface().(time.Time); ok {

				tm, err := time.ParseInLocation(val, "2006-01-02 15:04:05", time.Local)
				if err != nil && !c.config.ignoreNotMatchedProperty {
					return err
				}
				dst.Set(reflect.ValueOf(tm))
			}
		default:
			if c.config.ignoreNotMatchedProperty {
				return nil
			}
			return fmt.Errorf("not supported type src:%v dst:%v", src.Kind(), dst.Kind())
		}

	case reflect.Bool: // primitive type

		switch dst.Kind() {
		case reflect.Bool:
			dst.Set(src)
		case reflect.Int, reflect.Int32, reflect.Int64:
			srcVal := src.Interface().(bool)
			if srcVal {
				dst.Set(reflect.ValueOf(1))
			} else {
				dst.Set(reflect.ValueOf(0))
			}
		default:
			if c.config.ignoreNotMatchedProperty {
				return nil
			}
			return fmt.Errorf("not supported type src:%v dst:%v", src.Kind(), dst.Kind())
		}

	case reflect.Struct:
		if tm, ok := src.Interface().(time.Time); ok {

			switch dst.Interface().(type) {
			case time.Time:
				// dst.Set(dst)
				dst.Set(reflect.ValueOf(tm))
			case string:
				dst.Set(reflect.ValueOf(tm.Format("2006-01-02 15:04:05")))
			default:
				if c.config.ignoreNotMatchedProperty {
					return nil
				}
				return fmt.Errorf("not supported type src:%v dst:%v", src.Kind(), dst.Kind())
			}
			return nil
		}

		if dst.Kind() != reflect.Struct {
			if c.config.ignoreNotMatchedProperty {
				return nil
			}
			return fmt.Errorf("not supported type src:%v dst:%v", src.Kind(), dst.Kind())
		}
		for i := 0; i < src.NumField(); i++ {
			// fieldName := src.Type().Field(i).Name
			// fmt.Println("fieldName", fieldName)
			srcField := src.Field(i)
			dstField := dst.FieldByName(src.Type().Field(i).Name)

			if !src.Type().Field(i).IsExported() {
				continue
			}

			err := c.copy(srcField, dstField)
			if !c.config.ignoreNotMatchedProperty && err != nil {
				return err
			}
		}
	case reflect.Pointer: // *int -> int,*int
		if src.IsNil() {
			return nil
		}
		if dst.Kind() != reflect.Pointer {
			if !c.config.ignoreNotMatchedProperty {
				return fmt.Errorf("not supported type src:%v dst:%v", src.Kind(), dst.Kind())
			} else {
				return nil
			}
		}
		// dst.Set(reflect.New(dstType.Elem())) error _> can not set
		// val
		// val.Type *int
		// val.Type.Elem() -> int
		if dst.Kind() == reflect.Pointer && dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}

		err := c.copy(src.Elem(), dst.Elem())
		if err != nil {
			return err
		}
	case reflect.Slice:
		if dst.Kind() != reflect.Slice {
			if !c.config.ignoreNotMatchedProperty {
				return fmt.Errorf("not supported type src:%v dst:%v", src.Kind(), dst.Kind())
			} else {
				return nil
			}
		}

		dst.Set(reflect.MakeSlice(reflect.SliceOf(dst.Type().Elem()), src.Len(), src.Cap()))

		for i := 0; i < src.Len(); i++ {
			err := c.copy(src.Index(i), dst.Index(i))
			if !c.config.ignoreNotMatchedProperty && err != nil {
				return err
			}
		}
	}
	return nil
}

func CopyProperties(src, dst any, options ...option) error {

	var cfg = &Config{}
	for _, option := range options {
		option(cfg)
	}
	var copyer = &copyer{config: cfg}

	return copyer.copy(reflect.ValueOf(src), reflect.ValueOf(dst))
}
