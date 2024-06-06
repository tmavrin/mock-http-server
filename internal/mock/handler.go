package mock

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

		if h.Request != nil {
			err := validate.RequestBody(r, h.Request)
			if err != nil {
				response.JSON(w, cmp.Or(http.StatusBadRequest, http.StatusInternalServerError), err.Error())

				return
			}
		}

		if string(h.Response) != "" {
			response.JSON(w, http.StatusOK, h.Response)

			return
		}
	})
}

var errorMap = map[string]preparedError{}

type preparedError struct {
	Status   int             `json:"status"`
	Method   string          `json:"method"`
	Path     string          `json:"path"`
	Response json.RawMessage `json:"response"`
}

func ErrorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var prepError preparedError

		err := json.NewDecoder(r.Body).Decode(&prepError)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, err.Error())

			return
		}

		if !strings.HasPrefix(prepError.Path, "/") {
			prepError.Path = "/" + prepError.Path
		}

		prepError.Method = strings.ToUpper(prepError.Method)

		errorMap[fmt.Sprintf("%s_%s", prepError.Method, prepError.Path)] = prepError
	})
}
