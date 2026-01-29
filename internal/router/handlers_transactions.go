package router

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/payments-core/internal/domain"
)

type createTransactionRequest struct {
	SourceAccountID      int64  `json:"source_account_id"`
	DestinationAccountID int64  `json:"destination_account_id"`
	Amount               string `json:"amount"`
}

func (s *Server) handleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	amt, err := domain.ParseMoney(req.Amount)
	if err != nil {
		writeDomainErr(w, err)
		return
	}

	txID := uuid.NewString()
	err = s.svc.Transfer(
		r.Context(),
		domain.AccountID(req.SourceAccountID),
		domain.AccountID(req.DestinationAccountID),
		amt,
		txID,
	)
	if err != nil {
		writeDomainErr(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
