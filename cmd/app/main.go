package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ukawop/brandscout_test_task/internal/config"
	"github.com/ukawop/brandscout_test_task/internal/controller"
	"github.com/ukawop/brandscout_test_task/internal/repository"
	"github.com/ukawop/brandscout_test_task/internal/service"
)

func main() {
	cfg := config.NewConfig(10)
	logger := log.New(os.Stdout, "APP: ", log.LstdFlags|log.Lshortfile)
	httpLogger := log.New(os.Stdout, "HTTP: ", log.LstdFlags|log.Lshortfile)
	repo := repository.NewQuoteRepository()
	service := service.NewQuoteService(repo)

	router := controller.NewRouter(httpLogger, repo, service, cfg.MaxHandlerWorkers)

	logger.Printf("Starting async server on %s...", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
}
