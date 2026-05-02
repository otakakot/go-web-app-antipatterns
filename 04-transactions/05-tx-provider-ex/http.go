package main

import (
	"encoding/json"
	"net/http"
)

func NewHTTPHandler(
	usePointsAsDiscountHandler UsePointsAsDiscountHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /use-points", func(w http.ResponseWriter, r *http.Request) {
		type payload struct {
			UserID int `json:"user_id"`
			Points int `json:"points"`
		}

		var p payload
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = usePointsAsDiscountHandler.Handle(r.Context(), UsePointsAsDiscount(p))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return mux
}
