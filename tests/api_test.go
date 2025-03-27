package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/Nivedithabp/receipt-processor/models"
	"github.com/Nivedithabp/receipt-processor/routes"
)

func TestProcessReceipt(t *testing.T) {
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
		},
		Total: "35.35",
	}

	payload, _ := json.Marshal(receipt)
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status code 200 but got %v", status)
	}

	var response map[string]string
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}

	if _, exists := response["id"]; !exists {
		t.Errorf("Expected 'id' in the response but did not find it")
	}
}

func TestGetPoints(t *testing.T) {
	// First, process a receipt to get an ID
	receipt := models.Receipt{
		Retailer:     "Walmart",
		PurchaseDate: "2022-02-02",
		PurchaseTime: "14:30",
		Items: []models.Item{
			{ShortDescription: "Pepsi 12oz", Price: "2.25"},
			{ShortDescription: "Dasani Water", Price: "1.50"},
		},
		Total: "10.00",
	}

	payload, _ := json.Marshal(receipt)
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	router.ServeHTTP(recorder, req)

	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	receiptID, exists := response["id"]
	if !exists {
		t.Fatalf("No receipt ID returned after processing")
	}

	// Now, get points using the generated ID
	pointsURL := "/receipts/" + receiptID + "/points"
	req, err = http.NewRequest("GET", pointsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status code 200 but got %v", status)
	}

	var pointsResponse map[string]int
	err = json.Unmarshal(recorder.Body.Bytes(), &pointsResponse)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}

	if _, exists := pointsResponse["points"]; !exists {
		t.Errorf("Expected 'points' in the response but did not find it")
	}
}
