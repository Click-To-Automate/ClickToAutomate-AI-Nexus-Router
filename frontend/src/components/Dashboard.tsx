import { useState } from 'react';
import { ProvidersSettings } from './ProvidersSettings';
import { UsageLogs } from './UsageLogs';

export function Dashboard() {
  const [activeTab, setActiveTab] = useState('providers');

  return (
    <div className="app-container">
      <div className="sidebar">
        <div style={{ padding: '0 1rem', marginBottom: '2rem' }}>
          <div style={{ fontSize: '1.5rem', marginBottom: '0.25rem' }}>🌌</div>
          <h2 style={{ fontSize: '1rem', margin: 0 }}>CTA AI Nexus</h2>
        </div>
        
        <div 
          className={`nav-item ${activeTab === 'providers' ? 'active' : ''}`}
          onClick={() => setActiveTab('providers')}
        >
          Providers
        </div>
        
        <div 
          className={`nav-item ${activeTab === 'prompt-studio' ? 'active' : ''}`}
          onClick={() => setActiveTab('prompt-studio')}
        >
          Prompt Studio
        </div>
        
        <div 
          className={`nav-item ${activeTab === 'logs' ? 'active' : ''}`}
          onClick={() => setActiveTab('logs')}
        >
          Usage & Telemetry
        </div>
      </div>

      <div className="main-content">
        {activeTab === 'providers' && <ProvidersSettings />}
        
        {activeTab === 'prompt-studio' && (
          <div style={{ textAlign: 'center', marginTop: '4rem', color: 'var(--text-muted)' }}>
            <h2>Prompt Studio</h2>
            <p>Coming Soon in Phase 7</p>
          </div>
        )}
        
        {activeTab === 'logs' && <UsageLogs />}
      </div>
    </div>
  );
}
