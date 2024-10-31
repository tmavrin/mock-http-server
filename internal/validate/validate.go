package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/tmavrin/mock-http-server/internal/config"
)

var valid = validator.New()

func RequestBody(httpReq any, req *config.Request) error {
	if req.Validate != nil {
		return fieldsValidation(httpReq, req.Validate)
	}

	if len(req.Match) > 0 {
		return bodyMatching(httpReq, req)
	}

	return nil
}
