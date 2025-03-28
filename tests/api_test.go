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

func TestProcessReceiptAPI(t *testing.T) {
	tests := []struct {
		name           string
		payload        models.Receipt
		expectedStatus int
	}{
		{
			name: "Valid Receipt - Target",
			payload: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
					{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
					{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
					{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
					{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
				},
				Total: "35.35",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Valid Receipt - M&M Corner Market",
			payload: models.Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Items: []models.Item{
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
				},
				Total: "9.00",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid Receipt - Missing Retailer",
			payload: models.Receipt{
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
				},
				Total: "35.35",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid Receipt - Empty Items",
			payload: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{},
				Total:        "35.35",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(test.payload)
			req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(payloadBytes))
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if status := recorder.Code; status != test.expectedStatus {
				t.Errorf("Expected status code %d but got %d", test.expectedStatus, status)
			}

			if test.expectedStatus == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil || response["id"] == "" {
					t.Errorf("Expected valid receipt ID but got error or empty string")
				}
			}
		})
	}
}

func TestGetPointsAPI(t *testing.T) {
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	// Process a valid receipt to generate an ID
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "35.35",
	}

	payloadBytes, _ := json.Marshal(receipt)
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(payloadBytes))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	var response map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &response)
	receiptID := response["id"]

	tests := []struct {
		name           string
		receiptID      string
		expectedStatus int
		expectedPoints int
	}{
		{
			name:           "Valid Receipt ID - Target",
			receiptID:      receiptID,
			expectedStatus: http.StatusOK,
			expectedPoints: 28, // Corrected points for Target receipt
		},
		{
			name:           "Invalid Receipt ID",
			receiptID:      "invalid-id",
			expectedStatus: http.StatusNotFound,
			expectedPoints: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/receipts/"+test.receiptID+"/points", nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if status := recorder.Code; status != test.expectedStatus {
				t.Errorf("Expected status code %d but got %d", test.expectedStatus, status)
			}

			if test.expectedStatus == http.StatusOK {
				var pointsResponse map[string]int
				json.Unmarshal(recorder.Body.Bytes(), &pointsResponse)
				if pointsResponse["points"] != test.expectedPoints {
					t.Errorf("Expected %d points but got %d", test.expectedPoints, pointsResponse["points"])
				}
			}
		})
	}
}
