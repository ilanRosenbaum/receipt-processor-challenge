package service

import (
	"receipt-processor/internal/models"
	"testing"
)

func TestReceiptProcessing(t *testing.T) {
	t.Run("Validation", func(t *testing.T) {
		tests := []struct {
			name      string
			receipt   models.Receipt
			wantError bool
			errorMsg  string
		}{
			{
				name: "valid receipt with basic data",
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
				name: "invalid retailer with special characters",
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
				name: "incorrect date format",
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
				name: "incorrect time format",
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
				name: "invalid total format",
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
				name: "empty items list",
				receipt: models.Receipt{
					Retailer:     "Target",
					PurchaseDate: "2022-01-01",
					PurchaseTime: "13:01",
					Items:        []models.Item{},
					Total:        "1.00",
				},
				wantError: true,
				errorMsg:  "at least one item is required",
			},
			{
				name: "invalid item description with special characters",
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
				name: "invalid item price format",
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
	})

	t.Run("Points Calculation", func(t *testing.T) {
		tests := []struct {
			name    string
			receipt models.Receipt
			want    int64
			desc    string
		}{
			{
				name: "all rules combined",
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
				want: 105,
				desc: "6 retailer + 50 round + 25 multiple + 5 pairs + 10 afternoon + 3 desc + 6 odd day",
			},
			{
				name: "description length rule",
				receipt: models.Receipt{
					Retailer:     "Target",
					PurchaseDate: "2022-01-01",
					PurchaseTime: "13:01",
					Items: []models.Item{
						{ShortDescription: "123456789", Price: "5.00"},
					},
					Total: "5.00",
				},
				want: 88,
				desc: "6 retailer + 50 round + 25 multiple + 1 desc + 6 odd day",
			},
			{
				name: "special characters in retailer name",
				receipt: models.Receipt{
					Retailer:     "M&M Corner-Market",
					PurchaseDate: "2022-01-01",
					PurchaseTime: "14:30",
					Items: []models.Item{
						{ShortDescription: "Item", Price: "1.00"},
					},
					Total: "1.00",
				},
				want: 105,
				desc: "14 retailer + 50 round + 25 multiple + 10 afternoon + 6 odd day",
			},
			{
				name: "error handling - invalid total",
				receipt: models.Receipt{
					Retailer:     "Target",
					PurchaseDate: "2022-01-01",
					PurchaseTime: "13:01",
					Items: []models.Item{
						{ShortDescription: "Item", Price: "1.00"},
					},
					Total: "invalid",
				},
				want: 12,
				desc: "6 retailer + 6 odd day",
			},
			{
				name: "error handling - invalid price",
				receipt: models.Receipt{
					Retailer:     "Target",
					PurchaseDate: "2022-01-01",
					PurchaseTime: "14:30",
					Items: []models.Item{
						{ShortDescription: "123", Price: "invalid"},
					},
					Total: "1.00",
				},
				want: 97,
				desc: "6 retailer + 50 round + 25 multiple + 10 afternoon + 6 odd day",
			},
			{
				name: "time boundaries - at 14:00",
				receipt: models.Receipt{
					Retailer:     "Target",
					PurchaseDate: "2022-01-01",
					PurchaseTime: "14:00",
					Items: []models.Item{
						{ShortDescription: "Item", Price: "1.00"},
					},
					Total: "1.00",
				},
				want: 87,
				desc: "6 retailer + 50 round + 25 multiple + 6 odd day",
			},
			{
				name: "time boundaries - at 16:00",
				receipt: models.Receipt{
					Retailer:     "Target",
					PurchaseDate: "2022-01-01",
					PurchaseTime: "16:00",
					Items: []models.Item{
						{ShortDescription: "Item", Price: "1.00"},
					},
					Total: "1.00",
				},
				want: 87,
				desc: "6 retailer + 50 round + 25 multiple + 6 odd day",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := CalculatePoints(tt.receipt)
				if got != tt.want {
					t.Errorf("CalculatePoints() = %v, want %v\nBreakdown: %s", got, tt.want, tt.desc)
				}
			})
		}
	})
}
