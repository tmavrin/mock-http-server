package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type config struct {
	PrepareErrorPath string    `json:"prepareErrorPath"`
	Handlers         []handler `json:"handlers"`
}

func main() {
	dat, err := os.ReadFile(cmp.Or(os.Getenv("CONFIG_PATH"), "./config.json"))
	if err != nil {
		fmt.Print(err)
	}

	var cfg config
	err = json.Unmarshal(dat, &cfg)
	if err != nil {
		fmt.Print(err)
	}

	router := parseRoutes(cfg)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", cmp.Or(os.Getenv("PORT"), "8080")),
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("running server: %s", err)
		}
	}()

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM)
	<-quitChan
	log.Println("received interrupt, shutting down...")
}
