package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ukawop/brandscout_test_task/internal/config"
	"github.com/ukawop/brandscout_test_task/internal/controller"
)

func main() {
	cfg := config.NewConfig()

	logger := log.New(os.Stdout, "APP: ", log.LstdFlags|log.Lshortfile)

	router := controller.NewRouter()

	logger.Printf("Сервер запущен на %s...\n", cfg.ServerAddress)
	logger.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
}
