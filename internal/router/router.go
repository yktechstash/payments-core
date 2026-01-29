package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) Router() http.Handler {

	r := mux.NewRouter()

	r.HandleFunc("/accounts", s.handleCreateAccount).Methods("POST")
	r.HandleFunc("/accounts/{account_id}", s.handleGetAccount).Methods("GET")
	r.HandleFunc("/transactions", s.handleCreateTransaction).Methods("POST")

	// basic health
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods("GET")

	return r
}
