import React, { useState, useEffect } from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const API_BASE = 'http://localhost:30128'; // Ensure port matches UI or uses config

export const Analytics: React.FC = () => {
  const [data, setData] = useState<{ usage: Record<string, any>, latencies: Record<string, number> }>({ usage: {}, latencies: {} });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${API_BASE}/v1/usage`)
      .then(r => r.json())
      .then(d => {
        setData(d || { usage: {}, latencies: {} });
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div style={{ padding: '2rem', color: 'var(--text-muted)' }}>Loading analytics...</div>;
  }

  // Transform data for charts
  const usageKeys = Object.keys(data.usage || {});
  
  const tokenData = usageKeys.map(key => ({
    name: key,
    Used: data.usage[key].tokens_used || 0,
    Saved: data.usage[key].tokens_saved || 0,
  })).sort((a, b) => (b.Used + b.Saved) - (a.Used + a.Saved));

  const latencyKeys = Object.keys(data.latencies || {});
  const latencyData = latencyKeys.map(key => ({
    name: key,
    LatencyMs: data.latencies[key] || 0,
  })).sort((a, b) => a.LatencyMs - b.LatencyMs);

  const totalTokensSaved = usageKeys.reduce((acc, k) => acc + (data.usage[k].tokens_saved || 0), 0);
  // Estimate cost savings: ~$3.00 per 1M tokens approx
  const estimatedSavings = (totalTokensSaved / 1000000) * 3.00;

  return (
    <div style={{ padding: '2rem', height: '100%', overflowY: 'auto', background: 'var(--bg-primary)' }}>
      <div style={{ marginBottom: '2rem' }}>
        <h2>Performance & Analytics</h2>
        <p style={{ color: 'var(--text-muted)' }}>Monitor your token usage, latency, and cache savings across all providers.</p>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '1.5rem', marginBottom: '2.5rem' }}>
        <div style={{ background: 'var(--bg-secondary)', padding: '1.5rem', borderRadius: '12px', border: '1px solid var(--border-color)' }}>
          <div style={{ fontSize: '0.85rem', color: 'var(--text-muted)', textTransform: 'uppercase', letterSpacing: '0.05em', marginBottom: '0.5rem' }}>Tokens Saved via Cache</div>
          <div style={{ fontSize: '2rem', fontWeight: 700, color: 'var(--accent-primary)' }}>{totalTokensSaved.toLocaleString()}</div>
        </div>
        <div style={{ background: 'var(--bg-secondary)', padding: '1.5rem', borderRadius: '12px', border: '1px solid var(--border-color)' }}>
          <div style={{ fontSize: '0.85rem', color: 'var(--text-muted)', textTransform: 'uppercase', letterSpacing: '0.05em', marginBottom: '0.5rem' }}>Estimated Cost Savings</div>
          <div style={{ fontSize: '2rem', fontWeight: 700, color: '#27ae60' }}>${estimatedSavings.toFixed(2)}</div>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '2rem', marginBottom: '2rem' }}>
        <div style={{ background: 'var(--bg-secondary)', padding: '1.5rem', borderRadius: '12px', border: '1px solid var(--border-color)', height: '350px' }}>
          <h3 style={{ marginBottom: '1rem', fontSize: '1.1rem' }}>Token Usage vs Savings</h3>
          <ResponsiveContainer width="100%" height="90%">
            <BarChart data={tokenData} margin={{ top: 10, right: 10, left: 0, bottom: 20 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="var(--border-color)" vertical={false} />
              <XAxis dataKey="name" tick={{fill: 'var(--text-muted)'}} axisLine={false} tickLine={false} />
              <YAxis tick={{fill: 'var(--text-muted)'}} axisLine={false} tickLine={false} />
              <Tooltip contentStyle={{background: 'var(--bg-primary)', border: '1px solid var(--border-color)', borderRadius: '8px'}} />
              <Legend />
              <Bar dataKey="Used" stackId="a" fill="var(--text-secondary)" radius={[0, 0, 4, 4]} />
              <Bar dataKey="Saved" stackId="a" fill="var(--accent-primary)" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div style={{ background: 'var(--bg-secondary)', padding: '1.5rem', borderRadius: '12px', border: '1px solid var(--border-color)', height: '350px' }}>
          <h3 style={{ marginBottom: '1rem', fontSize: '1.1rem' }}>Moving Average Latency (ms)</h3>
          <ResponsiveContainer width="100%" height="90%">
            <BarChart data={latencyData} margin={{ top: 10, right: 10, left: 0, bottom: 20 }} layout="vertical">
              <CartesianGrid strokeDasharray="3 3" stroke="var(--border-color)" horizontal={false} />
              <XAxis type="number" tick={{fill: 'var(--text-muted)'}} axisLine={false} tickLine={false} />
              <YAxis type="category" dataKey="name" tick={{fill: 'var(--text-muted)'}} axisLine={false} tickLine={false} width={80} />
              <Tooltip contentStyle={{background: 'var(--bg-primary)', border: '1px solid var(--border-color)', borderRadius: '8px'}} cursor={{fill: 'var(--bg-tertiary)'}} />
              <Bar dataKey="LatencyMs" fill="#3498db" radius={[0, 4, 4, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

    </div>
  );
};
