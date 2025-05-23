package controller

import (
	"io"
	"log"
)

// NewTestLogger создает логгер для тестов, который ничего не выводит
func NewTestLogger() *log.Logger {
	return log.New(io.Discard, "TEST: ", log.LstdFlags)
}
