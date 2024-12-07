package service

import (
	"fmt"
	"math"
	"receipt-processor/internal/models"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func CalculatePoints(receipt models.Receipt) int64 {
	var points int64 = 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	alphanumeric := regexp.MustCompile(`[a-zA-Z0-9]`)
	matches := alphanumeric.FindAllString(receipt.Retailer, -1)
	points += int64(len(matches))

	// Rule 2: 50 points if the total is a round dollar amount
	if strings.HasSuffix(receipt.Total, ".00") {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if total, err := strconv.ParseFloat(receipt.Total, 64); err == nil {
		if math.Mod(total*100, 25) == 0 {
			points += 25
		}
	}

	// Rule 4: 5 points for every two items
	points += int64(len(receipt.Items) / 2 * 5)

	// Rule 5: Points for items with description length multiple of 3
	for _, item := range receipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int64(math.Ceil(price * 0.2))
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd
	if date, err := time.Parse("2006-01-02", receipt.PurchaseDate); err == nil {
		if date.Day()%2 == 1 {
			points += 6
		}
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	if purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		afterTwo := time.Date(2000, 1, 1, 14, 0, 0, 0, time.UTC)
		beforeFour := time.Date(2000, 1, 1, 16, 0, 0, 0, time.UTC)
		compareTime := time.Date(2000, 1, 1, purchaseTime.Hour(), purchaseTime.Minute(), 0, 0, time.UTC)

		if compareTime.After(afterTwo) && compareTime.Before(beforeFour) {
			points += 10
		}
	}

	return points
}

func ValidateReceipt(receipt models.Receipt) error {
	// Validate retailer
	if !regexp.MustCompile(`^[\w\s\-&]+$`).MatchString(receipt.Retailer) {
		return fmt.Errorf("invalid retailer name")
	}

	// Validate purchase date
	if _, err := time.Parse("2006-01-02", receipt.PurchaseDate); err != nil {
		return fmt.Errorf("invalid purchase date")
	}

	// Validate purchase time
	if _, err := time.Parse("15:04", receipt.PurchaseTime); err != nil {
		return fmt.Errorf("invalid purchase time")
	}

	// Validate total
	if !regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(receipt.Total) {
		return fmt.Errorf("invalid total")
	}

	// Validate items
	if len(receipt.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}

	for _, item := range receipt.Items {
		if !regexp.MustCompile(`^[\w\s\-]+$`).MatchString(item.ShortDescription) {
			return fmt.Errorf("invalid item description")
		}
		if !regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(item.Price) {
			return fmt.Errorf("invalid item price")
		}
	}

	return nil
}
