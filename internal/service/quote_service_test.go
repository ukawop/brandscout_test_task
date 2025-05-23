package service

import (
	"testing"

	"github.com/ukawop/brandscout_test_task/internal/models"
	"github.com/ukawop/brandscout_test_task/internal/repository"
)

func TestQuoteService_CreateQuote(t *testing.T) {
	repo := repository.NewQuoteRepository()
	srv := NewQuoteService(repo)

	quoteReq := models.QuoteRequest{
		Author: "Service Author",
		Quote:  "Service Quote",
	}

	created := srv.CreateQuote(quoteReq)

	if created.ID == 0 {
		t.Error("Expected non-zero ID")
	}
	if created.Author != quoteReq.Author {
		t.Errorf("Expected author %s, got %s", quoteReq.Author, created.Author)
	}
}

func TestQuoteService_GetAllQuotes(t *testing.T) {
	repo := repository.NewQuoteRepository()
	srv := NewQuoteService(repo)

	srv.CreateQuote(models.QuoteRequest{Author: "A1", Quote: "Q1"})
	srv.CreateQuote(models.QuoteRequest{Author: "A2", Quote: "Q2"})

	quotes := srv.GetAllQuotes()

	if len(quotes) != 2 {
		t.Fatalf("Expected 2 quotes, got %d", len(quotes))
	}
}

func TestQuoteService_GetQuotesByAuthor(t *testing.T) {
	repo := repository.NewQuoteRepository()
	srv := NewQuoteService(repo)

	srv.CreateQuote(models.QuoteRequest{Author: "Filter", Quote: "Q1"})
	srv.CreateQuote(models.QuoteRequest{Author: "Filter", Quote: "Q2"})
	srv.CreateQuote(models.QuoteRequest{Author: "Other", Quote: "Q3"})

	quotes := srv.GetQuotesByAuthor("Filter")

	if len(quotes) != 2 {
		t.Fatalf("Expected 2 quotes, got %d", len(quotes))
	}
	for _, q := range quotes {
		if q.Author != "Filter" {
			t.Errorf("Expected author Filter, got %s", q.Author)
		}
	}
}

func TestQuoteService_GetRandomQuote(t *testing.T) {
	repo := repository.NewQuoteRepository()
	srv := NewQuoteService(repo)

	_, err := srv.GetRandomQuote()
	if err == nil {
		t.Error("Expected error for empty repository")
	}

	srv.CreateQuote(models.QuoteRequest{Author: "R1", Quote: "Q1"})
	quote, err := srv.GetRandomQuote()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if quote.Author != "R1" {
		t.Errorf("Expected author R1, got %s", quote.Author)
	}
}

func TestQuoteService_DeleteQuote(t *testing.T) {
	repo := repository.NewQuoteRepository()
	srv := NewQuoteService(repo)

	quote := srv.CreateQuote(models.QuoteRequest{Author: "ToDelete", Quote: "DeleteMe"})

	err := srv.DeleteQuote(quote.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	err = srv.DeleteQuote(999)
	if err == nil {
		t.Error("Expected error for non-existent quote")
	}
}
