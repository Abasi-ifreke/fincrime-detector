package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type Transaction struct {
	ID              string    `json:"id"`
	AccountID       string    `json:"account_id"`
	TransactionType string    `json:"transaction_type"`
	Amount          float64   `json:"amount"`
	Timestamp       time.Time `json:"timestamp"`
}

type DetectionRule struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Priority    int     `json:"priority"`
	Condition   string  `json:"condition"`
	Weight      float64 `json:"weight"`
}

type Alert struct {
	ID            string    `json:"id"`
	TransactionID string    `json:"transaction_id"`
	AccountID     string    `json:"account_id"`
	Reason        string    `json:"reason"`
	Score         float64   `json:"score"`
	Timestamp     time.Time `json:"timestamp"`
}

var (
	transactionsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "fincrime_detector_transactions_processed_total",
		Help: "Total number of transactions processed.",
	})
	alertsGenerated = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "fincrime_detector_alerts_generated_total",
		Help: "Total number of alerts generated.",
	})
	logger = logrus.New()
)

func init() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	logDir := "/var/log/app"
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	logFile, err := os.OpenFile(logDir+"/backend.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger.SetOutput(logFile)

	prometheus.MustRegister(transactionsProcessed)
	prometheus.MustRegister(alertsGenerated)

	initDatabase()
	initElasticsearch()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/transactions", authenticate(handleTransaction))
	mux.HandleFunc("/api/alerts", authenticate(handleGetAlerts))
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", handleHealth)

	// CORS support
	c := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			allowedOrigins := []string{
				"http://localhost",
				"http://localhost:80",
				"http://frontend",
				"http://frontend:80",
			}

			for _, o := range allowedOrigins {
				if origin == o {
					return true
				}
			}
			return false
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Infof("Starting fincrime-detector application on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != os.Getenv("API_USER") || password != os.Getenv("API_PASSWORD") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		next(w, r)
	}
}
