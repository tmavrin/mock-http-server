package mock

import (
	"cmp"
	"encoding/json"
	"fmt"
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

		var reqBody map[string]any

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		}

		if h.Request != nil {
			err := validate.RequestBody(reqBody, h.Request)
			if err != nil {
				response.JSON(w, cmp.Or(http.StatusBadRequest, http.StatusInternalServerError), err.Error())

				return
			}
		}

		responseData := make(map[string]any)

		if string(h.Response) != "" {
			err := json.Unmarshal(h.Response, &responseData)
			if err != nil {
				response.JSON(w, http.StatusInternalServerError, err.Error())
			}
		}

		for k, v := range h.ResponseEcho {
			requestValue := findValue(reqBody, k)

			setValue(responseData, v, requestValue)
		}

		response.JSON(w, http.StatusOK, responseData)
	})
}
