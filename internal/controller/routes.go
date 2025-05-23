package controller

import (
	"log"

	"github.com/gorilla/mux"

	"github.com/ukawop/brandscout_test_task/internal/repository"
	"github.com/ukawop/brandscout_test_task/internal/service"
)

func NewRouter(logger *log.Logger, repo *repository.QuoteRepository, service *service.QuoteService, maxWorkers int) *mux.Router {
	handler := NewAsyncHandler(service, logger, maxWorkers)

	router := mux.NewRouter()

	router.HandleFunc("/quotes", handler.wrapHandler("create")).Methods("POST")
	router.HandleFunc("/quotes", handler.wrapHandler("getAll")).Methods("GET")
	router.HandleFunc("/quotes/random", handler.wrapHandler("getRandom")).Methods("GET")
	router.HandleFunc("/quotes/{id}", handler.wrapHandler("delete")).Methods("DELETE")

	return router
}
