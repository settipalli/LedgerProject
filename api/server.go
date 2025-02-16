package api

import (
	"context"
	"github.com/gorilla/mux"
	"ledgerproject/config"
	"ledgerproject/ledger"
	"net/http"
)

type Server struct {
	router *mux.Router
	ledger ledger.LedgerService
	config *config.Config
	server *http.Server
}

func NewServer(l ledger.LedgerService, c *config.Config) *Server {
	r := mux.NewRouter()
	s := &Server{
		router: r,
		ledger: l,
		config: c,
		server: &http.Server{
			Addr:              c.ServerPort,
			Handler:           r,
			ReadTimeout:       c.ReadTimeout,
			WriteTimeout:      c.WriteTimeout,
			IdleTimeout:       c.IdleTimeout,
			ReadHeaderTimeout: c.ReadHeaderTimeout,
			MaxHeaderBytes:    c.MaxHeaderBytes,
		},
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
	return s.server.ListenAndServe()
}

// Graceful shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
