package router

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/payments-core/internal/domain"
	"net/http"
	"strconv"
)

type createAccountRequest struct {
	AccountID      int64  `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}

func (s *Server) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	m, err := domain.ParseMoney(req.InitialBalance)
	if err != nil {
		writeDomainErr(w, err)
		return
	}

	err = s.svc.CreateAccount(r.Context(), domain.AccountID(req.AccountID), m)
	if err != nil {
		writeDomainErr(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type getAccountResponse struct {
	AccountID int64  `json:"account_id"`
	Balance   string `json:"balance"`
}

func (s *Server) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["account_id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid account_id")
		return
	}

	acc, err := s.svc.GetAccount(r.Context(), domain.AccountID(id))
	if err != nil {
		writeDomainErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, getAccountResponse{
		AccountID: int64(acc.AccountID),
		Balance:   acc.Balance.String(),
	})
}
