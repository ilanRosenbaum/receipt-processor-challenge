package service

import (
    "testing"
    "receipt-processor/internal/models"
)

func TestValidateReceipt_CompleteCoverage(t *testing.T) {
    tests := []struct {
        name       string
        receipt    models.Receipt
        wantError  bool
        errorMsg   string
    }{
        {
            name: "Valid Receipt",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "13:01",
                Items: []models.Item{
                    {ShortDescription: "Item", Price: "1.00"},
                },
                Total: "1.00",
            },
            wantError: false,
        },
        {
            name: "Invalid Retailer Characters",
            receipt: models.Receipt{
                Retailer:     "Target!!!",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "13:01",
                Items: []models.Item{
                    {ShortDescription: "Item", Price: "1.00"},
                },
                Total: "1.00",
            },
            wantError: true,
            errorMsg:  "invalid retailer name",
        },
        {
            name: "Invalid Purchase Date Format",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "01-01-2022",
                PurchaseTime: "13:01",
                Items: []models.Item{
                    {ShortDescription: "Item", Price: "1.00"},
                },
                Total: "1.00",
            },
            wantError: true,
            errorMsg:  "invalid purchase date",
        },
        {
            name: "Invalid Purchase Time Format",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "1:1",
                Items: []models.Item{
                    {ShortDescription: "Item", Price: "1.00"},
                },
                Total: "1.00",
            },
            wantError: true,
            errorMsg:  "invalid purchase time",
        },
        {
            name: "Invalid Total Format",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "13:01",
                Items: []models.Item{
                    {ShortDescription: "Item", Price: "1.00"},
                },
                Total: "1",
            },
            wantError: true,
            errorMsg:  "invalid total",
        },
        {
            name: "Empty Items Array",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "13:01",
                Items:        []models.Item{},
                Total:       "1.00",
            },
            wantError: true,
            errorMsg:  "at least one item is required",
        },
        {
            name: "Invalid Item Description",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "13:01",
                Items: []models.Item{
                    {ShortDescription: "Item!!!", Price: "1.00"},
                },
                Total: "1.00",
            },
            wantError: true,
            errorMsg:  "invalid item description",
        },
        {
            name: "Invalid Item Price Format",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "13:01",
                Items: []models.Item{
                    {ShortDescription: "Item", Price: "1"},
                },
                Total: "1.00",
            },
            wantError: true,
            errorMsg:  "invalid item price",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateReceipt(tt.receipt)
            if (err != nil) != tt.wantError {
                t.Errorf("ValidateReceipt() error = %v, wantError %v", err, tt.wantError)
                return
            }
            if tt.wantError && err != nil && err.Error() != tt.errorMsg {
                t.Errorf("ValidateReceipt() error message = %v, want %v", err.Error(), tt.errorMsg)
            }
        })
    }
}

func TestCalculatePoints_CompleteCoverage(t *testing.T) {
    tests := []struct {
        name     string
        receipt  models.Receipt
        expected int64
        explanation string
    }{
        {
            name: "All Rules Applied",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "14:30",
                Items: []models.Item{
                    {ShortDescription: "123", Price: "1.00"},
                    {ShortDescription: "456", Price: "2.00"},
                    {ShortDescription: "789", Price: "3.00"},
                },
                Total: "6.00",
            },
            expected: 105, // 6 (retailer) + 50 (round) + 25 (multiple 0.25) + 5 (pair) + 
                         // 6 (odd day) + 10 (afternoon) + 1 (ceil(1.00*0.2)) + 1 (ceil(2.00*0.2)) + 1 (ceil(3.00*0.2))
            explanation: "All point rules applied",
        },
        {
            name: "Invalid Total Parse",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "13:01",
                Items: []models.Item{
                    {ShortDescription: "Item", Price: "1.00"},
                },
                Total: "invalid",
            },
            expected: 12, // 6 (retailer) + 6 (odd day)
            explanation: "Should handle invalid total parse",
        },
        {
            name: "Invalid Item Price Parse",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "14:30",
                Items: []models.Item{
                    {ShortDescription: "123", Price: "invalid"},
                },
                Total: "1.00",
            },
            expected: 97, // 6 (retailer) + 50 (round) + 25 (multiple 0.25) + 6 (odd day) + 10 (afternoon)
            explanation: "Should handle invalid item price parse",
        },
        {
            name: "Invalid Date Parse",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "invalid",
                PurchaseTime: "14:30",
                Items: []models.Item{
                    {ShortDescription: "Item", Price: "1.00"},
                },
                Total: "1.00",
            },
            expected: 91, // 6 (retailer) + 50 (round) + 25 (multiple 0.25) + 10 (afternoon)
            explanation: "Should handle invalid date parse",
        },
        {
            name: "Invalid Time Parse",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "invalid",
                Items: []models.Item{
                    {ShortDescription: "Item", Price: "1.00"},
                },
                Total: "1.00",
            },
            expected: 87, // 6 (retailer) + 50 (round) + 25 (multiple 0.25) + 6 (odd day)
            explanation: "Should handle invalid time parse",
        },
        {
            name: "Time Edge Case - 14:00",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "14:00",
                Items: []models.Item{
                    {ShortDescription: "Item", Price: "1.00"},
                },
                Total: "1.00",
            },
            expected: 87, // 6 (retailer) + 50 (round) + 25 (multiple 0.25) + 6 (odd day)
            explanation: "Should handle 14:00 time edge case",
        },
        {
            name: "Time Edge Case - 16:00",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "16:00",
                Items: []models.Item{
                    {ShortDescription: "Item", Price: "1.00"},
                },
                Total: "1.00",
            },
            expected: 87, // 6 (retailer) + 50 (round) + 25 (multiple 0.25) + 6 (odd day)
            explanation: "Should handle 16:00 time edge case",
        },
        {
            name: "Special Characters in Retailer",
            receipt: models.Receipt{
                Retailer:     "M&M Corner-Market",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "14:30",
                Items: []models.Item{
                    {ShortDescription: "Item", Price: "1.00"},
                },
                Total: "1.00",
            },
            expected: 105, // 14 (alphanumeric chars) + 50 (round) + 25 (multiple 0.25) + 6 (odd day) + 10 (afternoon)
            explanation: "Should handle special characters in retailer name",
        },
        {
            name: "Description Multiple of 3 Points",
            receipt: models.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "13:01",
                Items: []models.Item{
                    {ShortDescription: "123456789", Price: "5.00"},
                },
                Total: "5.00",
            },
            expected: 88, // 6 (retailer) + 50 (round) + 6 (odd day) + 1 (ceil(5.00*0.2)) + 25 (multiple of 0.25)
            explanation: "Should calculate description length points correctly",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            points := CalculatePoints(tt.receipt)
            if points != tt.expected {
                t.Errorf("CalculatePoints() = %v, want %v\nExplanation: %s", 
                    points, tt.expected, tt.explanation)
            }
        })
    }
}