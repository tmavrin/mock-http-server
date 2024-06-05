package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type preparedError struct {
	Status   int             `json:"status"`
	Method   string          `json:"method"`
	Path     string          `json:"path"`
	Response json.RawMessage `json:"response"`
}

var errorMap map[string]preparedError = map[string]preparedError{}

func prepareError() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var prepError preparedError

		err := json.NewDecoder(r.Body).Decode(&prepError)
		if err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if !strings.HasPrefix(prepError.Path, "/") {
			prepError.Path = "/" + prepError.Path
		}

		prepError.Method = strings.ToUpper(prepError.Method)

		errorMap[fmt.Sprintf("%s_%s", prepError.Method, prepError.Path)] = prepError
	})
}
