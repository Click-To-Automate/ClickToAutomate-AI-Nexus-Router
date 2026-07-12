import { useState } from 'react';

interface OnboardingProps {
  onComplete: () => void;
}

export function Onboarding({ onComplete }: OnboardingProps) {
  const [email, setEmail] = useState('');
  const [subscribe, setSubscribe] = useState(true);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return;
    
    // In a real app, send to webhook here
    console.log("Captured Email:", email, "Subscribe:", subscribe);
    
    localStorage.setItem('hasOnboarded', 'true');
    onComplete();
  };

  return (
    <div style={{ display: 'flex', height: '100vh', alignItems: 'center', justifyContent: 'center' }}>
      <div className="glass-panel" style={{ width: '100%', maxWidth: '450px' }}>
        <div style={{ textAlign: 'center', marginBottom: '2.5rem' }}>
          <img src="https://www.clicktoautomate.in/icon-192.png" alt="Logo" style={{ width: '64px', height: '64px', marginBottom: '1rem' }} />
          <h1 style={{ fontSize: '2rem', letterSpacing: '-0.03em' }}>AI Nexus Router</h1>
          <p style={{ margin: 0 }}>Local AI Gateway Initialization</p>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="input-group">
            <label className="input-label" htmlFor="email">Email Address</label>
            <input 
              type="email" 
              id="email"
              className="text-input" 
              placeholder="you@example.com" 
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required 
            />
          </div>

          <label className="checkbox-group">
            <input 
              type="checkbox" 
              checked={subscribe}
              onChange={(e) => setSubscribe(e.target.checked)}
            />
            <span className="checkbox-label">Keep me updated on AI Nexus Router releases and features</span>
          </label>

          <label className="checkbox-group" style={{ opacity: 0.5 }}>
            <input type="checkbox" disabled />
            <span className="checkbox-label">Sync telemetry with ClickToAutomate Cloud (Coming Soon)</span>
          </label>

          <button type="submit" className="btn-primary" style={{ marginTop: '1rem' }}>
            Initialize AI Nexus Router
          </button>
        </form>
      </div>
    </div>
  );
}
