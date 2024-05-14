package handlers

import (
	"fmt"
	"net/http"
)

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := r.URL.Query().Get("id")

		if clientID == "" {
			http.Error(w, "Client ID not provided", http.StatusBadRequest)
			return
		}

		websocketURL := fmt.Sprintf("ws://%s/ws/%s", r.Host, clientID)
		homeTemplate.Execute(w, websocketURL)
	}

}
