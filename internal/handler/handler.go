package handler

import (
	"github.com/gorilla/mux"
)

type Handler interface {
	Register(r *mux.Router)
}
