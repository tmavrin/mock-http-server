package handler

import (
	"fmt"
	"net/http"

	"github.com/tmavrin/mock-http-server/internal/config"
)

func checkQueryParams(r *http.Request, config map[string]config.QueryParam) error {
	for key, c := range config {
		if c.Required && !r.URL.Query().Has(key) {
			return fmt.Errorf("required query param %s is not present", key)
		} else if !c.Required && r.URL.Query().Get(key) == "" {
			return nil
		} else if c.Value != "" && r.URL.Query().Get(key) != c.Value {
			return fmt.Errorf("query param %s value '%s' does not match '%s'", key, r.URL.Query().Get(key), c.Value)
		}
	}

	return nil
}
