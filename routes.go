package main

import (
	"fmt"
	"net/http"
	"strings"
)

func parseRoutes(cfg config) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle(fmt.Sprintf("POST %s", cfg.PrepareErrorPath), logMiddleware(prepareError()))

	for _, h := range cfg.Handlers {
		if !strings.HasPrefix(h.Path, "/") {
			h.Path = "/" + h.Path
		}

		h.Method = strings.ToUpper(h.Method)

		switch h.Method {
		case http.MethodPost:
			router.Handle(fmt.Sprintf("POST %s", h.Path), logMiddleware(handlePost(h)))

		case http.MethodGet:
			router.Handle(fmt.Sprintf("GET %s", h.Path), logMiddleware(handleGet(h)))

		case http.MethodPut:
			router.Handle(fmt.Sprintf("PUT %s", h.Path), logMiddleware(handlePost(h)))

		case http.MethodDelete:
			router.Handle(fmt.Sprintf("DELETE %s", h.Path), logMiddleware(handleGet(h)))
		}
	}

	return router
}
