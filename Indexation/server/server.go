package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const user = "admin"
const password = "Complexpass#123"
const searchPath = "http://localhost:4080/api/emails/_search"

type MyServer struct {
	server *http.Server
}

type QuerySearch struct {
	Term      string `json:"term"`
	From      string `json:"from"`
	MaxResult string `json:"max_results"`
}

type QueryAll struct {
	From      string `json:"from"`
	MaxResult string `json:"max_results"`
}

func NewServer(mux *chi.Mux) *MyServer {
	s := &http.Server{
		Addr:           ":9000",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &MyServer{server: s}
}

func (s *MyServer) Run() {
	s.server.ListenAndServe()
}
