package tests

import (
	"testing"

	"github.com/Nivedithabp/receipt-processor/models"
	"github.com/Nivedithabp/receipt-processor/utils"
)

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name     string
		receipt  models.Receipt
		expected int
	}{
		{
			name: "Valid Receipt - Target",
			receipt: models.Receipt{
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
			expected: 28,
		},
		{
			name: "Valid Receipt - M&M Corner Market",
			receipt: models.Receipt{
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
			expected: 109,
		},
		{
			name: "Empty Items List",
			receipt: models.Receipt{
				Retailer:     "Walmart",
				PurchaseDate: "2022-02-02",
				PurchaseTime: "10:45",
				Items:        []models.Item{},
				Total:        "10.00",
			},
			expected: 82, // 50 points for round dollar + 6 points for odd day + 1 for retailer name length
		},
		{
			name: "Single Item with Description Length Multiple of 3",
			receipt: models.Receipt{
				Retailer:     "Aldi",
				PurchaseDate: "2022-03-15",
				PurchaseTime: "15:01",
				Items: []models.Item{
					{ShortDescription: "Apple", Price: "5.00"},
				},
				Total: "5.00",
			},
			expected: 95, // 5 points for retailer + 10 points (2PM - 4PM) + 1 pair (5) + 2 points for description
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			points := utils.CalculatePoints(test.receipt)
			if points != test.expected {
				t.Errorf("Expected %d points but got %d", test.expected, points)
			}
		})
	}
}
