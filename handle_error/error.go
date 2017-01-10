package errhandle

import (
	"fmt"
	"reflect"
	"testing"
)

func val() (int, error) {
	return 1, fmt.Errorf("ttest error")
}

func read(i int) (string, error) {
	return fmt.Sprint(i), fmt.Errorf("error in read")
	// return fmt.Sprint(i), nil
}

func TestMust(t *testing.T) {

	f := read
	must(&f, func(err error) {
		fmt.Println("there is an error")
	})
	s, _ := f(3) // 需要在外面调用，要不不太好确定输入参数,这样就可以忽略error了
	fmt.Println("s :::: ", s)

	f1 := val
	must(&f1, nil)
	f1() // 需要在外面调用，要不不太好确定输入参数
	fmt.Println("f1 :::: ")
}

// f 必须为函数的指针，因为要新建个函数，必须要传入指针。callback为具体的error处理方法
func must(f interface{}, callback func(error)) {
	if reflect.TypeOf(f).Kind() != reflect.Ptr {
		return
	}
	var fv reflect.Value
	fv = reflect.ValueOf(reflect.ValueOf(f).Elem().Interface())
	errHandle := func(in []reflect.Value) []reflect.Value { // 这个函数是新生成的函数的具体实现。这里我们只要处理error,其他保持和源函数一致
		results := fv.Call(in)
		if len(results) > 0 {
			errV := results[len(results)-1]
			if errV.Type().Name() == "error" && !errV.IsNil() {
				if callback != nil {
					callback(errV.Interface().(error))
				} else {
					panic(errV.Interface())
				}
			}
		}
		return results
	}
	makeErrHandle := func(fptr interface{}) {
		fn := reflect.ValueOf(fptr).Elem()

		v := reflect.MakeFunc(fn.Type(), errHandle)
		fn.Set(v)
	}

	makeErrHandle(f)

}
