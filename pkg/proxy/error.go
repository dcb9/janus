package proxy

import (
	"encoding/json"
	"fmt"
)

const (
	ErrInvalid int = 150
)

type (
	ProxyError struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
	}

	ProxyErrors struct {
		errors []ProxyError
	}
)

func (e ProxyError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

func (errs ProxyErrors) Error() string {
	str := fmt.Sprintf("errors ")
	for _, e := range errs.errors {
		str += fmt.Sprint("%s ", e.Error())
	}
	return str
}

func (e ProxyErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.errors)
}
