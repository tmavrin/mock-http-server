package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		_, err = w.Write(b)
		if err != nil {
			log.Printf("writing response: %s\n", err)
		}

		l := fmt.Sprintf("[%s] response | %d | %s", time.Now().Format(time.RFC3339Nano), statusCode, string(b))

		log.Println(l)
	}
}
