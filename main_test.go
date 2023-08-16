package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestProcessReceiptHandler(t *testing.T) {
	// Prepare a test receipt JSON
	receiptJSON := []byte(`{
		"retailer": "Test Retailer",
		"purchaseDate": "2023-08-16",
		"purchaseTime": "15:30",
		"items": [
			{"shortDescription": "Test Item 1", "price": "9.99"},
			{"shortDescription": "Test Item 2", "price": "4.50"}
		],
		"total": "14.49"
	}`)

	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(receiptJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(processReceiptHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var response ProcessedReceipt
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Error decoding response JSON: %v", err)
	}

	if response.ID == "" {
		t.Errorf("Empty ID in response")
	}
}

func TestGetPointsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/receipts/test-id/points", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/receipts/{id}/points", getPointsHandler)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var response PointsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Error decoding response JSON: %v", err)
	}

	if response.Points != 100 { // Replace with the expected points
		t.Errorf("Expected points 100, got %d", response.Points)
	}
}
