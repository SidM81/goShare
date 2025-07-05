package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SidM81/goShare/config"
	"github.com/SidM81/goShare/models"
)

func ShareHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ShareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	var file models.File
	if err := db.Where("name = ?", req.Name).First(&file).Error; err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Create shareable URL
	shareURL := fmt.Sprintf("http://localhost:8080/download?name=%s", req.Name)

	resp := models.ShareResponse{ShareURL: shareURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
