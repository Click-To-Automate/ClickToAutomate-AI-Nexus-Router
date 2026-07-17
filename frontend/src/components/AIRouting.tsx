import { useState, useEffect } from 'react';
import { API_BASE } from '../api';

export function AIRouting() {
  const [settings, setSettings] = useState({
    ai_profile: 'balanced',
    preferred_models: '[]'
  });
  const [status, setStatus] = useState('');
  
  // Available connected providers
  const [providers, setProviders] = useState<any[]>([]);
  // Available models from all providers
  const [models, setModels] = useState<any[]>([]);
  // The ordered list of active models defined by user (format: "providerID::modelID")
  const [preferredModels, setPreferredModels] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  // State for expanded provider sections
  const [expandedProviders, setExpandedProviders] = useState<Record<string, boolean>>({});

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const [settingsRes, keysRes, provRes, modRes] = await Promise.all([
        fetch(`${API_BASE}/v1/settings`),
        fetch(`${API_BASE}/v1/keys`),
        fetch(`${API_BASE}/v1/providers`),
        fetch(`${API_BASE}/v1/models`)
      ]);
      
      const sData = await settingsRes.json();
      const keysData = await keysRes.json();
      const provData = await provRes.json();
      const modData = await modRes.json();
      
      const savedProfile = sData.ai_profile || 'balanced';
      const savedPreferred = sData.preferred_models || '[]';
      
      setSettings({
        ai_profile: savedProfile,
        preferred_models: savedPreferred
      });
      
      let parsedPreferred: string[] = [];
      try {
        parsedPreferred = JSON.parse(savedPreferred);
      } catch(e) {}
      
      // Determine which providers actually have keys
      const activeKeys = keysData.keys || [];
      const connectedProviderIds = new Set(activeKeys.map((k: any) => k.provider_id));
      
      const allProviders = provData.providers || [];
      
      // Map providers with basic details
      const connectedProviders = allProviders
        .filter((p: any) => connectedProviderIds.has(p.id))
        .map((p: any) => ({
          id: p.id,
          name: p.name || p.id,
          logo: `https://www.google.com/s2/favicons?domain=${p.base_url.replace(/https?:\/\//, '').split('/')[0]}&sz=64`
        }));
        
      setProviders(connectedProviders);
      setModels(modData.data || []);
      
      // For old configs that just had provider IDs, we filter those out or they will be ignored backend
      setPreferredModels(parsedPreferred);
      
      // Initialize all connected providers as collapsed by default
      const initialExpanded: Record<string, boolean> = {};
      connectedProviders.forEach((p: any) => {
        initialExpanded[p.id] = false;
      });
      setExpandedProviders(initialExpanded);

    } catch (err) {
      console.error('Failed to fetch data', err);
    } finally {
      setLoading(false);
    }
  };

  const saveSettings = async (profile: string, prefModels: string[]) => {
    const newSettings = {
      ai_profile: profile,
      preferred_models: JSON.stringify(prefModels)
    };
    setSettings(newSettings);
    
    try {
      await fetch(`${API_BASE}/v1/settings`, {
        method: 'POST',
        body: JSON.stringify(newSettings),
      });
      setStatus('Routing settings saved!');
      setTimeout(() => setStatus(''), 3000);
    } catch (e) {
      console.error(e);
      setStatus('Error saving settings');
    }
  };

  const moveUp = (index: number) => {
    if (index === 0) return;
    const newOrder = [...preferredModels];
    const temp = newOrder[index - 1];
    newOrder[index - 1] = newOrder[index];
    newOrder[index] = temp;
    setPreferredModels(newOrder);
    saveSettings(settings.ai_profile, newOrder);
  };

  const moveDown = (index: number) => {
    if (index === preferredModels.length - 1) return;
    const newOrder = [...preferredModels];
    const temp = newOrder[index + 1];
    newOrder[index + 1] = newOrder[index];
    newOrder[index] = temp;
    setPreferredModels(newOrder);
    saveSettings(settings.ai_profile, newOrder);
  };
  
  const addModel = (modelKey: string) => {
    if (preferredModels.includes(modelKey)) return;
    const newOrder = [...preferredModels, modelKey];
    setPreferredModels(newOrder);
    saveSettings(settings.ai_profile, newOrder);
  };

  const removeModel = (modelKey: string) => {
    const newOrder = preferredModels.filter(m => m !== modelKey);
    setPreferredModels(newOrder);
    saveSettings(settings.ai_profile, newOrder);
  };
  
  const toggleAutomatic = () => {
    const newOrder = preferredModels.length > 0 ? [] : [];
    setPreferredModels(newOrder);
    saveSettings(settings.ai_profile, newOrder);
  };

  const toggleProviderExpand = (providerId: string) => {
    setExpandedProviders(prev => ({
      ...prev,
      [providerId]: !prev[providerId]
    }));
  };

  const PROFILES = [
    {
      id: 'fast',
      icon: '🚀',
      name: 'Fast Coding',
      desc: 'Prioritize response speed. Uses minimal reasoning and context processing. Ideal for rapid iteration and autocomplete.'
    },
    {
      id: 'balanced',
      icon: '⚖️',
      name: 'Balanced',
      desc: 'Balance response speed, reasoning quality, and token usage. Moderate context compression. Recommended default.'
    },
    {
      id: 'max',
      icon: '🧠',
      name: 'Max Token Optimization',
      desc: 'Prioritize reasoning quality and token usage. Aggressive context and image compression. Accepts slower responses for better accuracy on large codebases.'
    }
  ];

  if (loading) return <div style={{ padding: '2rem' }}>Loading Routing settings...</div>;

  return (
    <div className="fade-in" style={{ padding: '2rem', paddingBottom: '4rem' }}>
      <div style={{ marginBottom: '3rem' }}>
        <h2>AI Routing & Performance</h2>
        <p style={{ color: 'var(--text-muted)' }}>Configure AI performance profiles and control strict fallback priorities for seamless coding.</p>
      </div>

      <div style={{ marginBottom: '3rem' }}>
        <h3 style={{ marginBottom: '1rem', color: 'var(--text-primary)' }}>AI Performance Mode</h3>
        <div style={{ display: 'grid', gap: '1rem', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))' }}>
          {PROFILES.map(prof => (
            <div 
              key={prof.id}
              onClick={() => saveSettings(prof.id, preferredModels)}
              style={{ 
                background: 'var(--bg-primary)', 
                padding: '1.5rem', 
                borderRadius: '12px', 
                border: settings.ai_profile === prof.id ? '2px solid var(--accent-primary)' : '1px solid var(--border-color)', 
                cursor: 'pointer',
                position: 'relative',
                boxShadow: settings.ai_profile === prof.id ? '0 4px 12px var(--accent-glow)' : 'none',
                transition: 'all 0.2s ease'
              }}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '0.75rem' }}>
                <span style={{ fontSize: '1.5rem' }}>{prof.icon}</span>
                <h4 style={{ margin: 0 }}>{prof.name}</h4>
              </div>
              <p style={{ margin: 0, fontSize: '0.85rem', color: 'var(--text-muted)', lineHeight: 1.5 }}>
                {prof.desc}
              </p>
              {settings.ai_profile === prof.id && (
                <div style={{ position: 'absolute', top: '1rem', right: '1rem', color: 'var(--accent-primary)' }}>
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>

      <div>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
          <h3 style={{ margin: 0, color: 'var(--text-primary)' }}>Preferred Models Priority</h3>
          <button 
            className="btn-secondary" 
            onClick={toggleAutomatic}
            style={{ fontSize: '0.8rem', padding: '0.4rem 0.8rem' }}
          >
            {preferredModels.length > 0 ? 'Switch to Auto-Routing (Based on Intent)' : 'Enable Strict Priority Routing'}
          </button>
        </div>
        <p style={{ color: 'var(--text-muted)', fontSize: '0.9rem', marginBottom: '1.5rem' }}>
          Select and order the exact AI models CTA-AI-Nexus is allowed to use. The first model becomes Priority 1. 
          If a provider fails due to expired tokens, limits, or outages, it automatically switches to the next priority model.
        </p>

        {providers.length === 0 ? (
          <div style={{ padding: '2rem', background: 'var(--bg-primary)', borderRadius: '12px', border: '1px solid var(--border-color)', textAlign: 'center' }}>
            <p style={{ color: 'var(--text-muted)' }}>No providers configured. Please add API keys in the Connections tab first.</p>
          </div>
        ) : (
          <div style={{ display: 'flex', gap: '2rem', flexWrap: 'wrap' }}>
            {/* Active Priority List */}
            <div style={{ flex: '1 1 300px' }}>
              <h4 style={{ marginBottom: '1rem', color: 'var(--text-secondary)' }}>Priority Queue</h4>
              {preferredModels.length === 0 ? (
                <div style={{ padding: '2rem', background: 'var(--bg-tertiary)', borderRadius: '12px', border: '1px dashed var(--border-color)', textAlign: 'center' }}>
                  <p style={{ color: 'var(--text-muted)', fontSize: '0.9rem' }}>Strict priority disabled. Router will automatically choose the best model for each task.</p>
                </div>
              ) : (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                  {preferredModels.map((modelKey, index) => {
                    const parts = modelKey.split('::');
                    const provId = parts[0];
                    const modId = parts[1] || parts[0]; // fallback for old configs
                    
                    const p = providers.find(prov => prov.id === provId);
                    
                    return (
                      <div key={modelKey} className="priority-list-item">
                        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                          <div style={{ 
                            width: '24px', height: '24px', borderRadius: '50%', background: 'var(--accent-primary)', 
                            color: '#fff', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '0.8rem', fontWeight: 600, flexShrink: 0
                          }}>
                            {index + 1}
                          </div>
                          {p && p.logo ? (
                            <img src={p.logo} alt={p.name} width="24" height="24" style={{ borderRadius: '4px', background: '#fff', flexShrink: 0 }} />
                          ) : (
                            <div style={{ width: '24px', height: '24px', borderRadius: '4px', background: 'var(--border-color)', flexShrink: 0 }}></div>
                          )}
                          <div style={{ display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
                            <span style={{ fontWeight: 600, fontSize: '0.85rem', color: 'var(--text-primary)', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{modId}</span>
                            <span style={{ fontSize: '0.7rem', color: 'var(--text-muted)' }}>{p ? p.name : provId}</span>
                          </div>
                        </div>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', flexShrink: 0 }}>
                          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.1rem' }}>
                            <button onClick={() => moveUp(index)} disabled={index === 0} className="icon-btn-small">▲</button>
                            <button onClick={() => moveDown(index)} disabled={index === preferredModels.length - 1} className="icon-btn-small">▼</button>
                          </div>
                          <button onClick={() => removeModel(modelKey)} className="icon-btn-small" style={{ color: '#ef4444', marginLeft: '0.5rem', padding: '0.3rem' }}>✕</button>
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </div>

            {/* Available Models Grouped by Provider */}
            <div style={{ flex: '1 1 300px' }}>
              <h4 style={{ marginBottom: '1rem', color: 'var(--text-secondary)' }}>Available Models</h4>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                {providers.map(p => {
                  // Find models for this provider and sort alphabetically
                  const providerModels = models.filter(m => m.owned_by === p.id).sort((a, b) => a.id.localeCompare(b.id));
                  const isExpanded = expandedProviders[p.id];
                  
                  return (
                    <div key={p.id} style={{ background: 'var(--bg-primary)', border: '1px solid var(--border-color)', borderRadius: '8px', overflow: 'hidden' }}>
                      {/* Provider Header */}
                      <div 
                        onClick={() => toggleProviderExpand(p.id)}
                        style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '0.8rem 1rem', cursor: 'pointer', background: isExpanded ? 'var(--bg-tertiary)' : 'transparent' }}
                      >
                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                          {p.logo ? (
                            <img src={p.logo} alt={p.name} width="20" height="20" style={{ borderRadius: '4px', background: '#fff' }} />
                          ) : null}
                          <span style={{ fontWeight: 600 }}>{p.name}</span>
                          <span style={{ fontSize: '0.75rem', background: 'var(--bg-secondary)', padding: '0.1rem 0.4rem', borderRadius: '4px', color: 'var(--text-muted)' }}>
                            {providerModels.length} models
                          </span>
                        </div>
                        <span style={{ color: 'var(--text-muted)' }}>{isExpanded ? '▼' : '▶'}</span>
                      </div>
                      
                      {/* Models List */}
                      {isExpanded && (
                        <div style={{ borderTop: '1px solid var(--border-color)', maxHeight: '300px', overflowY: 'auto' }}>
                          {providerModels.length === 0 ? (
                            <div style={{ padding: '1rem', textAlign: 'center', color: 'var(--text-muted)', fontSize: '0.85rem' }}>No models discovered yet.</div>
                          ) : (
                            providerModels.map(m => {
                              const modelKey = `${p.id}::${m.id}`;
                              const isSelected = preferredModels.includes(modelKey);
                              
                              if (isSelected) return null;
                              
                              return (
                                <div 
                                  key={modelKey} 
                                  onClick={() => addModel(modelKey)} 
                                  className="available-item" 
                                  style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '0.6rem 1rem', paddingLeft: '2.5rem', cursor: 'pointer', borderBottom: '1px solid var(--border-color)' }}
                                >
                                  <span style={{ fontSize: '0.85rem' }}>{m.id}</span>
                                  <span style={{ color: 'var(--accent-primary)', fontSize: '1rem', fontWeight: 600 }}>+</span>
                                </div>
                              );
                            })
                          )}
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
        )}
      </div>
      
      {status && (
        <div style={{ position: 'fixed', bottom: '2rem', right: '2rem', background: 'var(--accent-primary)', color: '#fff', padding: '0.75rem 1.5rem', borderRadius: '8px', boxShadow: '0 4px 12px rgba(0,0,0,0.15)', zIndex: 100 }} className="fade-in">
          {status}
        </div>
      )}
    </div>
  );
}
