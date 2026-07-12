import { useState, useEffect } from 'react';

interface ProviderKey {
  provider_id: string;
  masked_key: string;
  api_key: string; // The full key is passed just for deletion matching on backend
}

const HARDCODED_META: Record<string, any> = {
  'openai': { 
    id: 'openai', 
    name: 'OpenAI', 
    url: 'https://platform.openai.com', 
    color: '#10a37f', 
    tier: 'Paid', 
    models: 14, 
    location: 'San Francisco, CA', 
    description: 'Industry-leading language models including GPT-4o and o1 with advanced reasoning capabilities.',
    integration: 'api_key',
    pricing: 'paid',
    account: 'needed'
  },
  'anthropic': { 
    id: 'anthropic', 
    name: 'Anthropic', 
    url: 'https://console.anthropic.com', 
    color: '#d35400', 
    tier: 'Paid', 
    models: 6, 
    location: 'San Francisco, CA', 
    description: 'Creators of the Claude 3.5 model family, prioritizing safety, enormous context windows, and code generation.',
    integration: 'api_key',
    pricing: 'paid',
    account: 'needed'
  },
  'gemini': { 
    id: 'gemini', 
    name: 'Google Gemini', 
    url: 'https://aistudio.google.com', 
    color: '#3498db', 
    tier: 'Free', 
    models: 8, 
    location: 'Mountain View, CA', 
    description: "Google's natively multimodal models with 2-million token context windows and massive free tiers.",
    integration: 'api_key',
    pricing: 'free',
    account: 'needed'
  },
  'groq': { 
    id: 'groq', 
    name: 'Groq', 
    url: 'https://console.groq.com', 
    color: '#2ecc71', 
    tier: 'Free', 
    models: 12, 
    location: 'Mountain View, CA', 
    description: 'Lightning-fast LPU inference engine delivering hundreds of tokens per second for open-source models.',
    integration: 'api_key',
    pricing: 'free',
    account: 'needed'
  },
  'openrouter': { 
    id: 'openrouter', 
    name: 'OpenRouter', 
    url: 'https://openrouter.ai', 
    color: '#9b59b6', 
    tier: 'Free', 
    models: 250, 
    location: 'San Francisco, CA', 
    description: 'A unified API gateway providing seamless access to hundreds of open-source and proprietary LLMs.',
    integration: 'api_key',
    pricing: 'free',
    account: 'needed'
  },
  'mistral': { 
    id: 'mistral', 
    name: 'Mistral', 
    url: 'https://console.mistral.ai', 
    color: '#f39c12', 
    tier: 'Free', 
    models: 10, 
    location: 'Paris, France', 
    description: "Europe's leading AI lab offering highly efficient, compact, and powerful open-weights models.",
    integration: 'api_key',
    pricing: 'free',
    account: 'needed'
  },
  'cerebras': { 
    id: 'cerebras', 
    name: 'Cerebras', 
    url: 'https://cloud.cerebras.ai', 
    color: '#e74c3c', 
    tier: 'Free', 
    models: 2, 
    location: 'Sunnyvale, CA', 
    description: 'Wafer-scale hardware inference providing unprecedented generation speeds for LLaMA 3 models.',
    integration: 'api_key',
    pricing: 'free',
    account: 'needed'
  },
  'llm7': { 
    id: 'llm7', 
    name: 'LLM7', 
    url: 'https://llm7.io', 
    color: '#171717', 
    tier: 'Free', 
    models: 4, 
    location: 'London, UK', 
    description: 'Enterprise-grade managed inference API for fine-tuned and highly specialized open models.',
    integration: 'api_key',
    pricing: 'free',
    account: 'needed'
  },
  'sambanova': { 
    id: 'sambanova', 
    name: 'SambaNova', 
    url: 'https://cloud.sambanova.ai', 
    color: '#e67e22', 
    tier: 'Free', 
    models: 6, 
    location: 'Palo Alto, CA', 
    description: 'SN40L AI chips powering massive open models like LLaMA 3.1 405B at full precision speeds.',
    integration: 'api_key',
    pricing: 'free',
    account: 'needed'
  },
  'deepseek': { 
    id: 'deepseek', 
    name: 'DeepSeek', 
    url: 'https://platform.deepseek.com', 
    color: '#2980b9', 
    tier: 'Free', 
    models: 4, 
    location: 'Hangzhou, China', 
    description: 'Revolutionary architecture delivering state-of-the-art coding and math reasoning at incredibly low costs.',
    integration: 'api_key',
    pricing: 'free',
    account: 'needed'
  },
  'together': { 
    id: 'together', 
    name: 'Together AI', 
    url: 'https://api.together.xyz', 
    color: '#34495e', 
    tier: 'Paid', 
    models: 100, 
    location: 'San Francisco, CA', 
    description: 'Cloud platform for running and fine-tuning open source models at scale.',
    integration: 'api_key',
    pricing: 'paid',
    account: 'needed'
  },
  'ollama': { 
    id: 'ollama', 
    name: 'Ollama', 
    url: 'http://localhost:11434', 
    color: '#2c3e50', 
    tier: 'Free', 
    models: 'Local', 
    location: 'Local Machine', 
    description: 'Run large language models locally on your own hardware.',
    integration: 'local',
    pricing: 'free',
    account: 'not_needed'
  },
  'github-models': { 
    id: 'github-models', 
    name: 'GitHub Models', 
    url: 'https://github.com/marketplace/models', 
    color: '#24292e', 
    tier: 'Free', 
    models: 15, 
    location: 'Global', 
    description: 'Free interactive playground and API to try top AI models.',
    integration: 'api_key',
    pricing: 'free',
    account: 'needed'
  }
};

export function ProvidersSettings() {
  const [savedKeys, setSavedKeys] = useState<ProviderKey[]>([]);
  const [providersList, setProvidersList] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [searchQuery, setSearchQuery] = useState('');
  const [filterIntegration, setFilterIntegration] = useState('all');
  const [filterPricing, setFilterPricing] = useState('all');
  const [filterAccount, setFilterAccount] = useState('all');
  
  // Modal State
  const [selectedProvider, setSelectedProvider] = useState<string | null>(null);
  const [newKeyInput, setNewKeyInput] = useState('');

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const [keysRes, provRes] = await Promise.all([
        fetch('http://localhost:20128/v1/keys'),
        fetch('http://localhost:20128/v1/providers')
      ]);
      const keysData = await keysRes.json();
      const provData = await provRes.json();
      setSavedKeys(keysData.keys || []);
      
      const dynProviders = (provData.providers || []).map((p: any) => {
        const meta = HARDCODED_META[p.id] || {};
        return {
          id: p.id,
          name: p.name || meta.name || p.id,
          url: meta.url || p.base_url || 'https://example.com',
          color: meta.color || '#' + Math.floor(Math.random()*16777215).toString(16).padEnd(6, '0'),
          tier: meta.tier || 'API Key',
          models: meta.models || 'Dynamic',
          location: meta.location || 'Global API',
          description: meta.description || `Generic API wrapper for ${p.name || p.id} compatible endpoints.`,
          integration: meta.integration || 'api_key',
          pricing: meta.pricing || 'paid',
          account: meta.account || 'needed'
        };
      });
      setProvidersList(dynProviders);
    } catch (err) {
      console.error("Failed to fetch data", err);
    } finally {
      setLoading(false);
    }
  };

  const fetchKeys = async () => {
    try {
      const res = await fetch('http://localhost:20128/v1/keys');
      const data = await res.json();
      setSavedKeys(data.keys || []);
    } catch (err) {
      console.error("Failed to fetch keys", err);
    }
  };

  const handleAddKey = async () => {
    if (!selectedProvider || !newKeyInput) return;

    try {
      await fetch('http://localhost:20128/v1/keys', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          provider_id: selectedProvider, 
          api_key: newKeyInput,
          action: 'add'
        })
      });
      setNewKeyInput('');
      fetchKeys();
    } catch (err) {
      console.error("Failed to save key", err);
    }
  };

  const handleDeleteKey = async (apiKey: string) => {
    if (!selectedProvider) return;

    try {
      await fetch('http://localhost:20128/v1/keys', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          provider_id: selectedProvider, 
          api_key: apiKey,
          action: 'delete'
        })
      });
      fetchKeys();
    } catch (err) {
      console.error("Failed to delete key", err);
    }
  };

  const closeModal = () => {
    setSelectedProvider(null);
    setNewKeyInput('');
  };

  if (loading) return <div style={{ padding: '2rem' }}>Loading providers...</div>;

  const getProviderKeys = (providerId: string) => {
    return savedKeys.filter(k => k.provider_id === providerId);
  };

  const filteredProviders = providersList.filter(p => {
    const matchesSearch = p.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          p.id.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesIntegration = filterIntegration === 'all' || p.integration === filterIntegration;
    const matchesPricing = filterPricing === 'all' || p.pricing === filterPricing;
    const matchesAccount = filterAccount === 'all' || p.account === filterAccount;
    
    return matchesSearch && matchesIntegration && matchesPricing && matchesAccount;
  });

  const activeProvider = providersList.find(p => p.id === selectedProvider);
  const activeKeys = selectedProvider ? getProviderKeys(selectedProvider) : [];

  return (
    <div className="fade-in" style={{ paddingBottom: '2rem' }}>
      <div style={{ marginBottom: '3rem', display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div>
          <h2>Provider Configurations</h2>
          <p style={{ color: 'var(--text-muted)' }}>
            Click a provider to add, manage, and delete API keys. The router automatically load-balances requests across all configured keys for a provider!
          </p>
          <div style={{ marginTop: '1rem', display: 'flex', flexWrap: 'wrap', gap: '1rem', alignItems: 'center' }}>
            <input 
              type="text" 
              className="text-input" 
              placeholder="Search 158+ providers..." 
              value={searchQuery}
              onChange={e => setSearchQuery(e.target.value)}
              style={{ width: '250px' }}
            />
            
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <span style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>Integration:</span>
              <select 
                value={filterIntegration} 
                onChange={e => setFilterIntegration(e.target.value)}
                style={{
                  padding: '0.4rem 0.8rem', background: 'var(--bg-primary)', border: '1px solid var(--border-color)', 
                  borderRadius: '8px', fontSize: '0.85rem', color: 'var(--text-primary)', cursor: 'pointer', outline: 'none'
                }}
              >
                <option value="all">All</option>
                <option value="api_key">API Key</option>
                <option value="oauth">OAuth</option>
                <option value="web_cookie">Web Cookie</option>
                <option value="local">Local</option>
              </select>
            </div>

            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <span style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>Pricing:</span>
              <select 
                value={filterPricing} 
                onChange={e => setFilterPricing(e.target.value)}
                style={{
                  padding: '0.4rem 0.8rem', background: 'var(--bg-primary)', border: '1px solid var(--border-color)', 
                  borderRadius: '8px', fontSize: '0.85rem', color: 'var(--text-primary)', cursor: 'pointer', outline: 'none'
                }}
              >
                <option value="all">All</option>
                <option value="free">Free / Free Tier</option>
                <option value="paid">Paid</option>
              </select>
            </div>

            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <span style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>Account:</span>
              <select 
                value={filterAccount} 
                onChange={e => setFilterAccount(e.target.value)}
                style={{
                  padding: '0.4rem 0.8rem', background: 'var(--bg-primary)', border: '1px solid var(--border-color)', 
                  borderRadius: '8px', fontSize: '0.85rem', color: 'var(--text-primary)', cursor: 'pointer', outline: 'none'
                }}
              >
                <option value="all">All</option>
                <option value="needed">Account Needed</option>
                <option value="not_needed">No Account Needed</option>
              </select>
            </div>
          </div>
        </div>
        <div className="view-toggle">
          <button className={`view-toggle-btn ${viewMode === 'grid' ? 'active' : ''}`} onClick={() => setViewMode('grid')}>Grid</button>
          <button className={`view-toggle-btn ${viewMode === 'list' ? 'active' : ''}`} onClick={() => setViewMode('list')}>List</button>
        </div>
      </div>

      <div className={viewMode === 'grid' ? 'provider-grid' : 'provider-list'}>
        {filteredProviders.map(p => {
          const keysCount = getProviderKeys(p.id).length;
          const isGrid = viewMode === 'grid';
          
          return (
            <div 
              key={p.id} 
              className={isGrid ? "provider-card-v2" : "provider-card-list"}
              onClick={() => setSelectedProvider(p.id)}
              style={{
                borderColor: keysCount > 0 ? 'rgba(34, 197, 94, 0.5)' : 'var(--border-color)',
                boxShadow: keysCount > 0 ? '0 0 15px rgba(34, 197, 94, 0.1), 0 4px 6px -1px rgba(0,0,0,0.02)' : ''
              }}
            >
              <div className="provider-card-top-bar" style={{ background: `linear-gradient(90deg, ${p.color}, transparent)` }} />
              
              <div className="provider-card-header">
                <div className="provider-card-avatar">
                  <img src={`/providers/${p.id}.png`} alt={p.name} onError={(e) => { e.currentTarget.src = 'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="%232d3748"/><text x="50" y="50" font-family="Arial" font-size="40" fill="white" text-anchor="middle" dy=".3em">' + p.name.charAt(0).toUpperCase() + '</text></svg>'; }} />
                </div>
                {isGrid && (
                  keysCount > 0 ? (
                    <span className="badge badge-active">ACTIVE</span>
                  ) : (
                    <span className="badge badge-inactive">NEW</span>
                  )
                )}
              </div>

              <div className="provider-card-body">
                <h3 className="provider-card-title">{p.name}</h3>
                <p className="provider-card-subtitle">
                  <a 
                    href={p.url} 
                    target="_blank" 
                    rel="noreferrer" 
                    style={{ color: 'var(--accent-primary)', textDecoration: 'none' }}
                    onClick={(e) => e.stopPropagation()} // Prevent opening modal when clicking link
                  >
                    Get API Key ↗
                  </a>
                </p>
                <p className="provider-card-desc">{p.description}</p>
              </div>

              <div className="provider-card-footer">
                {!isGrid && (
                  <div className="stat-block" style={{ textAlign: 'left', alignItems: 'flex-start' }}>
                    {keysCount > 0 ? (
                      <span className="badge badge-active">ACTIVE</span>
                    ) : (
                      <span className="badge badge-inactive">NEW</span>
                    )}
                  </div>
                )}
                <div className="stat-block">
                  <span className="stat-label">TIER</span>
                  <span className="stat-value">{p.tier}</span>
                </div>
                <div className="stat-block">
                  <span className="stat-label">MODELS</span>
                  <span className="stat-value">{p.models}</span>
                </div>
                <div className="stat-block">
                  <span className="stat-label">KEYS</span>
                  <span className="stat-value">{keysCount}</span>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Modal */}
      {selectedProvider && activeProvider && (
        <div className="modal-overlay" onClick={closeModal}>
          <div className="modal-content fade-in" onClick={e => e.stopPropagation()}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                <img src={`/providers/${activeProvider.id}.png`} alt={activeProvider.name} onError={(e) => { e.currentTarget.src = 'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><rect width="100" height="100" fill="%232d3748"/><text x="50" y="50" font-family="Arial" font-size="40" fill="white" text-anchor="middle" dy=".3em">' + activeProvider.name.charAt(0).toUpperCase() + '</text></svg>'; }} style={{ width: '32px', height: '32px', borderRadius: '8px' }} />
                <h2 style={{ margin: 0, fontFamily: "'Inter Tight', sans-serif" }}>{activeProvider.name} Config</h2>
              </div>
              <button className="btn-icon" onClick={closeModal} style={{ fontSize: '1.5rem', background: 'transparent', border: 'none', cursor: 'pointer', color: 'var(--text-muted)' }}>
                ×
              </button>
            </div>

            <div style={{ marginBottom: '2rem' }}>
              <h4 style={{ marginBottom: '1rem', color: 'var(--text-secondary)' }}>Active API Keys</h4>
              {activeKeys.length === 0 ? (
                <div style={{ padding: '1.5rem', background: 'var(--bg-tertiary)', borderRadius: '8px', textAlign: 'center', color: 'var(--text-muted)' }}>
                  No keys configured yet.
                </div>
              ) : (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                  {activeKeys.map((k, idx) => (
                    <div key={idx} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '1rem', border: '1px solid var(--border-color)', borderRadius: '8px', background: 'var(--bg-primary)' }}>
                      <span style={{ fontFamily: 'monospace', fontSize: '1.1rem', color: 'var(--text-primary)' }}>{k.masked_key}</span>
                      <button 
                        onClick={() => handleDeleteKey(k.api_key)}
                        style={{ background: 'transparent', border: 'none', color: '#ef4444', cursor: 'pointer', fontWeight: 600, fontSize: '0.85rem' }}
                      >
                        Remove
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div style={{ borderTop: '1px solid var(--border-color)', paddingTop: '2rem' }}>
              <h4 style={{ marginBottom: '1rem', color: 'var(--text-secondary)' }}>Add New Key</h4>
              <div style={{ display: 'flex', gap: '1rem' }}>
                <input 
                  type="password" 
                  className="text-input" 
                  placeholder={`Enter new ${activeProvider.name} API Key`}
                  value={newKeyInput}
                  onChange={(e) => setNewKeyInput(e.target.value)}
                  style={{ flex: 1 }}
                  onKeyDown={(e) => { if(e.key === 'Enter') handleAddKey(); }}
                />
                <button 
                  className="btn-secondary" 
                  onClick={handleAddKey}
                  disabled={!newKeyInput}
                >
                  Add Key
                </button>
              </div>
              <p style={{ marginTop: '0.75rem', fontSize: '0.8rem', color: 'var(--text-muted)' }}>
                <a href={activeProvider.url} target="_blank" rel="noreferrer" style={{ color: 'var(--accent-primary)', textDecoration: 'none' }}>Get API Key ↗</a>
              </p>
            </div>

          </div>
        </div>
      )}
    </div>
  );
}
