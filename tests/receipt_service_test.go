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
			name: "Example 1: Basic Receipt",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
					{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
				},
				Total: "35.35",
			},
			expected: 28,
		},
		{
			name: "Example 2: Round Dollar Total",
			receipt: models.Receipt{
				Retailer:     "Walmart",
				PurchaseDate: "2022-02-02",
				PurchaseTime: "14:30",
				Items: []models.Item{
					{ShortDescription: "Pepsi 12oz", Price: "2.25"},
					{ShortDescription: "Dasani Water", Price: "1.50"},
				},
				Total: "10.00",
			},
			expected: 93, // Includes 50 points for round dollar, 25 for 0.25, and other rules
		},
		{
			name: "Example 3: Empty Items",
			receipt: models.Receipt{
				Retailer:     "Costco",
				PurchaseDate: "2022-03-15",
				PurchaseTime: "10:45",
				Items:        []models.Item{},
				Total:        "1.00",
			},
			expected: 7, // 6 points for date + 1 for retailer name
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

