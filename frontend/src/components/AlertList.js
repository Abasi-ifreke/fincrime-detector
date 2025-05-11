import React from 'react';

function AlertList({ alerts }) {
  if (!alerts || alerts.length === 0) {
    return <p>No recent alerts.</p>;
  }

  return (
    <ul>
      {alerts.map((alert) => (
        <li key={alert.id}>
          <strong>Alert ID:</strong> {alert.id}, <strong>Account:</strong> {alert.accountID}, <strong>Reason:</strong> {alert.reason}, <strong>Score:</strong> {alert.score.toFixed(2)}, <strong>Time:</strong> {new Date(alert.timestamp).toLocaleString()}
        </li>
      ))}
    </ul>
  );
}

export default AlertList;