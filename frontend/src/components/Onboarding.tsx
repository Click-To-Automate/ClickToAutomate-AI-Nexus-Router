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
        <div style={{ textAlign: 'center', marginBottom: '2rem' }}>
          <div style={{ fontSize: '3rem', marginBottom: '0.5rem' }}>🌌</div>
          <h1>ClickToAutomate</h1>
          <p>AI Nexus Router Initialization</p>
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
            <span className="checkbox-label">Keep me updated with AI news and exclusive offers</span>
          </label>

          <label className="checkbox-group" style={{ opacity: 0.5 }}>
            <input type="checkbox" disabled />
            <span className="checkbox-label">Sync with Web Dashboard (Coming Soon)</span>
          </label>

          <button type="submit" className="btn-primary" style={{ marginTop: '1rem' }}>
            Initialize Router
          </button>
        </form>
      </div>
    </div>
  );
}
