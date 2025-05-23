package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/ukawop/brandscout_test_task/internal/models"
	"github.com/ukawop/brandscout_test_task/internal/service"
)

type Handler struct {
	service *service.QuoteService
	logger  *log.Logger
}

func NewHandler(service *service.QuoteService, logger *log.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req models.QuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Ошибка декодирования запроса: %v", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	h.logger.Printf("Создание цитаты: автор '%s', текст '%s'", req.Author, req.Quote)
	quote := h.service.CreateQuote(req)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quote); err != nil {
		h.logger.Printf("Ошибка кодирования ответа: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	h.logger.Printf("Цитата создана (ID: %d). Время выполнения: %v", quote.ID, time.Since(start))
}

func (h *Handler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	author := r.URL.Query().Get("author")

	var quotes []models.Quote
	if author != "" {
		h.logger.Printf("Запрос цитат автора: '%s'", author)
		quotes = h.service.GetQuotesByAuthor(author)
	} else {
		h.logger.Println("Запрос всех цитат")
		quotes = h.service.GetAllQuotes()
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quotes); err != nil {
		h.logger.Printf("Ошибка кодирования ответа: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	h.logger.Printf("Возвращено %d цитат. Время выполнения: %v", len(quotes), time.Since(start))
}

func (h *Handler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	h.logger.Println("Запрос случайной цитаты")

	quote, err := h.service.GetRandomQuote()
	if err != nil {
		h.logger.Println("Случайная цитата не найдена")
		http.Error(w, "Нет доступных цитат", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quote); err != nil {
		h.logger.Printf("Ошибка кодирования ответа: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	h.logger.Printf("Возвращена случайная цитата (ID: %d). Время выполнения: %v", quote.ID, time.Since(start))
}

func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	vars := mux.Vars(r)
	id := vars["id"]

	h.logger.Printf("Попытка удаления цитаты с ID: %s", id)

	var intID int
	_, err := fmt.Sscanf(id, "%d", &intID)
	if err != nil {
		h.logger.Printf("Неверный формат ID: %s", id)
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteQuote(intID); err != nil {
		h.logger.Printf("Цитата с ID %d не найдена", intID)
		http.Error(w, "Цитата не найдена", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.logger.Printf("Цитата с ID %d удалена. Время выполнения: %v", intID, time.Since(start))
}
