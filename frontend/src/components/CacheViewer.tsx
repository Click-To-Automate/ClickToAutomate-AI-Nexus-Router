import { useState, useEffect } from 'react';

export function CacheViewer() {
  const [cacheEntries, setCacheEntries] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchCache = () => {
    fetch('http://localhost:20128/v1/cache')
      .then(r => r.json())
      .then(d => {
        setCacheEntries(d.cache || []);
        setLoading(false);
      })
      .catch(err => {
        console.error(err);
        setLoading(false);
      });
  };

  useEffect(() => {
    fetchCache();
  }, []);

  const clearCache = () => {
    if (!confirm('Are you sure you want to flush the semantic cache?')) return;
    fetch('http://localhost:20128/v1/cache', { method: 'DELETE' })
      .then(() => fetchCache());
  };

  if (loading) return <div style={{ padding: '2rem' }}>Loading memory bank...</div>;

  return (
    <div className="fade-in" style={{ padding: '2rem' }}>
      <div style={{ marginBottom: '3rem', display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div>
          <h2>Semantic Cache</h2>
          <p style={{ color: 'var(--text-muted)' }}>Intercepts identical prompts and returns cached responses instantly for 0 cost and 0 latency.</p>
        </div>
        <button className="btn-secondary" onClick={clearCache} style={{ borderColor: '#e74c3c', color: '#e74c3c' }}>Flush Memory</button>
      </div>

      {cacheEntries.length === 0 ? (
        <div style={{ background: 'var(--bg-primary)', padding: '3rem', borderRadius: '12px', textAlign: 'center', border: '1px dashed var(--border-color)' }}>
          <h3 style={{ color: 'var(--text-secondary)' }}>Cache is empty</h3>
          <p style={{ color: 'var(--text-muted)' }}>Send duplicate prompts from your IDE to see them intercepted here.</p>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
          {cacheEntries.map((c, i) => (
            <div key={i} style={{ background: 'var(--bg-primary)', padding: '1.5rem', borderRadius: '8px', border: '1px solid var(--border-color)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div style={{ flex: 1, marginRight: '1rem', overflow: 'hidden' }}>
                <div style={{ fontSize: '0.8rem', color: '#3498db', marginBottom: '0.5rem', fontFamily: 'monospace' }}>{c.hash.substring(0, 16)}...</div>
                <div style={{ whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis', color: 'var(--text-primary)' }}>
                  {c.prompt}
                </div>
              </div>
              <div style={{ textAlign: 'right' }}>
                <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>{new Date(c.created_at).toLocaleString()}</div>
                <div style={{ color: '#2ecc71', fontWeight: 600, marginTop: '0.25rem' }}>{c.tokens_saved} Tokens Saved</div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
