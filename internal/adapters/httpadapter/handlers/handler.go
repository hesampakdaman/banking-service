package handlers

import (
	"github.com/hesampakdaman/banking-service/internal/ports"
)

// httpHandler handles HTTP requests for banking operations.
type httpHandler struct {
	service ports.BankService
}

func NewHTTPHandler(service ports.BankService) *httpHandler {
	return &httpHandler{service: service}
}
