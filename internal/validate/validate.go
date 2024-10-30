package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/tmavrin/mock-http-server/internal/config"
)

var valid = validator.New()

func RequestBody(httpReq map[string]any, req *config.Request) error {
	if len(req.Validate) > 0 {
		return fieldsValidation(httpReq, req)
	}

	if len(req.Match) > 0 {
		return bodyMatching(httpReq, req)
	}

	return nil
}

func fieldsValidation(httpReq map[string]any, req *config.Request) error {
	err := validateObject(httpReq, req.Validate)
	if err != nil {
		return err
	}

	return nil
}

func validateObject(httpReq map[string]any, req map[string]any) error {
	for k, v := range req {
		switch v := v.(type) {
		case string:
			if httpReq[k] == nil && !strings.Contains(v, "required") {
				break
			}

			err := valid.Var(httpReq[k], v)
			if err != nil {
				msg := err.Error()
				msg = strings.Replace(msg,
					"Key: '' Error:Field validation for ''",
					fmt.Sprintf("Key: '%s' Error:Field validation for '%s'", k, k),
					1,
				)

				return errors.New(msg)
			}

		case map[string]any:
			h, ok := httpReq[k].(map[string]any)
			if !ok {
				return fmt.Errorf("object of Key: '%s' not present", k)
			}

			err := validateObject(h, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func bodyMatching(httpReq map[string]any, req *config.Request) error {
	var expected any

	if err := json.Unmarshal([]byte(req.Match), &expected); err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(httpReq)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, httpReq) {
		return fmt.Errorf(
			"request does not match expected: %s , actual: %s",
			req.Match,
			string(jsonBytes),
		)
	}

	return nil
}
