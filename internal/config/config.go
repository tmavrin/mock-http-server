package config

import (
	"encoding/json"
)

const (
	Path             = "./example/mock-server.json"
	ErrorHandlerPath = "/prepare-error"
)

type Config struct {
	PrepareErrorPath string    `json:"prepareErrorPath" yaml:"prepareErrorPath"`
	Handlers         []Handler `json:"handlers" yaml:"handlers"`
}

type Handler struct {
	Path         string                `json:"path" yaml:"path"`
	Method       Method                `json:"method" yaml:"method"`
	Response     json.RawMessage       `json:"response" yaml:"response"`
	ResponseEcho map[string]string     `json:"responseEcho" yaml:"responseEcho"`
	Request      *Request              `json:"request" yaml:"request"`
	QueryParams  map[string]QueryParam `json:"queryParams" yaml:"queryParams"`
}

type Request struct {
	Validate any             `json:"validate" yaml:"validate"`
	Match    json.RawMessage `json:"match" yaml:"match"`
}

type QueryParam struct {
	Required bool   `json:"required" yaml:"required"`
	Value    string `json:"value" yaml:"value"`
}
