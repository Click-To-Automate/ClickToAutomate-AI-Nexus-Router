import { useState, useEffect } from 'react';

export function Settings() {
  const [settings, setSettings] = useState({
    port: '20128',
    dark_mode: 'true',
  });
  const [status, setStatus] = useState('');

  useEffect(() => {
    fetch('http://localhost:20128/v1/settings')
      .then(r => r.json())
      .then(d => setSettings({ port: d.port || '20128', dark_mode: d.dark_mode || 'true' }))
      .catch(console.error);
  }, []);

  const saveSettings = () => {
    fetch('http://localhost:20128/v1/settings', {
      method: 'POST',
      body: JSON.stringify(settings),
    }).then(() => {
      setStatus('Saved! Note: Changing port requires an app restart.');
      setTimeout(() => setStatus(''), 3000);
    });
  };

  return (
    <div className="fade-in" style={{ padding: '2rem' }}>
      <div style={{ marginBottom: '3rem' }}>
        <h2>System Settings</h2>
        <p style={{ color: 'var(--text-muted)' }}>Configure the core behavior of your AI Nexus Router proxy.</p>
      </div>

      <div style={{ background: 'var(--bg-primary)', padding: '2rem', borderRadius: '12px', border: '1px solid var(--border-color)', maxWidth: '600px' }}>
        <div style={{ marginBottom: '1.5rem' }}>
          <label style={{ display: 'block', marginBottom: '0.5rem', fontWeight: 600 }}>Local Proxy Port</label>
          <input 
            type="text" 
            className="text-input" 
            value={settings.port} 
            onChange={e => setSettings({...settings, port: e.target.value})} 
            style={{ width: '100%' }}
          />
          <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', marginTop: '0.5rem' }}>
            The port this app listens on (default: 20128). Your IDEs (Cursor, Windsurf) must point to http://localhost:[port]/v1
          </p>
        </div>

        <div style={{ marginBottom: '2rem' }}>
          <label style={{ display: 'block', marginBottom: '0.5rem', fontWeight: 600 }}>Theme</label>
          <select 
            className="text-input" 
            value={settings.dark_mode} 
            onChange={e => setSettings({...settings, dark_mode: e.target.value})}
            style={{ width: '100%' }}
          >
            <option value="true">Dark Mode</option>
            <option value="false">Light Mode (Coming Soon)</option>
          </select>
        </div>

        <button className="btn-primary" onClick={saveSettings}>Save Settings</button>
        {status && <span style={{ marginLeft: '1rem', color: '#10a37f' }}>{status}</span>}
      </div>
    </div>
  );
}
