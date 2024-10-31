package handler

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
		// check for prepared error
		prepError, ok := errorMap[fmt.Sprintf("%s_%s", h.Method, h.Path)]
		if ok {
			delete(errorMap, fmt.Sprintf("%s_%s", h.Method, h.Path))
			response.JSON(w, prepError.Status, prepError.Response)

			return
		}

		var reqBody any

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

		if len(h.ResponseEcho) > 0 {
			err := setResponseEcho(h.ResponseEcho, reqBody, responseData)
			if err != nil {
				response.JSON(w, cmp.Or(http.StatusBadRequest, http.StatusInternalServerError), err.Error())

				return
			}
		}

		response.JSON(w, http.StatusOK, responseData)
	})
}
