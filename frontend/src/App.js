import React, { useState, useEffect } from 'react';
import TransactionForm from './components/TransactionForm';
import AlertList from './components/AlertList';
import './App.css';

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL;
const API_USER = process.env.REACT_APP_API_USER;
const API_PASSWORD = process.env.REACT_APP_API_PASSWORD;

console.log("API_BASE_URL:", API_BASE_URL); // ✅ Helps you confirm it's set correctly

function App() {
  const [alerts, setAlerts] = useState([]);
  const [loadingAlerts, setLoadingAlerts] = useState(false);
  const [errorAlerts, setErrorAlerts] = useState(null);
  const [submissionMessage, setSubmissionMessage] = useState(null);
  const [submissionError, setSubmissionError] = useState(null);

  const handleTransactionSubmit = async (transactionData) => {
    setSubmissionMessage(null);
    setSubmissionError(null);
    try {
      const response = await fetch(`${API_BASE_URL}/transactions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Basic ' + btoa(`${API_USER}:${API_PASSWORD}`),
        },
        body: JSON.stringify(transactionData),
        credentials: 'include', // Optional for CORS debugging
      });
      if (response.ok) {
        setSubmissionMessage('Transaction submitted successfully!');
        fetchAlerts(); // Refresh alerts
        setTimeout(() => setSubmissionMessage(null), 3000);
      } else {
        const errorData = await response.json().catch(() => null);
        setSubmissionError(`Failed to submit transaction: ${errorData?.message || response.statusText}`);
      }
    } catch (error) {
      setSubmissionError(`Error submitting transaction: ${error.message}`);
    }
  };

  const fetchAlerts = async () => {
    setLoadingAlerts(true);
    setErrorAlerts(null);
    try {
      const response = await fetch(`${API_BASE_URL}/alerts`, {
        headers: {
          'Authorization': 'Basic ' + btoa(`${API_USER}:${API_PASSWORD}`),
        },
        credentials: 'include', // Optional for CORS debugging
      });
      if (response.ok) {
        const data = await response.json();
        setAlerts(data);
      } else {
        setErrorAlerts(`Failed to fetch alerts: ${response.statusText}`);
      }
    } catch (error) {
      setErrorAlerts(`Error fetching alerts: ${error.message}`);
    } finally {
      setLoadingAlerts(false);
    }
  };

  useEffect(() => {
    fetchAlerts();
  }, []);

  return (
    <div className="App">
      <h1>Financial Crime Detector</h1>
      <TransactionForm onSubmit={handleTransactionSubmit} />
      {submissionMessage && <p className="success">{submissionMessage}</p>}
      {submissionError && <p className="error">{submissionError}</p>}
      <h2>Recent Alerts</h2>
      {loadingAlerts && <p>Loading alerts...</p>}
      {errorAlerts && <p className="error">{errorAlerts}</p>}
      <AlertList alerts={alerts} />
    </div>
  );
}

export default App;
