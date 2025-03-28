package utils

import (
	//"fmt"
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

	// 1 point per alphanumeric character in retailer name
	retailerPoints := countAlphanumeric(receipt.Retailer)
	points += retailerPoints
	//fmt.Printf("Retailer Points: %d\n", retailerPoints)

	// 50 points if total is a round dollar amount with no cents
	if isRoundDollar(receipt.Total) {
		points += 50
		//fmt.Println("Round Dollar Points: 50")
	}

	// 25 points if total is multiple of 0.25
	if isMultipleOf25(receipt.Total) {
		points += 25
		//fmt.Println("Multiple of 0.25 Points: 25")
	}

	// 5 points for every 2 items - Only apply if there are items
	if len(receipt.Items) > 0 {
		itemPairsPoints := (len(receipt.Items) / 2) * 5
		points += itemPairsPoints
		//fmt.Printf("Item Pairs Points: %d\n", itemPairsPoints)
	}

	// Points for item description length multiple of 3
	for _, item := range receipt.Items {
		desc := strings.TrimSpace(item.ShortDescription)
		if len(desc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			itemPoints := int(math.Ceil(price * 0.2))
			points += itemPoints
			//fmt.Printf("Item Desc Multiple of 3 Points (%.2f * 0.2 = %d): %d\n", price, itemPoints, itemPoints)
		}
	}

	// 6 points if day of purchase is odd
	if isDayOdd(receipt.PurchaseDate) {
		points += 6
		//fmt.Println("Odd Day Points: 6")
	}

	// 10 points if purchase time is between 2 PM and 4 PM
	if isBetweenTwoAndFour(receipt.PurchaseTime) {
		points += 10
		//fmt.Println("Time Between 2-4 PM Points: 10")
	}

	//fmt.Printf("Total Points: %d\n", points)
	return points
}

// countAlphanumeric counts alphanumeric characters
func countAlphanumeric(s string) int {
	re := regexp.MustCompile(`[a-zA-Z0-9]`)
	return len(re.FindAllString(s, -1))
}

// isRoundDollar checks if total is a round dollar
func isRoundDollar(total string) bool {
	price, _ := strconv.ParseFloat(total, 64)
	return price == float64(int(price))
}

// isMultipleOf25 checks if total is a multiple of 0.25
func isMultipleOf25(total string) bool {
	price, _ := strconv.ParseFloat(total, 64)
	return math.Mod(price, 0.25) == 0
}

// isDayOdd checks if purchase date is odd
func isDayOdd(dateStr string) bool {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}
	return date.Day()%2 != 0
}

// isBetweenTwoAndFour checks if time is between 2 PM and 4 PM
func isBetweenTwoAndFour(timeStr string) bool {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return false
	}
	return t.Hour() == 14 || (t.Hour() == 15 && t.Minute() <= 59)
}
