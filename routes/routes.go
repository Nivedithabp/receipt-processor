package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Nivedithabp/receipt-processor/models"
	"github.com/Nivedithabp/receipt-processor/services"
)

// RegisterRoutes registers API routes
func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/receipts/process", ProcessReceiptHandler).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", GetPointsHandler).Methods("GET")
}

// ProcessReceiptHandler handles receipt submission
func ProcessReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id := services.ProcessReceipt(receipt)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// GetPointsHandler retrieves points for a receipt
func GetPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	points, exists := services.GetPoints(id)
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"points": points})
}
