package main

import (
	"canoe/internal/model/data"
	"reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	val := reflect.ValueOf(data.User{})
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		println(f.String())
	}
}

func TestCreateUser(t *testing.T) {

}
