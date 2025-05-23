package service

import (
	"github.com/ukawop/brandscout_test_task/internal/models"
	"github.com/ukawop/brandscout_test_task/internal/repository"
)

type QuoteService struct {
	repo *repository.QuoteRepository
}

func NewQuoteService(repo *repository.QuoteRepository) *QuoteService {
	return &QuoteService{repo: repo}
}

func (s *QuoteService) CreateQuote(quote models.QuoteRequest) models.Quote {
	return s.repo.Create(quote)
}

func (s *QuoteService) GetAllQuotes() []models.Quote {
	return s.repo.GetAll()
}

func (s *QuoteService) GetQuotesByAuthor(author string) []models.Quote {
	return s.repo.GetByAuthor(author)
}

func (s *QuoteService) GetRandomQuote() (models.Quote, error) {
	return s.repo.GetRandom()
}

func (s *QuoteService) DeleteQuote(id int) error {
	return s.repo.Delete(id)
}
