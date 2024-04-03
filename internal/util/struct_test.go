package util

import (
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

// type MyInt int

type Page struct {
	PageNo int
}

type Request struct {
	Name      string
	Age       *int
	Page      *Page
	Page2     *Page
	Status    []int
	BirthDay  time.Time
	StartTime time.Time
}

type Page2 struct {
	PageNo int
}

type Request2 struct {
	Name  string
	Age   *int
	Page  *Page2
	Page2 *struct {
		PageNo int
	}
	Status    []int
	BirthDay  time.Time
	StartTime string
}

// go  test ./...
// go test xx/...
// go test xx...
func TestCopy(t *testing.T) {

	t.Run("复制结构体内 基本类型字段", func(t *testing.T) {

		// var r1 = &Request{Name: "xx", Age: &age, Page: &Page{PageNo: 10}, Status: []int{1, 2, 3, 4}}
		var r1 = &Request{Name: "xx"}

		var r2 = &Request2{}
		err := CopyProperties(r1, r2, IgnoreNotMatchedProperty())

		if err != nil {
			panic(err)
		}

		if r2.Name != "xx" {
			t.Fatalf("expected  xx , but got %s", r2.Name)
		}
	})

	t.Run("复制结构体内 基本类型指针字段", func(t *testing.T) {
		var age = 22
		var r1 = &Request{Age: &age}

		var r2 = &Request2{}
		CopyProperties(r1, r2)

		if r2.Age == nil || *r2.Age != 22 {
			t.Fatalf("expected  22 , but got %v", r2.Age)
		}
	})

	t.Run("复制结构体内 结构体指针字段", func(t *testing.T) {
		var r1 = &Request{Page: &Page{PageNo: 20}}

		var r2 = &Request2{}
		CopyProperties(r1, r2)

		if r2.Page == nil || *&r2.Page.PageNo != 20 {
			t.Fatalf("expected  20 , but got %v", r2.Page)
		}
	})

	t.Run("复制结构体内 切片字段", func(t *testing.T) {
		var r1 = &Request{Status: []int{1, 2, 3}}

		var r2 = &Request2{}
		CopyProperties(r1, r2)

		if len(r2.Status) != 3 {
			t.Fatalf("expected len is  3 , but got %d", len(r2.Status))
		}
	})

	t.Run("复制基本类型测试Int", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("[TODO ]int 不可复制类型，不应该复制和报错")
			}
		}()
		var o int = 1
		var b int = 2

		CopyProperties(o, b)
	})

	t.Run("复制匿名结构体", func(t *testing.T) {
		var r1 = &Request{Page2: &Page{
			PageNo: 10,
		}}

		var r2 = &Request2{}
		err := CopyProperties(r1, r2, IgnoreNotMatchedProperty())
		if err != nil {
			panic(err)
		}

		if r2.Page2.PageNo == 0 {
			t.Fatalf("expected pageNo is  0 , but got %d", r2.Page2.PageNo)
		}
	})

	t.Run("复制time->time", func(t *testing.T) {
		var r1 = &Request{BirthDay: time.Now()}

		var r2 = &Request2{}
		err := CopyProperties(r1, r2, IgnoreNotMatchedProperty())
		if err != nil {
			panic(err)
		}

		if r2.BirthDay.IsZero() {
			t.Fatal("expected BirthDay is not zero ,bug got zero ")

		}
	})

	t.Run("复制time->string", func(t *testing.T) {
		var r1 = &Request{StartTime: time.Now()}

		var r2 = &Request2{}
		err := CopyProperties(r1, r2, IgnoreNotMatchedProperty())
		if err != nil {
			panic(err)
		}

		if r2.StartTime == "" {
			t.Fatal("expected StartTime is not empty ,bug got empty ")
		}
	})
}

func copyAny(sv, dv reflect.Value) {

	if sv.Kind() != dv.Kind() { // TODO int <->int64
		return
	}
	dv.CanSet()

	switch sv.Kind() {
	case reflect.Struct:
		copyStruct(sv, dv)
	case reflect.Pointer:
		if sv.IsNil() { // 空值 忽略
			return
		}
		copyPtr(sv, dv)
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.String, reflect.Float64:
		dv.Set(sv)
	case reflect.Slice:
		copySlice(sv, dv)
	}

}

// []int -> []int
// []User{}->[]User
// []*User ->[]*User
func copySlice(src, dst reflect.Value) {
	if src.Len() == 0 {
		return
	}

	// srcEleType := src.Type().Kind()

	// first ->  reflect.Sliceof()
	dst.Set(reflect.MakeSlice(reflect.SliceOf(dst.Type().Elem()), src.Len(), src.Cap()))

	for i := 0; i < src.Len(); i++ {
		copyAny(src.Index(i), dst.Index(i))
	}

}

// v1 支持
// *int ->*int
// *Page1 ->Page2
func copyPtr(src, dst reflect.Value) { //ret ,  dst 是nil 的情况
	// 如果指针的底层类型相同

	switch src.Elem().Kind() {

	case reflect.Struct:

		if dst.IsNil() {
			// 空的指针， dst.Elem.Type 是 invalid , dst.Type.Elem() 	是可以的
			dst.Set(reflect.New(dst.Type().Elem()))
		}

		copyStruct(src.Elem(), dst.Elem())

	case reflect.Int, reflect.Int32, reflect.Int64:

		dst.Type().Elem().String()
		fmt.Println(dst.Type().Name())

		dst.Set(src)

	default:
		fmt.Println("what")
	}

}

func copyStruct(src, dst reflect.Value) {

	for i := 0; i < src.NumField(); i++ {
		if src.Field(i).Kind() == reflect.Ptr && src.Field(i).IsNil() {
			continue
		}
		fmt.Println(src.Type().Field(i).Name)
		copyAny(src.Field(i), dst.FieldByName(src.Type().Field(i).Name))
	}

}

// 1. 指针成员如何操纵值
func TestStruct(t *testing.T) {

	// var r1 = &Request{Name: "xx", Page: &Page{PageNo: 10}}
	var r1 = &Request{Name: "xx"}

	rv := reflect.ValueOf(r1)

	rv = reflect.Indirect(rv)

	for i := 0; i < rv.NumField(); i++ {

		val := rv.Field(i)
		fieldStruct := rv.Type().Field(i)
		fmt.Println(fieldStruct)

		// *int 赋值
		if val.Kind() == reflect.Ptr && val.IsNil() {
			ele := reflect.Indirect(val)
			fmt.Println(ele.Kind()) // 用 value 取 elem 是invalid
			switch ele.Kind() {
			case reflect.Int:
				ele.Set(reflect.ValueOf(1))

			}

			// kd := fieldStruct.Type.Elem().Kind() // int
			fmt.Println(fieldStruct.Type.Elem())

			nv := reflect.New(fieldStruct.Type.Elem())
			val.Set(nv)
			// switch kd {

			// case reflect.Int:
			// 	var a int = 1
			// 	val.Set(reflect.ValueOf(&a))
			// }

			continue
		}
		if val.Kind() == reflect.String {
			val.Set(reflect.ValueOf("aa"))
		}

	}

	fmt.Println(r1)

}

// func setDefault()

func TestCopyByUnsafe(t *testing.T) {

	var r1 = &Request{Name: "xx"}

	r2 := (*Request2)(unsafe.Pointer(r1))

	fmt.Printf("%v", r2)

}

func TestNewStruct(t *testing.T) {

	var r = new(Request)

	rv := reflect.ValueOf(r)
	fmt.Println(rv.IsNil())
}

func TestSet(t *testing.T) {

	type User struct {
		Name string
	}

	var user = &User{Name: "John Doe"}
	var uv = reflect.ValueOf(user).Elem()
	fmt.Println(uv.CanSet())
}

func TestIndirect(t *testing.T) {

	var age *int
	v := reflect.ValueOf(&age)
	v = reflect.Indirect(v)
	fmt.Println(v.Kind() == reflect.Pointer)
}
