package handlers

import (
    "encoding/json"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "net/http"
    "receipt-processor/internal/models"
    "receipt-processor/internal/service"
    "receipt-processor/internal/store"
)

type ReceiptHandler struct {
    store *store.ReceiptStore
}

func NewReceiptHandler(store *store.ReceiptStore) *ReceiptHandler {
    return &ReceiptHandler{store: store}
}

func (h *ReceiptHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
    var receipt models.Receipt
    if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    if err := service.ValidateReceipt(receipt); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    id := uuid.New().String()
    points := service.CalculatePoints(receipt)
    
    h.store.SaveReceipt(id, receipt, points)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(models.ReceiptResponse{ID: id})
}

func (h *ReceiptHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    
    points, exists := h.store.GetPoints(id)
    if !exists {
        http.Error(w, "Receipt not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(models.PointsResponse{Points: points})
}