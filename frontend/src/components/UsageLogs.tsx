import { useState, useEffect } from 'react';
import { API_BASE } from '../api';
import { Analytics } from './Analytics';

type UsageMap = {
  [providerId: string]: { count: number; tokens_saved: number; tokens_used: number };
};

// Simple mapping for nice display names and icons (removed emojis per request)
const PROVIDER_META: Record<string, { name: string; icon: string; color: string }> = {
  'mistral': { name: 'Mistral', icon: '', color: '#f39c12' },
  'cerebras': { name: 'Cerebras', icon: '', color: '#e74c3c' },
  'groq': { name: 'Groq', icon: '', color: '#2ecc71' },
  'gemini': { name: 'Google Gemini', icon: '', color: '#3498db' },
  'openai': { name: 'OpenAI', icon: '', color: '#10a37f' },
  'anthropic': { name: 'Anthropic', icon: '', color: '#d35400' },
  'openrouter': { name: 'OpenRouter', icon: '', color: '#9b59b6' },
  'together': { name: 'Together AI', icon: '', color: '#34495e' },
  'ollama': { name: 'Ollama', icon: '', color: '#bdc3c7' },
  'github-models': { name: 'GitHub Models', icon: '', color: '#111827' },
  'sambanova': { name: 'SambaNova', icon: '', color: '#e67e22' },
  'deepseek': { name: 'DeepSeek', icon: '', color: '#2980b9' },
};

export function UsageLogs() {
  const [usage, setUsage] = useState<UsageMap>({});
  const [savedKeys, setSavedKeys] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [viewMode, setViewMode] = useState<'grid' | 'list' | 'charts'>('grid');

  const fetchData = async () => {
    try {
      const [usageRes, keysRes] = await Promise.all([
        fetch(`${API_BASE}/v1/usage`),
        fetch(`${API_BASE}/v1/keys`)
      ]);
      
      if (usageRes.ok) {
        const uData = await usageRes.json();
        setUsage(uData);
      }
      if (keysRes.ok) {
        const kData = await keysRes.json();
        setSavedKeys(kData.keys || []);
      }
    } catch (error) {
      console.error('Failed to fetch telemetry data', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, []);

  const totalRequests = Object.values(usage).reduce((sum, data) => sum + (data.count || 0), 0);
  const totalSavedTokens = Object.values(usage).reduce((sum, data) => sum + (data.tokens_saved || 0), 0);
  const totalUsedTokens = Object.values(usage).reduce((sum, data) => sum + (data.tokens_used || 0), 0);

  // Determine all "active" providers (those with usage OR configured keys)
  const activeProviderIds = new Set(Object.keys(usage));
  savedKeys.forEach(k => activeProviderIds.add(k.provider_id));

  const activeProvidersList = Array.from(activeProviderIds).map(id => ({
    id,
    count: usage[id]?.count || 0,
    saved: usage[id]?.tokens_saved || 0,
    used: usage[id]?.tokens_used || 0
  })).sort((a, b) => b.count - a.count);

  return (
    <div className="fade-in">
      <div style={{ marginBottom: '2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div>
          <h1 style={{ fontSize: '1.5rem', margin: '0 0 0.5rem 0', fontFamily: "'Inter Tight', sans-serif" }}>Usage Telemetry</h1>
          <p style={{ color: 'var(--text-muted)', margin: 0 }}>
            Live breakdown of where your requests are being routed.
          </p>
        </div>
        <div className="view-toggle">
          <button className={`view-toggle-btn ${viewMode === 'grid' ? 'active' : ''}`} onClick={() => setViewMode('grid')}>Grid</button>
          <button className={`view-toggle-btn ${viewMode === 'list' ? 'active' : ''}`} onClick={() => setViewMode('list')}>List</button>
          <button className={`view-toggle-btn ${viewMode === 'charts' ? 'active' : ''}`} onClick={() => setViewMode('charts')}>Charts</button>
        </div>
      </div>

      <div style={{
        background: 'var(--bg-primary)',
        border: '1px solid var(--border-color)',
        borderRadius: '8px',
        padding: '1.5rem',
        marginBottom: '2rem',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        boxShadow: '0 1px 3px rgba(0,0,0,0.02)'
      }}>
        <div style={{ display: 'flex', gap: '3rem' }}>
          <div>
            <div style={{ color: 'var(--text-muted)', fontSize: '0.9rem', marginBottom: '0.25rem', fontWeight: 600 }}>TOTAL REQUESTS ROUTED</div>
            <div style={{ fontSize: '2.5rem', fontWeight: 'bold', color: 'var(--text-primary)', fontFamily: "'Inter Tight', sans-serif", letterSpacing: '-0.02em' }}>{totalRequests}</div>
          </div>
          <div style={{ paddingLeft: '3rem', borderLeft: '1px solid var(--border-color)' }}>
            <div style={{ color: 'var(--text-muted)', fontSize: '0.9rem', marginBottom: '0.25rem', fontWeight: 600 }}>TOTAL TOKENS USED</div>
            <div style={{ fontSize: '2.5rem', fontWeight: 'bold', color: 'var(--text-primary)', fontFamily: "'Inter Tight', sans-serif", letterSpacing: '-0.02em' }}>{totalUsedTokens.toLocaleString()}</div>
          </div>
          <div style={{ paddingLeft: '3rem', borderLeft: '1px solid var(--border-color)' }}>
            <div style={{ color: 'var(--accent)', fontSize: '0.9rem', marginBottom: '0.25rem', fontWeight: 600 }}>TOTAL TOKENS SAVED</div>
            <div style={{ fontSize: '2.5rem', fontWeight: 'bold', color: 'var(--accent)', fontFamily: "'Inter Tight', sans-serif", letterSpacing: '-0.02em' }}>{totalSavedTokens.toLocaleString()}</div>
          </div>
        </div>
      </div>

      {loading && activeProvidersList.length === 0 ? (
        <div style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '2rem' }}>Loading telemetry...</div>
      ) : activeProvidersList.length === 0 ? (
        <div style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '2rem', background: 'var(--bg-primary)', borderRadius: '12px', border: '1px solid var(--border-color)' }}>
          No keys configured and no requests routed yet.
        </div>
      ) : viewMode === 'charts' ? (
        <div style={{ marginTop: '1rem', borderTop: '1px solid var(--border-color)', paddingTop: '1rem' }}>
          <Analytics />
        </div>
      ) : (
        <div className={viewMode === 'grid' ? 'provider-grid' : 'provider-list'}>
          {activeProvidersList.map(({ id: providerId, count, saved, used }) => {
              const meta = PROVIDER_META[providerId] || { name: providerId, icon: '', color: '#888' };
              const percentage = totalRequests > 0 ? Math.round((count / totalRequests) * 100) : 0;
              const isGrid = viewMode === 'grid';
              
              return (
                <div key={providerId} className={isGrid ? 'provider-card-v2' : 'provider-card-list'}>
                  <div className="provider-card-top-bar" style={{ background: meta.color }}></div>
                  
                  <div className="provider-card-header">
                    <div className="provider-card-avatar">
                      <img 
                        src={`/providers/${providerId}.png`} 
                        alt={meta.name} 
                        onError={(e) => {
                          e.currentTarget.onerror = null;
                          e.currentTarget.src = `data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><rect width='100' height='100' fill='${encodeURIComponent(meta.color)}' rx='18'/><text x='50' y='62' font-family='Arial' font-size='46' font-weight='bold' fill='white' text-anchor='middle'>${meta.name.charAt(0).toUpperCase()}</text></svg>`;
                        }}
                      />
                    </div>
                    {isGrid && (
                      <span className="badge" style={{ background: 'var(--bg-tertiary)', color: 'var(--text-secondary)' }}>
                        #{percentage}%
                      </span>
                    )}
                  </div>

                  <div className="provider-card-body">
                    <h3 className="provider-card-title">{meta.name}</h3>
                    <p className="provider-card-subtitle" style={{ color: 'var(--text-primary)' }}>
                      <strong style={{ fontSize: isGrid ? '1.75rem' : '1.25rem', color: 'var(--text-primary)' }}>{count}</strong> requests
                    </p>
                    <p style={{ margin: '0.5rem 0 0 0', fontSize: '0.85rem', color: 'var(--text-primary)', fontWeight: 600 }}>
                      {used.toLocaleString()} tokens used
                    </p>
                    <p style={{ margin: '0.25rem 0 0 0', fontSize: '0.85rem', color: 'var(--accent)', fontWeight: 600 }}>
                      {saved.toLocaleString()} tokens saved
                    </p>
                    
                    {isGrid && (
                      <div style={{ width: '100%', height: '6px', background: 'var(--bg-tertiary)', borderRadius: '3px', overflow: 'hidden', marginTop: '1rem' }}>
                        <div style={{ width: `${percentage}%`, height: '100%', background: meta.color, borderRadius: '3px', transition: 'width 0.5s ease-out' }}></div>
                      </div>
                    )}
                  </div>

                  <div className="provider-card-footer">
                    {!isGrid && (
                      <div className="stat-block" style={{ textAlign: 'left', alignItems: 'flex-start' }}>
                        <div style={{ width: '100px', height: '6px', background: 'var(--bg-tertiary)', borderRadius: '3px', overflow: 'hidden' }}>
                          <div style={{ width: `${percentage}%`, height: '100%', background: meta.color, borderRadius: '3px', transition: 'width 0.5s ease-out' }}></div>
                        </div>
                      </div>
                    )}
                    <div className="stat-block" style={{ gridColumn: isGrid ? 'span 3' : 'auto', textAlign: isGrid ? 'center' : 'right' }}>
                      <span className="stat-label">Share of Traffic</span>
                      <span className="stat-value">{percentage}%</span>
                    </div>
                  </div>
                </div>
              );
          })}
        </div>
      )}
    </div>
  );
}
