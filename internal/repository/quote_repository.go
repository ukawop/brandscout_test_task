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
	quotes      map[int]models.Quote
	authorIndex map[string]map[int]bool
	allIDs      map[int]bool
	mu          sync.RWMutex
	nextID      int
	rand        *rand.Rand
}

func NewQuoteRepository() *QuoteRepository {
	return &QuoteRepository{
		quotes:      make(map[int]models.Quote),
		authorIndex: make(map[string]map[int]bool),
		allIDs:      make(map[int]bool),
		nextID:      1,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
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

	r.quotes[r.nextID] = newQuote

	if r.authorIndex[quote.Author] == nil {
		r.authorIndex[quote.Author] = make(map[int]bool)
	}
	r.authorIndex[quote.Author][r.nextID] = true

	r.allIDs[r.nextID] = true
	r.nextID++

	return newQuote
}

func (r *QuoteRepository) GetAll() []models.Quote {
	r.mu.RLock()
	defer r.mu.RUnlock()

	quotes := make([]models.Quote, 0, len(r.quotes))
	for id := range r.allIDs {
		quotes = append(quotes, r.quotes[id])
	}
	return quotes
}

func (r *QuoteRepository) GetByAuthor(author string) []models.Quote {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids, exists := r.authorIndex[author]
	if !exists {
		return nil
	}

	quotes := make([]models.Quote, 0, len(ids))
	for id := range ids {
		quotes = append(quotes, r.quotes[id])
	}
	return quotes
}

func (r *QuoteRepository) GetRandom() (models.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.allIDs) == 0 {
		return models.Quote{}, ErrQuoteNotFound
	}

	randomIdx := r.rand.Intn(len(r.allIDs))
	i := 0
	for id := range r.allIDs {
		if i == randomIdx {
			return r.quotes[id], nil
		}
		i++
	}
	return models.Quote{}, ErrQuoteNotFound
}

func (r *QuoteRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	quote, exists := r.quotes[id]
	if !exists {
		return ErrQuoteNotFound
	}

	delete(r.quotes, id)
	delete(r.authorIndex[quote.Author], id)
	delete(r.allIDs, id)

	if len(r.authorIndex[quote.Author]) == 0 {
		delete(r.authorIndex, quote.Author)
	}

	return nil
}
