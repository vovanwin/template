package ui

import (
	"encoding/json"
	"net/http"
)

func (c *UIController) handleCentrifugoToken(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userID, stop := c.requireAuth(w, r)
	if stop {
		return
	}

	token, err := c.centrifugoClient.GenerateToken(userID)
	if err != nil {
		c.log.Error("generate centrifugo token", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"url":   c.centrifugoURL,
	})
}
