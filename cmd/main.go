package main

import (
	"cmp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tmavrin/mock-http-server/internal/config"
	"github.com/tmavrin/mock-http-server/internal/mock"
)

func main() {
	cfg := configPrepare()

	serverPort := cmp.Or(os.Getenv("PORT"), "8080")

	server := http.Server{
		Addr:              ":" + serverPort,
		Handler:           mock.Router(cfg),
		ReadHeaderTimeout: 30 * time.Second,
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

func configPrepare() *config.Config {
	configPath := cmp.Or(os.Getenv("CONFIG_PATH"), config.Path)
	configTypeEnv := cmp.Or(os.Getenv("CONFIG_TYPE"), "json")

	var configType config.Type

	configType.Parse(configTypeEnv)

	dat, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("reading config: %s", err)
	}

	cfg, err := config.Parse(configType, dat)
	if err != nil {
		log.Fatalf("parsing config: %s", err)
	}

	return cfg
}
