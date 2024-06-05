package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type request struct {
	Validate map[string]string `json:"validate"`
	Match    json.RawMessage   `json:"match"`
}

type handler struct {
	Path     string          `json:"path"`
	Method   string          `json:"method"`
	Response json.RawMessage `json:"response"`
	Request  request         `json:"request"`
}

func handleGet(h handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prepError, ok := errorMap[fmt.Sprintf("%s_%s", h.Method, h.Path)]
		if ok {
			delete(errorMap, fmt.Sprintf("%s_%s", h.Method, h.Path))
			writeJSONResponse(w, prepError.Status, prepError.Response)
			return
		}

		if string(h.Response) != "" {
			writeJSONResponse(w, http.StatusOK, h.Response)
			return
		}
	})
}

func handlePost(h handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prepError, ok := errorMap[fmt.Sprintf("%s_%s", h.Method, h.Path)]
		if ok {
			delete(errorMap, fmt.Sprintf("%s_%s", h.Method, h.Path))
			writeJSONResponse(w, prepError.Status, prepError.Response)
			return
		}

		status, data, err := handleRequestBody(r, h)
		if err != nil {
			writeJSONResponse(w, cmp.Or(status, http.StatusInternalServerError), err.Error())
			return
		}
		writeJSONResponse(w, status, cmp.Or[any](data, h.Response))
	})
}

func handleRequestBody(r *http.Request, h handler) (int, any, error) {
	if len(h.Request.Validate) > 0 {
		return handleBodyValidation(r, h)
	}

	return handleBodyMatching(r, h)
}

var valid = validator.New()

func handleBodyValidation(r *http.Request, h handler) (int, any, error) {
	req := make(map[string]any)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return 0, nil, err

	}

	for k, v := range h.Request.Validate {
		if req[k] == nil && !strings.Contains(v, "required") {
			continue
		}
		err = valid.Var(req[k], v)
		if err != nil {
			return http.StatusBadRequest, nil, err
		}
	}

	return http.StatusOK, h.Response, nil
}

func handleBodyMatching(r *http.Request, h handler) (int, any, error) {
	var expectedJSONAsInterface, actualJSONAsInterface any

	if err := json.Unmarshal([]byte(h.Request.Match), &expectedJSONAsInterface); err != nil {
		return 0, nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&actualJSONAsInterface); err != nil {
		return 0, nil, err
	}

	if !reflect.DeepEqual(expectedJSONAsInterface, actualJSONAsInterface) {
		return http.StatusBadRequest,
			fmt.Sprintf(
				"request does not match expected:\n expected: %s\n actual: %s",
				expectedJSONAsInterface,
				actualJSONAsInterface,
			),
			nil
	}
	return http.StatusOK, nil, nil
}
