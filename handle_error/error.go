package errhandle



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
	s, _ := f(3)
	fmt.Println("s :::: ", s)

	f1 := val
	must(&f1, nil)
	f1()
	fmt.Println("f1 :::: ")
}

// must use to handle error, if error not nil, will panic.Or use callback.
func must(f interface{}, callback func(error)) {
	if reflect.TypeOf(f).Kind() != reflect.Ptr {
		return
	}
	var fv reflect.Value
	fv = reflect.ValueOf(reflect.ValueOf(f).Elem().Interface())
	errHandle := func(in []reflect.Value) []reflect.Value {
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

