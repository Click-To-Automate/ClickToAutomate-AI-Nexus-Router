import { useState, useEffect } from 'react';

type UsageMap = {
  [providerId: string]: number;
};

// Simple mapping for nice display names and icons
const PROVIDER_META: Record<string, { name: string; icon: string; color: string }> = {
  'mistral': { name: 'Mistral', icon: '🌪️', color: '#f39c12' },
  'cerebras': { name: 'Cerebras', icon: '🧠', color: '#e74c3c' },
  'groq': { name: 'Groq', icon: '⚡', color: '#2ecc71' },
  'gemini': { name: 'Google Gemini', icon: '✨', color: '#3498db' },
  'openai': { name: 'OpenAI', icon: '🟢', color: '#10a37f' },
  'anthropic': { name: 'Anthropic', icon: '🎎', color: '#d35400' },
  'openrouter': { name: 'OpenRouter', icon: '🌌', color: '#9b59b6' },
  'together': { name: 'Together AI', icon: '🤝', color: '#34495e' },
  'ollama': { name: 'Ollama', icon: '🦙', color: '#bdc3c7' },
  'github-models': { name: 'GitHub Models', icon: '🐙', color: '#ffffff' },
  'sambanova': { name: 'SambaNova', icon: '💃', color: '#e67e22' },
  'deepseek': { name: 'DeepSeek', icon: '🐋', color: '#2980b9' },
};

export function UsageLogs() {
  const [usage, setUsage] = useState<UsageMap>({});
  const [loading, setLoading] = useState(true);

  const fetchUsage = async () => {
    try {
      const response = await fetch('http://localhost:20128/v1/usage');
      if (response.ok) {
        const data = await response.json();
        setUsage(data);
      }
    } catch (error) {
      console.error('Failed to fetch usage', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsage();
    // Poll every 5 seconds for live updates
    const interval = setInterval(fetchUsage, 5000);
    return () => clearInterval(interval);
  }, []);

  const totalRequests = Object.values(usage).reduce((sum, count) => sum + count, 0);

  return (
    <div className="fade-in">
      <div style={{ marginBottom: '2rem' }}>
        <h1 style={{ fontSize: '1.5rem', marginBottom: '0.5rem' }}>Usage Telemetry</h1>
        <p style={{ color: 'var(--text-muted)' }}>
          Live breakdown of where your requests are being routed.
        </p>
      </div>

      <div style={{
        background: 'var(--surface)',
        border: '1px solid var(--border)',
        borderRadius: '12px',
        padding: '1.5rem',
        marginBottom: '2rem',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
      }}>
        <div>
          <div style={{ color: 'var(--text-muted)', fontSize: '0.9rem', marginBottom: '0.25rem' }}>Total Requests Routed</div>
          <div style={{ fontSize: '2rem', fontWeight: 'bold', color: 'var(--primary)' }}>{totalRequests}</div>
        </div>
        <div style={{ fontSize: '2.5rem' }}>📊</div>
      </div>

      {loading && Object.keys(usage).length === 0 ? (
        <div style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '2rem' }}>Loading telemetry...</div>
      ) : Object.keys(usage).length === 0 ? (
        <div style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '2rem', background: 'var(--surface)', borderRadius: '12px', border: '1px solid var(--border)' }}>
          No requests have been routed yet. Try sending a message!
        </div>
      ) : (
        <div style={{ display: 'grid', gap: '1rem', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))' }}>
          {Object.entries(usage)
            .sort(([, a], [, b]) => b - a)
            .map(([providerId, count]) => {
              const meta = PROVIDER_META[providerId] || { name: providerId, icon: '🔌', color: '#888' };
              const percentage = totalRequests > 0 ? Math.round((count / totalRequests) * 100) : 0;
              
              return (
                <div key={providerId} style={{
                  background: 'var(--surface)',
                  border: '1px solid var(--border)',
                  borderRadius: '12px',
                  padding: '1.5rem',
                  display: 'flex',
                  flexDirection: 'column',
                  gap: '1rem',
                  position: 'relative',
                  overflow: 'hidden'
                }}>
                  {/* Subtle top border color matching provider */}
                  <div style={{ position: 'absolute', top: 0, left: 0, right: 0, height: '4px', background: meta.color }}></div>
                  
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                      <span style={{ fontSize: '1.5rem' }}>{meta.icon}</span>
                      <span style={{ fontWeight: '500' }}>{meta.name}</span>
                    </div>
                    <span style={{ fontSize: '1.25rem', fontWeight: 'bold', color: 'var(--text)' }}>{count}</span>
                  </div>

                  <div style={{ width: '100%', height: '8px', background: 'var(--bg)', borderRadius: '4px', overflow: 'hidden' }}>
                    <div style={{ width: `${percentage}%`, height: '100%', background: meta.color, borderRadius: '4px' }}></div>
                  </div>
                  
                  <div style={{ display: 'flex', justifyContent: 'flex-end', fontSize: '0.8rem', color: 'var(--text-muted)' }}>
                    {percentage}% of traffic
                  </div>
                </div>
              );
          })}
        </div>
      )}
    </div>
  );
}
