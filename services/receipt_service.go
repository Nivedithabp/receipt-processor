package services

import (
	"sync"

	"github.com/google/uuid"
	"github.com/Nivedithabp/receipt-processor/models"
	"github.com/Nivedithabp/receipt-processor/utils"
)

var (
	receiptStore = make(map[string]models.Receipt)
	pointsStore  = make(map[string]int)
	mutex        = &sync.Mutex{}
)

// ProcessReceipt processes a receipt and returns a unique ID
func ProcessReceipt(receipt models.Receipt) string {
	id := uuid.New().String()
	points := utils.CalculatePoints(receipt)

	mutex.Lock()
	receiptStore[id] = receipt
	pointsStore[id] = points
	mutex.Unlock()

	return id
}

// GetPoints returns the points for a given receipt ID
func GetPoints(id string) (int, bool) {
	mutex.Lock()
	points, exists := pointsStore[id]
	mutex.Unlock()

	return points, exists
}
