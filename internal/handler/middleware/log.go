package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tmavrin/mock-http-server/internal/response"
)

func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}

		if len(body) == 0 {
			r.Body = io.NopCloser(bytes.NewReader(body))
			next.ServeHTTP(w, r)

			return
		}

		var txt json.RawMessage

		err = json.Unmarshal(body, &txt)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		}

		d, err := json.Marshal(txt)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		}

		r.Body = io.NopCloser(bytes.NewReader(body))

		l := fmt.Sprintf("[%s] request | %s %s | %s", time.Now().Format(time.RFC3339Nano), r.Method, r.URL.Path, string(d))

		log.Println(l)

		next.ServeHTTP(w, r)
	})
}
