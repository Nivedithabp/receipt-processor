package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Nivedithabp/receipt-processor/models"
	"github.com/Nivedithabp/receipt-processor/services"
)

// RegisterRoutes registers all routes for the API
func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/receipts/process", ProcessReceiptHandler).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", GetPointsHandler).Methods("GET")
}

// @Summary Process a receipt and generate an ID
// @Description Submits a receipt and returns a unique ID
// @Tags Receipts
// @Accept json
// @Produce json
// @Param receipt body models.Receipt true "Receipt JSON"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /receipts/process [post]
func ProcessReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil || !isValidReceipt(receipt) {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id := services.ProcessReceipt(receipt)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// @Summary Get points for a receipt
// @Description Returns points for a given receipt ID
// @Tags Receipts
// @Produce json
// @Param id path string true "Receipt ID"
// @Success 200 {object} map[string]int
// @Failure 404 {object} map[string]string
// @Router /receipts/{id}/points [get]
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

// isValidReceipt validates required fields
func isValidReceipt(receipt models.Receipt) bool {
	if receipt.Retailer == "" ||
		receipt.PurchaseDate == "" ||
		receipt.PurchaseTime == "" ||
		len(receipt.Items) == 0 ||
		receipt.Total == "" {
		return false
	}
	return true
}
