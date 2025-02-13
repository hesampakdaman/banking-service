package handlers

import (
	"errors"
	"net/http"

	"github.com/hesampakdaman/banking-service/internal/domain"
)

func domainErrToStatusCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidOwner),
		errors.Is(err, domain.ErrInvalidTransactionType),
		errors.Is(err, domain.ErrNegativeBalance),
		errors.Is(err, domain.ErrSelfTransfer),
		errors.Is(err, domain.ErrInvalidAmount):
		return http.StatusBadRequest

	case errors.Is(err, domain.ErrAccountAlreadyExists),
		errors.Is(err, domain.ErrInsufficientFunds):
		return http.StatusConflict

	case errors.Is(err, domain.ErrInvalidAccountID):
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}
