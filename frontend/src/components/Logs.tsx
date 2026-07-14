import React, { useState, useEffect } from 'react';
import './Logs.css';
import { API_BASE } from '../api';

interface LogEntry {
  timestamp: string;
  level: string;
  message: string;
}

const Logs: React.FC = () => {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchLogs = async () => {
      try {
        const response = await fetch(`${API_BASE}/v1/logs`);
        if (!response.ok) {
          throw new Error('Failed to fetch logs');
        }
        const data = await response.json();
        setLogs(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An unknown error occurred');
      } finally {
        setIsLoading(false);
      }
    };

    fetchLogs();
    
    // Poll for new logs every 5 seconds
    const interval = setInterval(fetchLogs, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="logs-container">
      <h1>Backend Logs</h1>
      {isLoading && <p>Loading logs...</p>}
      {error && <p className="error">Error: {error}</p>}
      <div className="logs-list">
        {logs.length === 0 && !isLoading && <p>No logs available.</p>}
        {logs.map((log, index) => (
          <div key={index} className={`log-entry ${log.level.toLowerCase()}`}>
            <span className="timestamp">{log.timestamp}</span>
            <span className="level">{log.level}</span>
            <span className="message">{log.message}</span>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Logs;