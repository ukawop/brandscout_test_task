package controller

import (
	"log"
	"os"

	"github.com/gorilla/mux"

	"github.com/ukawop/brandscout_test_task/internal/repository"
	"github.com/ukawop/brandscout_test_task/internal/service"
)

func NewRouter() *mux.Router {
	logger := log.New(os.Stdout, "HTTP: ", log.LstdFlags|log.Lshortfile)

	repo := repository.NewQuoteRepository()
	service := service.NewQuoteService(repo)
	handler := NewHandler(service, logger)

	router := mux.NewRouter()

	router.HandleFunc("/quotes", handler.CreateQuote).Methods("POST")
	router.HandleFunc("/quotes", handler.GetAllQuotes).Methods("GET")
	router.HandleFunc("/quotes/random", handler.GetRandomQuote).Methods("GET")
	router.HandleFunc("/quotes/{id}", handler.DeleteQuote).Methods("DELETE")

	return router
}
