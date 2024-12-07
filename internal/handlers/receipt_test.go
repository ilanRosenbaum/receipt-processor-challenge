package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"receipt-processor/internal/models"
	"receipt-processor/internal/store"
	"testing"
)

func TestProcessReceipt(t *testing.T) {
	tests := []struct {
		name         string
		receipt      models.Receipt
		expectedCode int
		expectedID   bool // whether we expect an ID in response
		invalidJSON  bool // whether to send invalid JSON
	}{
		{
			name: "Valid Receipt",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
				},
				Total: "6.49",
			},
			expectedCode: http.StatusOK,
			expectedID:   true,
		},
		{
			name: "Invalid Retailer",
			receipt: models.Receipt{
				Retailer:     "Target!!!", // invalid characters
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
				},
				Total: "6.49",
			},
			expectedCode: http.StatusBadRequest,
			expectedID:   false,
		},
		{
			name: "Invalid Date Format",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "01-01-2022", // wrong format
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
				},
				Total: "6.49",
			},
			expectedCode: http.StatusBadRequest,
			expectedID:   false,
		},
		{
			name: "Invalid Time Format",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "1:01 PM", // wrong format
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
				},
				Total: "6.49",
			},
			expectedCode: http.StatusBadRequest,
			expectedID:   false,
		},
		{
			name: "Empty Items Array",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items:        []models.Item{},
				Total:        "0.00",
			},
			expectedCode: http.StatusBadRequest,
			expectedID:   false,
		},
		{
			name:         "Invalid JSON",
			invalidJSON:  true,
			expectedCode: http.StatusBadRequest,
			expectedID:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := store.NewStore()
			handler := NewReceiptHandler(store)

			var body []byte
			if tt.invalidJSON {
				body = []byte(`{invalid json}`)
			} else {
				body, _ = json.Marshal(tt.receipt)
			}

			req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			handler.ProcessReceipt(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.expectedCode)
			}

			if tt.expectedID {
				var response models.ReceiptResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Errorf("couldn't decode response: %v", err)
				}
				if response.ID == "" {
					t.Error("expected non-empty ID")
				}
			}
		})
	}
}

func TestGetPoints(t *testing.T) {
	tests := []struct {
		name           string
		setupID        string
		setupPoints    int64
		requestID      string
		expectedCode   int
		expectedPoints int64
	}{
		{
			name:           "Valid Points Request",
			setupID:        "test-id-1",
			setupPoints:    100,
			requestID:      "test-id-1",
			expectedCode:   http.StatusOK,
			expectedPoints: 100,
		},
		{
			name:         "Non-existent Receipt",
			requestID:    "non-existent",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Empty ID",
			requestID:    "",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := store.NewStore()
			handler := NewReceiptHandler(store)

			// Setup test data if needed
			if tt.setupID != "" {
				store.SaveReceipt(tt.setupID, models.Receipt{}, tt.setupPoints)
			}

			// Create request with mux vars
			req := httptest.NewRequest("GET", "/receipts/{id}/points", nil)
			vars := map[string]string{
				"id": tt.requestID,
			}
			req = mux.SetURLVars(req, vars)

			rr := httptest.NewRecorder()
			handler.GetPoints(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.expectedCode)
			}

			if tt.expectedCode == http.StatusOK {
				var response models.PointsResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Errorf("couldn't decode response: %v", err)
				}
				if response.Points != tt.expectedPoints {
					t.Errorf("expected %d points, got %d", tt.expectedPoints, response.Points)
				}
			}
		})
	}
}
