package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"ledgerproject/models"
	"net/http"
)

func (s *Server) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.ledger.CreateAccount(account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) RecordTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var tx models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.ledger.RecordTransaction(tx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	balance, err := s.ledger.GetAccountBalance(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]models.Money{"balance": balance})
}

func (s *Server) GetTransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	history := s.ledger.GetTransactionHistory(accountID)
	json.NewEncoder(w).Encode(history)
}
