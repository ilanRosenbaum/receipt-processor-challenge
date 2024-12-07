package store

import (
	"receipt-processor/internal/models"
	"testing"
)

func TestReceiptStore(t *testing.T) {
	store := NewStore()

	testReceipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Test Item", Price: "10.00"},
		},
		Total: "10.00",
	}

	// Test storing receipt and points
	t.Run("Save and Retrieve Points", func(t *testing.T) {
		testID := "test-id-1"
		testPoints := int64(50)

		store.SaveReceipt(testID, testReceipt, testPoints)

		points, exists := store.GetPoints(testID)
		if !exists {
			t.Error("Receipt not found in store")
		}

		if points != testPoints {
			t.Errorf("Got points %d, want %d", points, testPoints)
		}
	})

	// Test retrieving non-existent receipt
	t.Run("Get Non-existent Receipt", func(t *testing.T) {
		_, exists := store.GetPoints("non-existent-id")
		if exists {
			t.Error("Expected non-existent receipt to return exists=false")
		}
	})

	// Test concurrent access
	t.Run("Concurrent Access", func(t *testing.T) {
		done := make(chan bool)

		// Launch multiple goroutines to test concurrent access
		for i := 0; i < 10; i++ {
			go func(index int) {
				id := string(rune('A' + index))
				store.SaveReceipt(id, testReceipt, int64(index))
				points, _ := store.GetPoints(id)
				if points != int64(index) {
					t.Errorf("Got points %d, want %d", points, index)
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}
