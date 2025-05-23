package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/ukawop/brandscout_test_task/internal/models"
	"github.com/ukawop/brandscout_test_task/internal/repository"
	"github.com/ukawop/brandscout_test_task/internal/service"
)

func TestAsyncHandler(t *testing.T) {
	repo := repository.NewQuoteRepository()
	service := service.NewQuoteService(repo)
	router := NewRouter(NewTestLogger(), repo, service, 10)

	t.Run("Create Quote Async", func(t *testing.T) {
		quoteReq := models.QuoteRequest{
			Author: "Test Author",
			Quote:  "Test Quote",
		}
		body, _ := json.Marshal(quoteReq)

		req := httptest.NewRequest("POST", "/quotes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var response models.Quote
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response.Author != quoteReq.Author {
			t.Errorf("Expected author %s, got %s", quoteReq.Author, response.Author)
		}
	})
}
func TestAsyncHandlerGetAll(t *testing.T) {
	repo := repository.NewQuoteRepository()
	service := service.NewQuoteService(repo)
	router := NewRouter(NewTestLogger(), repo, service, 10)
	t.Run("Get All Quotes Async", func(t *testing.T) {
		service.CreateQuote(models.QuoteRequest{Author: "Author1", Quote: "Quote1"})
		service.CreateQuote(models.QuoteRequest{Author: "Author2", Quote: "Quote2"})

		req := httptest.NewRequest("GET", "/quotes", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var response []models.Quote
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if len(response) != 2 {
			t.Errorf("Expected 2 quotes, got %d", len(response))
		}
	})
}
func TestAsyncHandlerGetByFilter(t *testing.T) {
	repo := repository.NewQuoteRepository()
	service := service.NewQuoteService(repo)
	router := NewRouter(NewTestLogger(), repo, service, 10)
	t.Run("Get Quotes Filtered By Author Async", func(t *testing.T) {
		service.CreateQuote(models.QuoteRequest{Author: "FilterAuthor", Quote: "Q1"})
		service.CreateQuote(models.QuoteRequest{Author: "FilterAuthor", Quote: "Q2"})
		service.CreateQuote(models.QuoteRequest{Author: "OtherAuthor", Quote: "Q3"})

		req := httptest.NewRequest("GET", "/quotes?author=FilterAuthor", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var response []models.Quote
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if len(response) != 2 {
			t.Errorf("Expected 2 quotes, got %d", len(response))
		}
		for _, q := range response {
			if q.Author != "FilterAuthor" {
				t.Errorf("Expected author FilterAuthor, got %s", q.Author)
			}
		}
	})
}

func TestAsyncHandlerGetRandom(t *testing.T) {
	repo := repository.NewQuoteRepository()
	service := service.NewQuoteService(repo)
	router := NewRouter(NewTestLogger(), repo, service, 10)
	t.Run("Get Random Quote Async", func(t *testing.T) {
		service.CreateQuote(models.QuoteRequest{Author: "RandomAuthor", Quote: "RandomQuote"})

		req := httptest.NewRequest("GET", "/quotes/random", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var response models.Quote
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response.Author != "RandomAuthor" {
			t.Errorf("Expected author RandomAuthor, got %s", response.Author)
		}
	})
}
func TestAsyncHandlerDelete(t *testing.T) {
	repo := repository.NewQuoteRepository()
	service := service.NewQuoteService(repo)
	router := NewRouter(NewTestLogger(), repo, service, 10)
	t.Run("Delete Quote Async", func(t *testing.T) {
		quote := service.CreateQuote(models.QuoteRequest{Author: "ToDelete", Quote: "DeleteMe"})

		req := httptest.NewRequest("DELETE", "/quotes/"+strconv.Itoa(quote.ID), nil)
		rr := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(quote.ID)})

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", rr.Code)
		}

		_, err := service.GetRandomQuote()
		if err == nil {
			t.Error("Quote was not deleted")
		}
	})

	t.Run("Delete Non-Existent Quote Async", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/quotes/999", nil)
		rr := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "999"})

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", rr.Code)
		}
	})
}
func TestAsyncHandlersStress(t *testing.T) {
	repo := repository.NewQuoteRepository()
	service := service.NewQuoteService(repo)
	router := NewRouter(NewTestLogger(), repo, service, 10)

	t.Run("Concurrent Requests Stress Test", func(t *testing.T) {
		start := time.Now()
		results := make(chan bool, 20)
		for i := 0; i < 20; i++ {
			go func(i int) {
				if i%2 == 0 {
					quoteReq := models.QuoteRequest{
						Author: "Concurrent Author " + strconv.Itoa(i),
						Quote:  "Concurrent Quote " + strconv.Itoa(i),
					}
					body, _ := json.Marshal(quoteReq)

					req := httptest.NewRequest("POST", "/quotes", bytes.NewReader(body))
					req.Header.Set("Content-Type", "application/json")
					rr := httptest.NewRecorder()

					router.ServeHTTP(rr, req)
					results <- rr.Code == http.StatusOK
				} else {
					req := httptest.NewRequest("GET", "/quotes", nil)
					rr := httptest.NewRecorder()

					router.ServeHTTP(rr, req)
					results <- rr.Code == http.StatusOK
				}
			}(i)
		}

		successCount := 0
		for i := 0; i < 20; i++ {
			if <-results {
				successCount++
			}
		}

		if successCount != 20 {
			t.Errorf("Expected 20 successful requests, got %d", successCount)
		}

		t.Logf("Processed 20 concurrent requests in %v", time.Since(start))
	})
}
