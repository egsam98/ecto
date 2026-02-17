package ecto

import (
	"encoding/json"
	"fmt"
)

var errRequired = Error("required")

type Error string

func Errorf(format string, args ...any) Error { return Error(fmt.Sprintf(format, args...)) }

func (e Error) Error() string { return string(e) }

type ListError []Error

func (e ListError) JSON() json.RawMessage {
	b, _ := json.Marshal(e)
	return b
}

func (e ListError) Error() string { return string(e.JSON()) }

type MapError map[string]error

func (e *MapError) Add(key string, err error) {
	if err == nil {
		return
	}
	if *e == nil {
		*e = make(map[string]error)
	}
	(*e)[key] = err
}

func (e MapError) JSON() json.RawMessage {
	b, _ := json.Marshal(e)
	return b
}

func (e MapError) Error() string { return string(e.JSON()) }
