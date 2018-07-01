package main

import "reflect"

func walk(x interface{}, fn func(input string)) {
	val := reflect.ValueOf(x)

	for i:=0; i<val.NumField(); i++ {
		field := val.Field(i)
		fn(field.String())
	}
}
