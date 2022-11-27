package server

import (
	"net/http"
)

type UnwellHandler struct{}

func (*UnwellHandler) Pattern() string {
	return "/unwell"
}

func (*UnwellHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func NewUnwellHandler() *UnwellHandler {
	return &UnwellHandler{}
}
