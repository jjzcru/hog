package handler

import (
	"encoding/json"
	"fmt"
	"github.com/jjzcru/hog/pkg/utils"
	"net/http"
)

func serverError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	utils.PrintError(err)
	_, _ = fmt.Fprintf(w, `{"error":"%s"}`, jsonEscape(err.Error()))
}

func notFoundError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	utils.PrintError(err)
	_, _ = fmt.Fprintf(w, `{"error":"%s"}`, jsonEscape(err.Error()))
}

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	// Trim the beginning and trailing " character
	return string(b[1 : len(b)-1])
}
