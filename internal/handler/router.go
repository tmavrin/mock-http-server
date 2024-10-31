package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tmavrin/mock-http-server/internal/config"
	"github.com/tmavrin/mock-http-server/internal/handler/middleware"
)

func Router(cfg *config.Config) *http.ServeMux {
	router := http.NewServeMux()

	if cfg.PrepareErrorPath == "" {
		cfg.PrepareErrorPath = config.ErrorHandlerPath
	}

	router.Handle(
		fmt.Sprintf("%s %s", http.MethodPost, cfg.PrepareErrorPath),
		middleware.RequestLog(
			ErrorHandler(),
		),
	)

	for _, h := range cfg.Handlers {
		if !strings.HasPrefix(h.Path, "/") {
			h.Path = "/" + h.Path
		}

		router.Handle(
			fmt.Sprintf("%s %s", h.Method, h.Path),
			middleware.RequestLog(
				Handler(h),
			),
		)
	}

	return router
}
