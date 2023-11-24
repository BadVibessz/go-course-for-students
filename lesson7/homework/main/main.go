package main

import (
	"reflect"
)

func main() {

	type args struct {
		v any
	}

	d := struct {
		args args
	}{
		args{
			v: new(any),
		},
	}

	val := reflect.ValueOf(&d).Elem()

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)

		if !val.Type().Field(i).IsExported() {
			// err
			println("ERR")
		} else {
			println(f.String())
			println(string(val.Type().Field(i).Tag.Get("yaml")))
		}
	}

}
