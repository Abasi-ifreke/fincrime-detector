package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := GetDB().Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Database not healthy: %v", err)))
		return
	}
	// Add Elasticsearch health check if needed
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handleTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var transaction Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		logger.Errorf("Error decoding transaction request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Infof("Received transaction: %+v", transaction) // Log the received transaction

	// **Process the transaction here:**
	// - Basic validation (e.g., are required fields present?)
	if transaction.AccountID == "" || transaction.TransactionType == "" || transaction.Amount == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	anomalyScore, explanation := simulateAIModel(transaction)
	finalScore := analyzeTransaction(transaction, anomalyScore)

	logger.Infof("AI Score: %.2f, Rule Score: %.2f, Final Score: %.2f for transaction %s", anomalyScore, finalScore-anomalyScore, finalScore, transaction.ID)

	if finalScore >= 0.7 {
		alert := Alert{
			ID:            fmt.Sprintf("alert-%d", time.Now().UnixNano()),
			TransactionID: transaction.ID,
			AccountID:     transaction.AccountID,
			Reason:        fmt.Sprintf("High anomaly score detected (AI: %.2f, Final Score: %.2f. Explanation: %s", anomalyScore, finalScore, explanation),
			Score:         finalScore,
			Timestamp:     time.Now(),
		}
		err := indexAlert(alert)
		if err != nil {
			logger.Errorf("Failed to index alert: %v", err)
		}
		alertsGenerated.Inc()
		logger.WithFields(logrus.Fields{
			"id":             alert.ID,
			"account_id":     alert.AccountID,
			"reason":         alert.Reason,
			"score":          alert.Score,
			"transaction_id": alert.TransactionID,
		}).Warn("Alert generated")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transaction processed"))
}

func handleGetAlerts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	alerts, err := queryAlerts()
	if err != nil {
		logger.Errorf("Failed to retrieve alerts: %v", err)
		http.Error(w, "Failed to retrieve alerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}
