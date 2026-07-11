import { useState, useEffect } from 'react';

interface ProviderKey {
  provider_id: string;
  masked_key: string;
  is_set: string;
}

const ALL_PROVIDERS = [
  { id: 'openai', name: 'OpenAI', url: 'https://platform.openai.com', freeTier: 'Paid' },
  { id: 'anthropic', name: 'Anthropic', url: 'https://console.anthropic.com', freeTier: 'Paid' },
  { id: 'gemini', name: 'Google Gemini', url: 'https://aistudio.google.com', freeTier: 'Free (60M tokens/mo)' },
  { id: 'groq', name: 'Groq', url: 'https://console.groq.com', freeTier: 'Free (15M tokens/mo)' },
  { id: 'openrouter', name: 'OpenRouter', url: 'https://openrouter.ai', freeTier: 'Free (Limited Models)' },
  { id: 'mistral', name: 'Mistral', url: 'https://console.mistral.ai', freeTier: 'Free (1B tokens/mo)' },
  { id: 'cerebras', name: 'Cerebras', url: 'https://cloud.cerebras.ai', freeTier: 'Free (30M tokens/mo)' },
  { id: 'llm7', name: 'LLM7', url: 'https://llm7.io', freeTier: 'Free (150M tokens/mo)' },
  { id: 'sambanova', name: 'SambaNova', url: 'https://cloud.sambanova.ai', freeTier: '$5 Signup Credit' },
  { id: 'deepseek', name: 'DeepSeek', url: 'https://platform.deepseek.com', freeTier: '5M Signup Credit' },
];

export function ProvidersSettings() {
  const [savedKeys, setSavedKeys] = useState<ProviderKey[]>([]);
  const [loading, setLoading] = useState(true);
  const [inputs, setInputs] = useState<Record<string, string>>({});

  useEffect(() => {
    fetchKeys();
  }, []);

  const fetchKeys = async () => {
    try {
      const res = await fetch('http://localhost:20128/v1/keys');
      const data = await res.json();
      setSavedKeys(data.keys || []);
    } catch (err) {
      console.error("Failed to fetch keys", err);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (providerId: string) => {
    const keyToSave = inputs[providerId];
    if (!keyToSave) return;

    try {
      await fetch('http://localhost:20128/v1/keys', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ provider_id: providerId, api_key: keyToSave })
      });
      // Clear input and refresh
      setInputs(prev => ({ ...prev, [providerId]: '' }));
      fetchKeys();
    } catch (err) {
      console.error("Failed to save key", err);
    }
  };

  if (loading) return <div>Loading...</div>;

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto' }}>
      <div style={{ marginBottom: '2rem' }}>
        <h2>Provider Configuration</h2>
        <p>Set your API keys for the LLM providers you want to route to. Keys are stored securely in your local SQLite database.</p>
      </div>

      {ALL_PROVIDERS.map(p => {
        const saved = savedKeys.find(k => k.provider_id === p.id);
        
        return (
          <div key={p.id} className="provider-card">
            <div className="provider-header">
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                <span className="provider-title">{p.name}</span>
                <a href={p.url} target="_blank" rel="noreferrer" style={{ fontSize: '0.75rem', color: 'var(--accent-primary)', textDecoration: 'none' }}>
                  Get API Key ↗
                </a>
                <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)', background: 'var(--bg-tertiary)', padding: '0.125rem 0.5rem', borderRadius: '4px' }}>
                  {p.freeTier}
                </span>
              </div>
              {saved ? (
                <span className="provider-badge">Configured</span>
              ) : (
                <span className="provider-badge inactive">Not Configured</span>
              )}
            </div>
            
            <div className="provider-actions">
              <div className="input-group">
                <label className="input-label">API Key</label>
                <input 
                  type="password" 
                  className="text-input" 
                  placeholder={saved ? saved.masked_key : `Enter ${p.name} API Key`}
                  value={inputs[p.id] || ''}
                  onChange={(e) => setInputs(prev => ({ ...prev, [p.id]: e.target.value }))}
                />
              </div>
              <button 
                className="btn-secondary" 
                onClick={() => handleSave(p.id)}
                disabled={!inputs[p.id]}
              >
                Save Key
              </button>
            </div>
          </div>
        );
      })}
    </div>
  );
}
