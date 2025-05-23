package repository

import (
	"testing"

	"github.com/ukawop/brandscout_test_task/internal/models"
)

func TestQuoteRepository_Create(t *testing.T) {
	repo := NewQuoteRepository()
	quoteReq := models.QuoteRequest{
		Author: "Test Author",
		Quote:  "Test Quote",
	}

	created := repo.Create(quoteReq)

	if created.ID != 1 {
		t.Errorf("Expected ID 1, got %d", created.ID)
	}
	if created.Author != quoteReq.Author {
		t.Errorf("Expected author %s, got %s", quoteReq.Author, created.Author)
	}
	if created.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestQuoteRepository_GetAll(t *testing.T) {
	repo := NewQuoteRepository()
	repo.Create(models.QuoteRequest{Author: "A1", Quote: "Q1"})
	repo.Create(models.QuoteRequest{Author: "A2", Quote: "Q2"})

	quotes := repo.GetAll()

	if len(quotes) != 2 {
		t.Fatalf("Expected 2 quotes, got %d", len(quotes))
	}
	if quotes[0].Author != "A1" {
		t.Errorf("Expected first author A1, got %s", quotes[0].Author)
	}
}

func TestQuoteRepository_GetByAuthor(t *testing.T) {
	repo := NewQuoteRepository()
	repo.Create(models.QuoteRequest{Author: "Author1", Quote: "Q1"})
	repo.Create(models.QuoteRequest{Author: "Author2", Quote: "Q2"})
	repo.Create(models.QuoteRequest{Author: "Author1", Quote: "Q3"})

	quotes := repo.GetByAuthor("Author1")

	if len(quotes) != 2 {
		t.Fatalf("Expected 2 quotes, got %d", len(quotes))
	}
	for _, q := range quotes {
		if q.Author != "Author1" {
			t.Errorf("Expected author Author1, got %s", q.Author)
		}
	}
}

func TestQuoteRepository_GetRandom(t *testing.T) {
	repo := NewQuoteRepository()

	_, err := repo.GetRandom()
	if err != ErrQuoteNotFound {
		t.Errorf("Expected ErrQuoteNotFound, got %v", err)
	}

	repo.Create(models.QuoteRequest{Author: "R1", Quote: "Q1"})
	repo.Create(models.QuoteRequest{Author: "R2", Quote: "Q2"})

	seen := make(map[int]bool)
	for i := 0; i < 10; i++ {
		quote, err := repo.GetRandom()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		seen[quote.ID] = true
	}

	if len(seen) < 1 {
		t.Error("Expected at least one unique quote")
	}
}

func TestQuoteRepository_Delete(t *testing.T) {
	repo := NewQuoteRepository()
	quote := repo.Create(models.QuoteRequest{Author: "ToDelete", Quote: "DeleteMe"})

	err := repo.Delete(quote.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	quotes := repo.GetAll()
	for _, q := range quotes {
		if q.ID == quote.ID {
			t.Error("Quote was not deleted")
		}
	}

	err = repo.Delete(999)
	if err != ErrQuoteNotFound {
		t.Errorf("Expected ErrQuoteNotFound, got %v", err)
	}
}
