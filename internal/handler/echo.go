package handler

import (
	"errors"
)

func setResponseEcho(config map[string]string, reqBody any, responseData any) error {
	for k, v := range config {
		requestValue := findValue(reqBody, k)

		switch data := responseData.(type) {
		case []any:
			for _, item := range data {
				castItem, ok := item.(map[string]any)
				if !ok {
					return errors.New("can not cast response as object")
				}

				setValue(castItem, v, requestValue)
			}
		case map[string]any:
			setValue(data, v, requestValue)
		default:
			return errors.New("unknown request, not an object nor array")
		}
	}

	return nil
}
