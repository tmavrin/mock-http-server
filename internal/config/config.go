package config

import (
	"encoding/json"
)

const (
	Path             = "./config/mock-server.json"
	ErrorHandlerPath = "/prepare-error"
)

type Config struct {
	PrepareErrorPath string    `json:"prepareErrorPath" yaml:"prepareErrorPath"`
	Handlers         []Handler `json:"handlers" yaml:"handlers"`
}

type Handler struct {
	Path     string          `json:"path" yaml:"path"`
	Method   Method          `json:"method" yaml:"method"`
	Response json.RawMessage `json:"response" yaml:"response"`
	Request  *Request        `json:"request" yaml:"request"`
}

type Request struct {
	Validate map[string]string `json:"validate" yaml:"validate"`
	Match    json.RawMessage   `json:"match" yaml:"match"`
}
