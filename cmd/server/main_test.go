package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "receipt-processor/internal/models"
)

func TestSetupServer(t *testing.T) {
    srv := setupServer()
    
    // Create test server
    testServer := httptest.NewServer(srv)
    defer testServer.Close()
    
    // Test health endpoint
    t.Run("Health Check", func(t *testing.T) {
        resp, err := http.Get(testServer.URL + "/health")
        if err != nil {
            t.Fatalf("Could not send GET request: %v", err)
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK {
            t.Errorf("Expected status OK; got %v", resp.StatusCode)
        }
    })

    // Test POST /receipts/process
    t.Run("Process Receipt", func(t *testing.T) {
        receipt := models.Receipt{
            Retailer:     "Target",
            PurchaseDate: "2022-01-01",
            PurchaseTime: "13:01",
            Items: []models.Item{
                {ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
            },
            Total: "6.49",
        }

        jsonData, err := json.Marshal(receipt)
        if err != nil {
            t.Fatalf("Failed to marshal receipt: %v", err)
        }

        resp, err := http.Post(testServer.URL+"/receipts/process", "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            t.Fatalf("Could not send POST request: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            t.Errorf("Expected status OK; got %v", resp.StatusCode)
        }

        var response models.ReceiptResponse
        if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        if response.ID == "" {
            t.Error("Expected non-empty receipt ID")
        }
    })

    // Test GET /receipts/{id}/points
    t.Run("Get Points", func(t *testing.T) {
        // First create a receipt
        receipt := models.Receipt{
            Retailer:     "Target",
            PurchaseDate: "2022-01-01",
            PurchaseTime: "13:01",
            Items: []models.Item{
                {ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
            },
            Total: "6.49",
        }

        jsonData, err := json.Marshal(receipt)
        if err != nil {
            t.Fatalf("Failed to marshal receipt: %v", err)
        }

        resp, err := http.Post(testServer.URL+"/receipts/process", "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            t.Fatalf("Could not send POST request: %v", err)
        }

        var receiptResponse models.ReceiptResponse
        if err := json.NewDecoder(resp.Body).Decode(&receiptResponse); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }
        resp.Body.Close()

        // Now get points for the receipt
        resp, err = http.Get(testServer.URL + "/receipts/" + receiptResponse.ID + "/points")
        if err != nil {
            t.Fatalf("Could not send GET request: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            t.Errorf("Expected status OK; got %v", resp.StatusCode)
        }

        var pointsResponse models.PointsResponse
        if err := json.NewDecoder(resp.Body).Decode(&pointsResponse); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        if pointsResponse.Points <= 0 {
            t.Error("Expected points to be greater than 0")
        }
    })

    // Test invalid endpoints
    t.Run("Invalid Endpoint", func(t *testing.T) {
        resp, err := http.Get(testServer.URL + "/invalid")
        if err != nil {
            t.Fatalf("Could not send GET request: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusNotFound {
            t.Errorf("Expected status NotFound; got %v", resp.StatusCode)
        }
    })
}

func TestMain(m *testing.M) {
    go func() {
        main()
    }()
    m.Run()
}