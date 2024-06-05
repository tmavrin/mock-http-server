package main

import (
	"fmt"
	"net/http"
	"strings"
)

func parseRoutes(cfg config) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc(fmt.Sprintf("POST %s", cfg.PrepareErrorPath), prepareError)

	for _, h := range cfg.Handlers {
		if !strings.HasPrefix(h.Path, "/") {
			h.Path = "/" + h.Path
		}

		h.Method = strings.ToUpper(h.Method)

		switch h.Method {
		case http.MethodPost:
			router.HandleFunc(fmt.Sprintf("POST %s", h.Path), handlePost(h))

		case http.MethodGet:
			router.HandleFunc(fmt.Sprintf("GET %s", h.Path), handleGet(h))

		case http.MethodPut:
			router.HandleFunc(fmt.Sprintf("PUT %s", h.Path), handlePost(h))

		case http.MethodDelete:
			router.HandleFunc(fmt.Sprintf("DELETE %s", h.Path), handleGet(h))
		}
	}

	return router
}
