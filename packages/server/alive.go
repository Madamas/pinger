package server

import (
	"net/http"
)

type AliveHandler struct{}

func (*AliveHandler) Pattern() string {
	return "/alive"
}

func (*AliveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func NewAliveHandler() *AliveHandler {
	return &AliveHandler{}
}
