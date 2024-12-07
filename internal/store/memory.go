package store

import (
    "sync"
    "receipt-processor/internal/models"
)

type ReceiptStore struct {
    receipts map[string]models.Receipt
    points   map[string]int64
    mutex    sync.RWMutex
}

func NewStore() *ReceiptStore {
    return &ReceiptStore{
        receipts: make(map[string]models.Receipt),
        points:   make(map[string]int64),
    }
}

func (s *ReceiptStore) SaveReceipt(id string, receipt models.Receipt, points int64) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    s.receipts[id] = receipt
    s.points[id] = points
}

func (s *ReceiptStore) GetPoints(id string) (int64, bool) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    points, exists := s.points[id]
    return points, exists
}