package mock

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/tmavrin/mock-http-server/internal/config"
	"github.com/tmavrin/mock-http-server/internal/response"
	"github.com/tmavrin/mock-http-server/internal/validate"
)

func Handler(h config.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prepError, ok := errorMap[fmt.Sprintf("%s_%s", h.Method, h.Path)]
		if ok {
			delete(errorMap, fmt.Sprintf("%s_%s", h.Method, h.Path))
			response.JSON(w, prepError.Status, prepError.Response)

			return
		}

		reqBody := map[string]any{}

		if r.Body != nil {
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			if err != nil && !errors.Is(err, io.EOF) {
				response.JSON(w, http.StatusInternalServerError, err.Error())

				return
			}
		}

		if h.Request != nil {
			err := validate.RequestBody(reqBody, h.Request)
			if err != nil {
				response.JSON(w, cmp.Or(http.StatusBadRequest, http.StatusInternalServerError), err.Error())

				return
			}
		}

		var responseData any

		if string(h.Response) != "" {
			err := json.Unmarshal(h.Response, &responseData)
			if err != nil {
				response.JSON(w, http.StatusInternalServerError, err.Error())

				return
			}
		}

		for k, v := range h.ResponseEcho {
			requestValue := findValue(reqBody, k)

			switch responseData.(type) {
			case []any:
				castArray, ok := responseData.([]any)
				if !ok {
					response.JSON(w, http.StatusInternalServerError, "can not cast response as array")

					return
				}

				for _, item := range castArray {
					castItem, ok := item.(map[string]any)
					if !ok {
						response.JSON(w, http.StatusInternalServerError, "can not cast response as object")

						return
					}

					setValue(castItem, v, requestValue)
				}
			case map[string]any:
				castItem, ok := responseData.(map[string]any)
				if !ok {
					response.JSON(w, http.StatusInternalServerError, "can not cast response as object")

					return
				}

				setValue(castItem, v, requestValue)
			default:
				response.JSON(w, http.StatusInternalServerError, "unknown request, not an object nor array")

				return
			}

		}

		response.JSON(w, http.StatusOK, responseData)
	})
}
