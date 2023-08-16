package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Items        []Item  `json:"items"`
	Total        string  `json:"total"`
}

type ProcessedReceipt struct {
	ID string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/receipts/process", processReceiptHandler).Methods(http.MethodPost)
	r.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods(http.MethodGet)

	http.Handle("/", r)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Calculate points based on receipt rules
	points := calculatePoints(&receipt)

	// Generate a random ID (for demonstration purposes)
	id := generateRandomID()

	response := ProcessedReceipt{ID: id}
	jsonResponse, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Fetch receipt based on ID (in-memory or database)
	// For demonstration purposes, assuming the receipt is already processed

	// Calculate points based on receipt rules
	points := calculatePoints(&receipt) // Replace 'receipt' with actual fetched receipt

	response := PointsResponse{Points: points}
	jsonResponse, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func calculatePoints(receipt *Receipt) int {
	points := len(receipt.Retailer)

	totalFloat, _ := strconv.ParseFloat(receipt.Total, 64)
	if totalFloat == math.Floor(totalFloat) && totalFloat > 0 {
		points += 50
	}

	totalFloat, _ = strconv.ParseFloat(receipt.Total, 64)
	if math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	}

	points += len(receipt.Items) / 2 * 5

	for _, item := range receipt.Items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			itemPrice, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(itemPrice * 0.2))
		}
	}

	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if purchaseTime.Hour() > 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 == 1 {
		points += 6
	}

	return points
}

func generateRandomID() string {
	return "7fb1377b-b223-49d9-a31a-5a02701dd310" // Replace with actual ID generation logic
}

func fetchReceiptByID(id string) Receipt {
	// Fetch receipt by ID from in-memory storage or database
	// For demonstration purposes, a mock receipt is returned
	return Receipt{
		Retailer:     "Mock Retailer",
		PurchaseDate: "2023-08-16",
		PurchaseTime: "15:30",
		Items: []Item{
			{ShortDescription: "Item 1", Price: "5.99"},
			{ShortDescription: "Item 2", Price: "2.50"},
		},
		Total: "8.49",
	}
}
