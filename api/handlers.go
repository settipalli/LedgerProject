package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"ledgerproject/logger"
	"ledgerproject/models"
	"net/http"
)

func (s *Server) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		// X-Forwarded-For header
		clientIP := r.Header.Get("X-Forwarded-For")
		// If not present, fallback to RemoteAddr
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}

		log.Error("Failed to decode account creation request",
			zap.Error(err),
			zap.String("remote_addr", clientIP))

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.ledger.CreateAccount(account); err != nil {
		log.Error("Failed to create account",
			zap.Error(err),
			zap.String("account_id", account.ID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("Account created successfully", zap.String("account_id", account.ID))
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) RecordTransactionHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	var tx models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		clientIP := r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}

		log.Error("Failed to decode record transaction request",
			zap.Error(err),
			zap.String("remote_addr", clientIP))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.ledger.RecordTransaction(tx); err != nil {
		log.Error("Failed to record transaction",
			zap.Error(err),
			zap.String("transaction_id", tx.ID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("Transaction recorded successfully", zap.String("transaction_id", tx.ID))
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	balance, err := s.ledger.GetAccountBalance(accountID)
	if err != nil {
		log.Error("Failed to get account balance",
			zap.Error(err),
			zap.String("account_id", accountID),
		)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Info("Account balance retrieved successfully", zap.String("account_id", accountID))
	json.NewEncoder(w).Encode(map[string]models.Money{"balance": balance})
}

func (s *Server) GetTransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	history := s.ledger.GetTransactionHistory(accountID)
	log.Info("Successfully generated transaction history", zap.String("account_id", accountID))
	json.NewEncoder(w).Encode(history)
}
