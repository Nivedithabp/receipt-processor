package models

// Receipt represents a receipt submitted by the user
// @Description Receipt model for processing
type Receipt struct {
	Retailer     string `json:"retailer" example:"Target"`
	PurchaseDate string `json:"purchaseDate" example:"2022-01-01"`
	PurchaseTime string `json:"purchaseTime" example:"13:01"`
	Items        []Item `json:"items"`
	Total        string `json:"total" example:"35.35"`
}

// Item represents an item in the receipt
// @Description Item model for the receipt
type Item struct {
	ShortDescription string `json:"shortDescription" example:"Mountain Dew 12PK"`
	Price            string `json:"price" example:"6.49"`
}

