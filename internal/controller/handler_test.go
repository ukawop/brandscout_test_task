package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ukawop/brandscout_test_task/internal/models"
	"github.com/ukawop/brandscout_test_task/internal/repository"
	"github.com/ukawop/brandscout_test_task/internal/service"
)

func TestHandler_CreateQuote(t *testing.T) {
	repo := repository.NewQuoteRepository()
	srv := service.NewQuoteService(repo)
	h := NewHandler(srv, NewTestLogger())

	t.Run("Success", func(t *testing.T) {
		quoteReq := models.QuoteRequest{
			Author: "Test Author",
			Quote:  "Test Quote",
		}
		body, _ := json.Marshal(quoteReq)

		req := httptest.NewRequest("POST", "/quotes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		h.CreateQuote(rr, req)

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

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/quotes", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		h.CreateQuote(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})
}

func TestHandler_GetAllQuotes(t *testing.T) {
	repo := repository.NewQuoteRepository()
	srv := service.NewQuoteService(repo)
	h := NewHandler(srv, NewTestLogger())

	srv.CreateQuote(models.QuoteRequest{Author: "A1", Quote: "Q1"})
	srv.CreateQuote(models.QuoteRequest{Author: "A2", Quote: "Q2"})

	t.Run("Get All", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/quotes", nil)
		rr := httptest.NewRecorder()

		h.GetAllQuotes(rr, req)

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

	t.Run("Filter by Author", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/quotes?author=A1", nil)
		rr := httptest.NewRecorder()

		h.GetAllQuotes(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var response []models.Quote
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if len(response) != 1 {
			t.Fatalf("Expected 1 quote, got %d", len(response))
		}
		if response[0].Author != "A1" {
			t.Errorf("Expected author A1, got %s", response[0].Author)
		}
	})
}

func TestHandler_GetRandomQuote(t *testing.T) {
	repo := repository.NewQuoteRepository()
	srv := service.NewQuoteService(repo)
	h := NewHandler(srv, NewTestLogger())

	t.Run("No Quotes", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/quotes/random", nil)
		rr := httptest.NewRecorder()

		h.GetRandomQuote(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", rr.Code)
		}
	})

	t.Run("With Quotes", func(t *testing.T) {
		srv.CreateQuote(models.QuoteRequest{Author: "Random", Quote: "Quote"})

		req := httptest.NewRequest("GET", "/quotes/random", nil)
		rr := httptest.NewRecorder()

		h.GetRandomQuote(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var response models.Quote
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response.Author != "Random" {
			t.Errorf("Expected author Random, got %s", response.Author)
		}
	})
}

func TestHandler_DeleteQuote(t *testing.T) {
	repo := repository.NewQuoteRepository()
	srv := service.NewQuoteService(repo)
	h := NewHandler(srv, NewTestLogger())

	quote := srv.CreateQuote(models.QuoteRequest{Author: "ToDelete", Quote: "Me"})

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/quotes/"+strconv.Itoa(quote.ID), nil)
		rr := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(quote.ID)})

		h.DeleteQuote(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", rr.Code)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/quotes/invalid", nil)
		rr := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})

		h.DeleteQuote(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/quotes/999", nil)
		rr := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": "999"})

		h.DeleteQuote(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", rr.Code)
		}
	})
}
