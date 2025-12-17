package server

import (
	"errors"
)

var (
	ErrEmptyUUID        = errors.New("empty uuid")
	ErrEmptyPassword    = errors.New("empty password")
	ErrPasswordTooShort = errors.New("password too short")
	ErrEmptyContent     = errors.New("empty content")
	ErrInvalidBinID     = errors.New("invalid bin id")
)

func (s *Server) routes() {
	s.router.
		HandleFunc("/liveness", s.logRequest(s.handleLiveness())).
		Methods("GET")
	s.router.
		HandleFunc("/readiness", s.logRequest(s.handleReadiness())).
		Methods("GET")
	s.router.
		HandleFunc("/api/v1/{bid}", s.logRequest(s.handleAPIGetBin())).
		Methods("GET")
	s.router.
		HandleFunc("/create", s.logRequest(s.handleBinCreate())).
		Methods("POST")
	s.router.
		HandleFunc("/{bid}/auth", s.logRequest(s.handleBinAuth())).
		Methods("POST")
	s.router.
		HandleFunc("/{bid}/update", s.logRequest(s.handleBinUpdate())).
		Methods("POST")
	s.router.
		HandleFunc("/{bid}", s.logRequest(s.handleBinShow())).
		Methods("GET")
	s.router.
		HandleFunc("/", s.logRequest(s.handleRoot())).
		Methods("GET")
	s.router.NotFoundHandler = s.logRequest(s.handleNotFound())
}
