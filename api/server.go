package api

import (
	"github.com/gorilla/mux"
	"ledgerproject/config"
	"ledgerproject/ledger"
	"net/http"
)

type Server struct {
	router *mux.Router
	ledger ledger.LedgerService
	config *config.Config
}

func NewServer(l ledger.LedgerService, c *config.Config) *Server {
	s := &Server{
		router: mux.NewRouter(),
		ledger: l,
		config: c,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/accounts", s.CreateAccountHandler).Methods("POST")
	s.router.HandleFunc("/transactions", s.RecordTransactionHandler).Methods("POST")
	s.router.HandleFunc("/accounts/{accountId}/balance", s.GetBalanceHandler).Methods("GET")
	s.router.HandleFunc("/accounts/{accountId}/history", s.GetTransactionHistoryHandler).Methods("GET")
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.config.ServerPort, s.router)
}
