package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type CustomError map[string]any

func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func parseIdParam(id string) (int, error) {
	parsedId, err := strconv.Atoi(id)
	if err != nil || parsedId < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return parsedId, nil
}
