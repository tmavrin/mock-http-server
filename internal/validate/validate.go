package validate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/tmavrin/mock-http-server/internal/config"
)

var valid = validator.New()

func RequestBody(r *http.Request, req *config.Request) error {
	if len(req.Validate) > 0 {
		return fieldsValidation(r, req)
	}

	if len(req.Match) > 0 {
		return bodyMatching(r, req)
	}

	return nil
}

func fieldsValidation(r *http.Request, req *config.Request) error {
	httpReq := make(map[string]any)

	err := json.NewDecoder(r.Body).Decode(&httpReq)
	if err != nil {
		return err
	}

	for k, v := range req.Validate {
		if httpReq[k] == nil && !strings.Contains(v, "required") {
			continue
		}

		err = valid.Var(httpReq[k], v)
		if err != nil {
			return err
		}
	}

	return nil
}

func bodyMatching(r *http.Request, req *config.Request) error {
	var expected, actual any

	if err := json.Unmarshal([]byte(req.Match), &expected); err != nil {
		return err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &actual); err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf(
			"request does not match expected: %s , actual: %s",
			req.Match,
			body,
		)
	}

	return nil
}
