package repository

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/ukawop/brandscout_test_task/internal/models"
)

var (
	ErrQuoteNotFound = errors.New("цитата не найдена")
)

type QuoteRepository struct {
	quotes []models.Quote
	mu     sync.RWMutex
	nextID int
}

func NewQuoteRepository() *QuoteRepository {
	return &QuoteRepository{
		quotes: make([]models.Quote, 0),
		nextID: 1,
	}
}

func (r *QuoteRepository) Create(quote models.QuoteRequest) models.Quote {
	r.mu.Lock()
	defer r.mu.Unlock()

	newQuote := models.Quote{
		ID:        r.nextID,
		Author:    quote.Author,
		Quote:     quote.Quote,
		CreatedAt: time.Now(),
	}

	r.quotes = append(r.quotes, newQuote)
	r.nextID++

	return newQuote
}

func (r *QuoteRepository) GetAll() []models.Quote {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.quotes
}

func (r *QuoteRepository) GetByAuthor(author string) []models.Quote {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []models.Quote
	for _, q := range r.quotes {
		if q.Author == author {
			filtered = append(filtered, q)
		}
	}

	return filtered
}

func (r *QuoteRepository) GetRandom() (models.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.quotes) == 0 {
		return models.Quote{}, ErrQuoteNotFound
	}

	rand.Seed(time.Now().UnixNano())
	return r.quotes[rand.Intn(len(r.quotes))], nil
}

func (r *QuoteRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, q := range r.quotes {
		if q.ID == id {
			r.quotes = append(r.quotes[:i], r.quotes[i+1:]...)
			return nil
		}
	}

	return ErrQuoteNotFound
}
