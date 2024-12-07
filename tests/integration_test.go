package tests

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "receipt-processor/internal/models"
    "receipt-processor/internal/handlers"
    "receipt-processor/internal/store"
    "github.com/gorilla/mux"
)


func setupRouter() http.Handler {
    store := store.NewStore()
    handler := handlers.NewReceiptHandler(store)
    
    router := mux.NewRouter()
    router.HandleFunc("/receipts/process", handler.ProcessReceipt).Methods("POST")
    router.HandleFunc("/receipts/{id}/points", handler.GetPoints).Methods("GET")
    return router
}

func TestIntegrationReceiptProcessing(t *testing.T) {
    router := setupRouter()
    server := httptest.NewServer(router)
    defer server.Close()

    t.Run("Full Receipt Processing Flow", func(t *testing.T) {
        receipt := models.Receipt{
            Retailer:     "Target",
            PurchaseDate: "2022-01-01",
            PurchaseTime: "13:01",
            Items: []models.Item{
                {ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
            },
            Total: "6.49",
        }

        // Process receipt
        receiptJSON, _ := json.Marshal(receipt)
        resp, err := http.Post(fmt.Sprintf("%s/receipts/process", server.URL), 
            "application/json", 
            bytes.NewBuffer(receiptJSON))

        if err != nil || resp.StatusCode != http.StatusOK {
            t.Fatalf("Failed to process receipt: %v", err)
        }

        var receiptResponse models.ReceiptResponse
        if err := json.NewDecoder(resp.Body).Decode(&receiptResponse); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }
        resp.Body.Close()

        // Get points
        resp, err = http.Get(fmt.Sprintf("%s/receipts/%s/points", server.URL, receiptResponse.ID))
        if err != nil || resp.StatusCode != http.StatusOK {
            t.Fatalf("Failed to get points: %v", err)
        }

        var pointsResponse models.PointsResponse
        if err := json.NewDecoder(resp.Body).Decode(&pointsResponse); err != nil {
            t.Fatalf("Failed to decode points response: %v", err)
        }
        resp.Body.Close()

        expectedPoints := int64(12) // 6 for retailer name + 6 for odd day
        if pointsResponse.Points != expectedPoints {
            t.Errorf("Expected %d points, got %d", expectedPoints, pointsResponse.Points)
        }
    })

    t.Run("Get Points for Non-existent Receipt", func(t *testing.T) {
        resp, err := http.Get(fmt.Sprintf("%s/receipts/nonexistent/points", server.URL))
        if err != nil || resp.StatusCode != http.StatusNotFound {
            t.Errorf("Expected 404 status, got %d", resp.StatusCode)
        }
    })
}