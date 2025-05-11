package main

import (
	"log"
)

func analyzeTransaction(transaction Transaction, aiScore float64) float64 {
	score := aiScore // Start with the AI score
	db := GetDB()

	// Fetch rules from PostgreSQL
	rows, err := db.Query("SELECT id, weight, condition FROM detection_rules")
	if err != nil {
		log.Printf("Failed to fetch detection rules: %v", err)
		return score // Return AI score if rules can't be fetched
	}
	defer rows.Close()

	for rows.Next() {
		var rule DetectionRule
		if err := rows.Scan(&rule.ID, &rule.Weight, &rule.Condition); err != nil {
			log.Printf("Failed to scan rule: %v", err)
			continue
		}
		// Simple condition evaluation (can be made more sophisticated)
		if rule.Condition == "large_amount" && transaction.Amount > 10000 {
			score += rule.Weight
			log.Printf("Transaction %s triggered rule %d", transaction.ID, rule.ID)
		}
		if rule.Condition == "blacklisted_account" && contains(getBlacklistedAccounts(), transaction.AccountID) {
			score += rule.Weight
			log.Printf("Transaction %s triggered rule %d", transaction.ID, rule.ID)
		}
	}

	return score
}

func getBlacklistedAccounts() []string {
	// In a real system, this might come from a database or config
	return []string{"BLCKLST001", "BLCKLST002"}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
