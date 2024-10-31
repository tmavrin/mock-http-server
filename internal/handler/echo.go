package handler

import (
	"errors"
)

func setResponseEcho(config map[string]string, reqBody any, responseData any) error {
	for k, v := range config {
		requestValue := findValue(reqBody, k)

		switch responseData.(type) {
		case []any:
			castArray, ok := responseData.([]any)
			if !ok {
				return errors.New("can not cast response as array")
			}

			for _, item := range castArray {
				castItem, ok := item.(map[string]any)
				if !ok {
					return errors.New("can not cast response as object")
				}

				setValue(castItem, v, requestValue)
			}
		case map[string]any:
			castItem, ok := responseData.(map[string]any)
			if !ok {
				return errors.New("can not cast response as object")
			}

			setValue(castItem, v, requestValue)
		default:
			return errors.New("unknown request, not an object nor array")
		}
	}

	return nil
}
