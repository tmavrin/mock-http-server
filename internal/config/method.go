package config

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Method string

const (
	Get    Method = http.MethodGet
	Post   Method = http.MethodPost
	Put    Method = http.MethodPut
	Patch  Method = http.MethodPatch
	Delete Method = http.MethodDelete
)

func (m Method) String() string {
	return string(m)
}

func (m *Method) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	switch Method(s) {
	case Get, Post, Put, Patch, Delete:
		*m = Method(s)

		return nil
	default:
		return errors.New("invalid method")
	}
}
