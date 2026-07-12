import { useState } from 'react';
import { ProvidersSettings } from './ProvidersSettings';
import { UsageLogs } from './UsageLogs';

export function Dashboard() {
  const [activeTab, setActiveTab] = useState('providers');
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);

  return (
    <div className="app-container">
      <div className={`sidebar ${isSidebarCollapsed ? 'collapsed' : ''}`}>
        <div className="brand-header">
          <img src="https://www.clicktoautomate.in/icon-192.png" alt="Logo" className="brand-logo" />
          <h2 className="brand-title">AI Nexus Router</h2>
          <button 
            className="sidebar-toggle-btn"
            onClick={() => setIsSidebarCollapsed(!isSidebarCollapsed)}
          >
            <svg viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round" width="18" height="18" stroke="currentColor" fill="none" strokeWidth="2">
              <polyline points={isSidebarCollapsed ? "13 17 18 12 13 7" : "11 17 6 12 11 7"}></polyline>
              <polyline points={isSidebarCollapsed ? "6 17 11 12 6 7" : "18 17 13 12 18 7"}></polyline>
            </svg>
          </button>
        </div>
        
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
          <div 
            className={`nav-item ${activeTab === 'providers' ? 'active' : ''}`}
            onClick={() => setActiveTab('providers')}
          >
            <svg viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round">
              <rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect>
              <rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect>
              <line x1="6" y1="6" x2="6.01" y2="6"></line>
              <line x1="6" y1="18" x2="6.01" y2="18"></line>
            </svg>
            <span className="nav-text">Provider Settings</span>
          </div>
          
          <div 
            className={`nav-item ${activeTab === 'prompt-studio' ? 'active' : ''}`}
            onClick={() => setActiveTab('prompt-studio')}
          >
            <svg viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round">
              <polyline points="16 18 22 12 16 6"></polyline>
              <polyline points="8 6 2 12 8 18"></polyline>
            </svg>
            <span className="nav-text">Prompt Studio</span>
          </div>
          
          <div 
            className={`nav-item ${activeTab === 'logs' ? 'active' : ''}`}
            onClick={() => setActiveTab('logs')}
          >
            <svg viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round">
              <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline>
            </svg>
            <span className="nav-text">Usage & Telemetry</span>
          </div>
        </div>

        <div style={{ marginTop: 'auto', paddingTop: '2rem' }}>
          <button 
            className="nav-item" 
            style={{ width: '100%', textAlign: 'left', background: 'transparent', border: 'none', color: 'var(--text-secondary)' }}
            onClick={() => {
              localStorage.removeItem('hasOnboarded');
              window.location.reload();
            }}
          >
            <svg viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
              <polyline points="16 17 21 12 16 7"></polyline>
              <line x1="21" y1="12" x2="9" y2="12"></line>
            </svg>
            <span className="nav-text">Log Out</span>
          </button>
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
