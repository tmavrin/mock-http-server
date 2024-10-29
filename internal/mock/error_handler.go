package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tmavrin/mock-http-server/internal/response"
)

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
