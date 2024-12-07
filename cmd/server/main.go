package main

import (
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "receipt-processor/internal/handlers"
    "receipt-processor/internal/store"
)

func setupServer() http.Handler {
    store := store.NewStore()
    handler := handlers.NewReceiptHandler(store)

    router := mux.NewRouter()
    router.HandleFunc("/receipts/process", handler.ProcessReceipt).Methods("POST")
    router.HandleFunc("/receipts/{id}/points", handler.GetPoints).Methods("GET")
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }).Methods("GET")

    return router
}

func main() {
    router := setupServer()
    log.Printf("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}