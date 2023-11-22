package main

import (
	"reflect"
)

func main() {

	d := struct {
		private int    `yaml:"private"`
		Public  string `yaml:"private,dura"`
	}{
		private: 1,
		Public:  "0",
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
