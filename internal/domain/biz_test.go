package domain_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestXxx(t *testing.T) {
	funcType := reflect.FuncOf([]reflect.Type{reflect.TypeOf(int(0)), reflect.TypeOf(int(0))}, []reflect.Type{reflect.TypeOf(int(0))}, false)

	// addFunc := func(a, b int) int {
	// 	return a + b
	// }

	newFunc := reflect.MakeFunc(funcType, func(args []reflect.Value) (results []reflect.Value) {

		fmt.Println("what")
		return []reflect.Value{reflect.ValueOf(1)}
	})

	result := newFunc.Call([]reflect.Value{reflect.ValueOf(1), reflect.ValueOf(2)})

	fmt.Println(result)
}

// ptr
type emptyStruct struct{}

func TestCtx(t *testing.T) {

	ctx := context.Background()
	ctx = context.WithValue(ctx, emptyStruct{}, "ss")

	fmt.Println(ctx.Value(emptyStruct{}))

}
