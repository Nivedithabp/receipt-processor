package utils

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Nivedithabp/receipt-processor/models"
)

// CalculatePoints calculates the points for the receipt
func CalculatePoints(receipt models.Receipt) int {
	points := 0

	// 1 point per alphanumeric character in the retailer name
	points += countAlphanumeric(receipt.Retailer)

	// 50 points if total is a round dollar
	if isRoundDollar(receipt.Total) {
		points += 50
	}

	// 25 points if total is multiple of 0.25
	if isMultipleOf25(receipt.Total) {
		points += 25
	}

	// 5 points for every 2 items
	points += (len(receipt.Items) / 2) * 5

	// Points for item description length multiple of 3
	for _, item := range receipt.Items {
		desc := strings.TrimSpace(item.ShortDescription)
		if len(desc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	// 6 points if day of purchase is odd
	if isDayOdd(receipt.PurchaseDate) {
		points += 6
	}

	// 10 points if purchase time is between 2 PM and 4 PM
	if isBetweenTwoAndFour(receipt.PurchaseTime) {
		points += 10
	}

	return points
}

func countAlphanumeric(s string) int {
	re := regexp.MustCompile(`[a-zA-Z0-9]`)
	return len(re.FindAllString(s, -1))
}

func isRoundDollar(total string) bool {
	price, _ := strconv.ParseFloat(total, 64)
	return price == float64(int(price))
}

func isMultipleOf25(total string) bool {
	price, _ := strconv.ParseFloat(total, 64)
	return math.Mod(price, 0.25) == 0
}

func isDayOdd(dateStr string) bool {
	date, _ := time.Parse("2006-01-02", dateStr)
	return date.Day()%2 != 0
}

func isBetweenTwoAndFour(timeStr string) bool {
	t, _ := time.Parse("15:04", timeStr)
	return t.Hour() == 14 || (t.Hour() == 15 && t.Minute() <= 59)
}
