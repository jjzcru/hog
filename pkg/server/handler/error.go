package handler

import (
	"fmt"
	"net/http"
)

func notImplemented(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
}

func serverError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
}

func notFoundError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
}
