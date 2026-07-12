import { useState, useEffect } from 'react';
import { API_BASE } from '../api';

export function Compression() {
  const [settings, setSettings] = useState({
    rtk_engine: 'true',
    caveman_engine: 'false',
  });
  const [status, setStatus] = useState('');

  useEffect(() => {
    fetch(`${API_BASE}/v1/settings`)
      .then(r => r.json())
      .then(d => setSettings({ rtk_engine: d.rtk_engine || 'true', caveman_engine: d.caveman_engine || 'false' }))
      .catch(console.error);
  }, []);

  const saveSettings = () => {
    fetch(`${API_BASE}/v1/settings`, {
      method: 'POST',
      body: JSON.stringify(settings),
    }).then(() => {
      setStatus('Compression config saved!');
      setTimeout(() => setStatus(''), 3000);
    });
  };

  const toggle = (key: 'rtk_engine' | 'caveman_engine') => {
    setSettings(s => ({ ...s, [key]: s[key] === 'true' ? 'false' : 'true' }));
  };

  return (
    <div className="fade-in" style={{ padding: '2rem' }}>
      <div style={{ marginBottom: '3rem' }}>
        <h2>Compression Stack</h2>
        <p style={{ color: 'var(--text-muted)' }}>Configure the dual-engine prompt condensation pipeline to save API costs.</p>
      </div>

      <div style={{ display: 'grid', gap: '1.5rem', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))' }}>
        <div style={{ background: 'var(--bg-primary)', padding: '2rem', borderRadius: '12px', border: '1px solid var(--border-color)', display: 'flex', flexDirection: 'column' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
            <h3 style={{ margin: 0, color: '#3498db' }}>RTK Engine</h3>
            <div 
              style={{
                width: '48px', height: '24px', borderRadius: '12px', background: settings.rtk_engine === 'true' ? '#10a37f' : '#cbd5e1',
                position: 'relative', cursor: 'pointer', transition: 'background 0.3s'
              }}
              onClick={() => toggle('rtk_engine')}
            >
              <div style={{
                width: '20px', height: '20px', borderRadius: '50%', background: '#fff',
                position: 'absolute', top: '2px', left: settings.rtk_engine === 'true' ? '26px' : '2px',
                transition: 'left 0.3s cubic-bezier(0.4, 0.0, 0.2, 1)', boxShadow: '0 1px 3px rgba(0,0,0,0.2)'
              }} />
            </div>
          </div>
          <p style={{ color: 'var(--text-muted)', fontSize: '0.9rem', flex: 1 }}>
            Automatically strips redundant terminal logs, traceback bloat, and noisy machine code from your IDE prompts. Reduces context by ~15-40% without losing intent.
          </p>
        </div>

        <div style={{ background: 'var(--bg-primary)', padding: '2rem', borderRadius: '12px', border: '1px solid var(--border-color)', display: 'flex', flexDirection: 'column' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
            <h3 style={{ margin: 0, color: '#e74c3c' }}>Caveman Engine</h3>
            <div 
              style={{
                width: '48px', height: '24px', borderRadius: '12px', background: settings.caveman_engine === 'true' ? '#e74c3c' : '#cbd5e1',
                position: 'relative', cursor: 'pointer', transition: 'background 0.3s'
              }}
              onClick={() => toggle('caveman_engine')}
            >
              <div style={{
                width: '20px', height: '20px', borderRadius: '50%', background: '#fff',
                position: 'absolute', top: '2px', left: settings.caveman_engine === 'true' ? '26px' : '2px',
                transition: 'left 0.3s cubic-bezier(0.4, 0.0, 0.2, 1)', boxShadow: '0 1px 3px rgba(0,0,0,0.2)'
              }} />
            </div>
          </div>
          <p style={{ color: 'var(--text-muted)', fontSize: '0.9rem', flex: 1 }}>
            Advanced NLP compression that aggressively condenses human conversational prose into semantic shorthand. Highly effective but can slightly alter tone. Reduces context by up to 90%.
          </p>
        </div>
      </div>

      <div style={{ marginTop: '2rem' }}>
        <button className="btn-primary" onClick={saveSettings}>Save Pipeline</button>
        {status && <span style={{ marginLeft: '1rem', color: '#10a37f' }}>{status}</span>}
      </div>
    </div>
  );
}
