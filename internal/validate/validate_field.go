package validate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func fieldsValidation(httpReq any, req any) error {
	switch validateConfig := req.(type) {
	case []any:
		if httpReq != nil {
			httpReq, ok := httpReq.([]any)
			if !ok {
				return errors.New("config request validation is defined as array, but sent request is not array")
			}

			return validateArray(httpReq, validateConfig)
		} else {
			return validateArray(nil, validateConfig)
		}

	case []map[string]any:
		if httpReq != nil {
			httpReq, ok := httpReq.([]map[string]any)
			if !ok {
				return errors.New("config request validation is defined as array, but sent request is not array")
			}

			return validateArrayMap(httpReq, validateConfig)
		} else {
			return validateArrayMap(nil, validateConfig)
		}

	case map[string]any:
		if httpReq != nil {
			httpReq, ok := httpReq.(map[string]any)
			if !ok {
				return errors.New("config request validation is defined as an object, but sent request is not an object")
			}

			return validateObject(httpReq, validateConfig)
		} else {
			return validateObject(nil, validateConfig)
		}
	}

	return nil
}

func validateArray(httpReq []any, req []any) error {
	for k, v := range req {
		switch v := v.(type) {
		// if we have string, that means its validation
		case string:
			if httpReq[k] == nil && !strings.Contains(v, "required") {
				break
			}

			if err := validateProperty(strconv.Itoa(k), httpReq[k], v); err != nil {
				return err
			}

		// we have subproperty as object or array
		case map[string]any, []any:
			if err := fieldsValidation(httpReq[k], v); err != nil {
				return err
			}

		default:
			return errors.New("validation rules must be a string")
		}
	}

	return nil
}

func validateArrayMap(httpReq []map[string]any, req []map[string]any) error {
	for k, v := range req {
		if err := fieldsValidation(httpReq[k], v); err != nil {
			return err
		}
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

			if err := validateProperty(k, httpReq[k], v); err != nil {
				return err
			}

		case map[string]any, []any, []map[string]any:
			return fieldsValidation(httpReq[k], v)

		default:
			return errors.New("validation rules must be a string")
		}
	}

	return nil
}

func validateProperty(key string, value any, rule string) error {
	err := valid.Var(value, rule)
	if err != nil {
		msg := err.Error()
		msg = strings.Replace(msg,
			"Key: '' Error:Field validation for ''",
			fmt.Sprintf("Key: '%s' Error:Field validation for '%s'", key, key),
			1,
		)

		return errors.New(msg)
	}

	return nil
}
