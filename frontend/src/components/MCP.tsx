import { useState, useEffect } from 'react';

export function MCP() {
  const [data, setData] = useState<{allowed_commands: string[], allowed_dirs: string[]}>({ allowed_commands: [], allowed_dirs: [] });
  const [loading, setLoading] = useState(true);
  const [newCmd, setNewCmd] = useState('');
  const [newDir, setNewDir] = useState('');

  useEffect(() => {
    fetch('http://localhost:20128/v1/mcp')
      .then(r => r.json())
      .then(d => {
        setData({
          allowed_commands: d.allowed_commands || [],
          allowed_dirs: d.allowed_dirs || []
        });
        setLoading(false);
      });
  }, []);

  const saveSettings = (newData: any) => {
    setData(newData);
    fetch('http://localhost:20128/v1/mcp', {
      method: 'POST',
      body: JSON.stringify(newData)
    });
  };

  const removeCommand = (idx: number) => {
    const arr = [...data.allowed_commands];
    arr.splice(idx, 1);
    saveSettings({...data, allowed_commands: arr});
  };

  const addCommand = () => {
    if(!newCmd.trim()) return;
    saveSettings({...data, allowed_commands: [...data.allowed_commands, newCmd.trim()]});
    setNewCmd('');
  };

  const removeDir = (idx: number) => {
    const arr = [...data.allowed_dirs];
    arr.splice(idx, 1);
    saveSettings({...data, allowed_dirs: arr});
  };

  const addDir = () => {
    if(!newDir.trim()) return;
    saveSettings({...data, allowed_dirs: [...data.allowed_dirs, newDir.trim()]});
    setNewDir('');
  };

  if (loading) return <div style={{ padding: '2rem' }}>Loading MCP...</div>;

  return (
    <div className="fade-in" style={{ padding: '2rem' }}>
      <div style={{ marginBottom: '3rem' }}>
        <h2>Local Agents / MCP</h2>
        <p style={{ color: 'var(--text-muted)' }}>Configure the Model Context Protocol to allow the proxy to execute local shell commands.</p>
      </div>

      <div style={{ background: 'var(--bg-primary)', padding: '2rem', borderRadius: '12px', border: '1px solid var(--border-color)', maxWidth: '600px' }}>
        <h3 style={{ marginBottom: '1rem', color: '#e74c3c' }}>Security Sandbox</h3>
        <p style={{ color: 'var(--text-muted)', marginBottom: '2rem', fontSize: '0.9rem' }}>
          For security, the AI router can only execute shell commands that exactly match or start with the whitelisted prefixes below.
        </p>

        <div style={{ marginBottom: '2rem' }}>
          <label style={{ display: 'block', marginBottom: '0.5rem', fontWeight: 600 }}>Whitelisted Commands</label>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem', marginBottom: '0.5rem' }}>
            {data.allowed_commands.map((cmd, idx) => (
              <div key={idx} style={{ display: 'flex', justifyContent: 'space-between', padding: '0.5rem 1rem', background: 'var(--bg-secondary)', borderRadius: '6px', border: '1px solid var(--border-color)' }}>
                <code style={{ color: 'var(--text-primary)' }}>{cmd}</code>
                <button onClick={() => removeCommand(idx)} style={{ background: 'transparent', border: 'none', color: '#e74c3c', cursor: 'pointer' }}>✖</button>
              </div>
            ))}
          </div>
          <div style={{ display: 'flex', gap: '0.5rem' }}>
            <input type="text" className="text-input" placeholder="e.g. git status" value={newCmd} onChange={e => setNewCmd(e.target.value)} onKeyDown={e => {if(e.key==='Enter') addCommand()}} />
            <button className="btn-secondary" onClick={addCommand}>Add</button>
          </div>
        </div>

        <div style={{ marginBottom: '1rem' }}>
          <label style={{ display: 'block', marginBottom: '0.5rem', fontWeight: 600 }}>Whitelisted Directories</label>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem', marginBottom: '0.5rem' }}>
            {data.allowed_dirs.map((dir, idx) => (
              <div key={idx} style={{ display: 'flex', justifyContent: 'space-between', padding: '0.5rem 1rem', background: 'var(--bg-secondary)', borderRadius: '6px', border: '1px solid var(--border-color)' }}>
                <code style={{ color: 'var(--text-primary)' }}>{dir}</code>
                <button onClick={() => removeDir(idx)} style={{ background: 'transparent', border: 'none', color: '#e74c3c', cursor: 'pointer' }}>✖</button>
              </div>
            ))}
          </div>
          <div style={{ display: 'flex', gap: '0.5rem' }}>
            <input type="text" className="text-input" placeholder="e.g. /tmp" value={newDir} onChange={e => setNewDir(e.target.value)} onKeyDown={e => {if(e.key==='Enter') addDir()}} />
            <button className="btn-secondary" onClick={addDir}>Add</button>
          </div>
        </div>
      </div>
    </div>
  );
}
