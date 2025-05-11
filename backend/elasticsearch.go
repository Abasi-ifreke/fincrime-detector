package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

var es *elasticsearch.Client

func initElasticsearch() {
	// Load the CA certificate
	caCertPath := os.Getenv("ELASTIC_CA_CERT")
	if caCertPath == "" {
		logger.Fatal("ELASTIC_CA_CERT env var is not set")
	}
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		logger.Fatalf("Error reading CA certificate from %s: %v", caCertPath, err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		logger.Fatalf("Failed to append CA certificate")
	}

	// Build the Elasticsearch client configuration.
	cfg := elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTICSEARCH_URL")},
		Username:  os.Getenv("ELASTICSEARCH_USER"),
		Password:  os.Getenv("ELASTICSEARCH_PASSWORD"),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Fatalf("Error creating Elasticsearch client: %v", err)
	}
	res, err := es.Info()
	if err != nil {
		logger.Fatalf("Error getting Elasticsearch info: %v", err)
	}
	defer res.Body.Close()
	logger.Printf("Connected to Elasticsearch: %s", res.String())

	// Ensure the alerts index exists
	ensureAlertsIndexExists()
}

func ensureAlertsIndexExists() {
	res, err := es.Indices.Exists([]string{"alerts"})
	if err != nil {
		logger.Fatalf("Error checking if 'alerts' index exists: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		mappings := `{
			"mappings": {
				"properties": {
					"id":            { "type": "keyword" },
					"transaction_id":{ "type": "keyword" },
					"account_id":    { "type": "keyword" },
					"reason":        { "type": "text" },
					"score":         { "type": "float" },
					"timestamp":     { "type": "date", "format": "epoch_millis" }
				}
			}
		}`
		createRes, err := es.Indices.Create("alerts", es.Indices.Create.WithBody(bytes.NewReader([]byte(mappings))))
		if err != nil {
			logger.Fatalf("Error creating 'alerts' index: %v", err)
		}
		defer createRes.Body.Close()
		if createRes.IsError() {
			logger.Fatalf("Failed to create 'alerts' index: %s", createRes.String())
		}
		logger.Println("'alerts' index created")
	} else {
		logger.Println("'alerts' index already exists")
	}
}

func indexAlert(alert Alert) error {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(alert)
	if err != nil {
		logger.Printf("Error encoding alert to JSON: %v", err)
		return err
	}

	res, err := es.Index(
		"alerts",
		&b,
		es.Index.WithDocumentID(alert.ID),
	)
	if err != nil {
		logger.Printf("Error indexing alert in Elasticsearch: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			logger.Printf("Error parsing Elasticsearch response: %v", err)
		} else {
			logger.Printf("Elasticsearch error: %v", e)
		}
		return err
	}
	return nil
}

func queryAlerts() ([]Alert, error) {
	res, err := es.Search(
		es.Search.WithIndex("alerts"),
		es.Search.WithSize(10),
		es.Search.WithBody(bytes.NewReader([]byte(`{
			"sort": [
				{ "timestamp": { "order": "desc" } }
			]
		}`))),
	)
	if err != nil {
		logger.Printf("Error querying Elasticsearch: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			logger.Printf("Error parsing Elasticsearch response: %v", err)
		} else {
			logger.Printf("Elasticsearch error: %v", e)
		}
		return nil, err
	}

	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		logger.Printf("Error decoding Elasticsearch response: %v", err)
		return nil, err
	}

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	var alerts []Alert
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]

		scoreValue := source.(map[string]interface{})["score"]
		var score float64
		switch v := scoreValue.(type) {
		case float64:
			score = v
		case string:
			score, err = strconv.ParseFloat(v, 64)
			if err != nil {
				logger.Printf("Error converting score string to float64: %v", err)
				continue
			}
		default:
			logger.Printf("Unexpected type for score: %T", v)
			continue
		}

		alert := Alert{
			ID:            source.(map[string]interface{})["id"].(string),
			TransactionID: source.(map[string]interface{})["transaction_id"].(string),
			AccountID:     source.(map[string]interface{})["account_id"].(string),
			Reason:        source.(map[string]interface{})["reason"].(string),
			Score:         score,
			Timestamp:     time.Unix(int64(source.(map[string]interface{})["timestamp"].(float64)/1000), 0),
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}
