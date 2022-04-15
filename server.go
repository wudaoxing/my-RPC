package myRPC

import "reflect"

type Server struct {
	addr  string
	funcs map[string]reflect.Value
}
