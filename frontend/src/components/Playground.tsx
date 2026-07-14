import { useState, useEffect, useRef } from 'react';
import { API_BASE } from '../api';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';

type ProviderGroup = {
  provider_id: string;
  provider_name: string;
  models: string[];
};

type ChatSession = {
  id: string;
  title: string;
  date: string;
  messages: { role: string; content: string }[];
};

const renderMessageContent = (content: string) => {
  const thinkStart = content.indexOf('<think>');
  const thinkEnd = content.indexOf('</think>');
  
  if (thinkStart !== -1 && thinkEnd !== -1) {
    const beforeThink = content.slice(0, thinkStart);
    const thinkContent = content.slice(thinkStart + 7, thinkEnd);
    const afterThink = content.slice(thinkEnd + 8);

    return (
      <>
        {beforeThink && <div className="markdown-body"><ReactMarkdown remarkPlugins={[remarkGfm]}>{beforeThink}</ReactMarkdown></div>}
        <details style={{ marginBottom: '1rem', background: 'rgba(0,0,0,0.02)', padding: '0.75rem 1rem', borderRadius: '8px', border: '1px solid var(--border-color)' }}>
          <summary style={{ cursor: 'pointer', color: 'var(--text-muted)', fontWeight: 500, fontSize: '0.85rem', userSelect: 'none' }}>
            Thinking Process...
          </summary>
          <div className="markdown-body" style={{ marginTop: '0.75rem', fontSize: '0.9rem', color: 'var(--text-secondary)' }}>
            <ReactMarkdown remarkPlugins={[remarkGfm]}>{thinkContent.trim()}</ReactMarkdown>
          </div>
        </details>
        {afterThink && <div className="markdown-body"><ReactMarkdown remarkPlugins={[remarkGfm]}>{afterThink}</ReactMarkdown></div>}
      </>
    );
  }

  // Handle streaming/unclosed think tag
  if (thinkStart !== -1 && thinkEnd === -1) {
    const thinkContent = content.slice(thinkStart + 7);
    return (
      <details open style={{ marginBottom: '1rem', background: 'rgba(0,0,0,0.02)', padding: '0.75rem 1rem', borderRadius: '8px', border: '1px solid var(--border-color)' }}>
        <summary style={{ cursor: 'pointer', color: 'var(--text-muted)', fontWeight: 500, fontSize: '0.85rem', userSelect: 'none' }}>
          Thinking Process...
        </summary>
        <div className="markdown-body" style={{ marginTop: '0.75rem', fontSize: '0.9rem', color: 'var(--text-secondary)' }}>
          <ReactMarkdown remarkPlugins={[remarkGfm]}>{thinkContent.trim()}</ReactMarkdown>
        </div>
      </details>
    );
  }

  return <div className="markdown-body"><ReactMarkdown remarkPlugins={[remarkGfm]}>{content}</ReactMarkdown></div>;
};

export function Playground() {
  const [messages, setMessages] = useState<{role: string, content: string}[]>([]);
  const [input, setInput] = useState('');
  const [model, setModel] = useState('cta-ai-nexus');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [showModelDropdown, setShowModelDropdown] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, loading]);

  const [providerGroups, setProviderGroups] = useState<ProviderGroup[]>([]);
  const [modelNameMap, setModelNameMap] = useState<Record<string, string>>({'cta-ai-nexus': 'Smart Routing'});

  const [memoryEnabled, setMemoryEnabled] = useState(() => {
    return localStorage.getItem('cta_memory_enabled') === 'true';
  });
  const [sessions, setSessions] = useState<ChatSession[]>(() => {
    const s = localStorage.getItem('cta_chat_sessions');
    return s ? JSON.parse(s) : [];
  });
  const [currentSessionId, setCurrentSessionId] = useState<string | null>(null);

  useEffect(() => {
    fetch(`${API_BASE}/v1/models`)
      .then(r => r.json())
      .then(d => {
        if(d.data) {
          const map: Record<string, string> = {'cta-ai-nexus': 'Smart Routing'};
          const grouped: Record<string, ProviderGroup> = {};
          
          // Parse OpenAI standard model list and group by owned_by
          d.data.forEach((m: any) => {
            const pid = m.owned_by || 'unknown';
            if (!grouped[pid]) {
              grouped[pid] = {
                provider_id: pid,
                provider_name: pid.charAt(0).toUpperCase() + pid.slice(1),
                models: []
              };
            }
            grouped[pid].models.push(m.id);
            map[m.id] = m.id;
          });
          
          setProviderGroups(Object.values(grouped));
          setModelNameMap(map);
        }
      })
      .catch(console.error);
  }, []);

  // Sync memory toggle
  useEffect(() => {
    localStorage.setItem('cta_memory_enabled', String(memoryEnabled));
  }, [memoryEnabled]);

  // Sync sessions on messages update if memory is enabled
  useEffect(() => {
    if (memoryEnabled && messages.length > 0) {
      setSessions((prev: ChatSession[]) => {
        let updated = [...prev];
        const existingIdx = updated.findIndex(s => s.id === currentSessionId);
        if (existingIdx >= 0) {
          updated[existingIdx].messages = messages;
        } else {
          // New session!
          const newId = Date.now().toString();
          const title = messages[0].content.substring(0, 30) + (messages[0].content.length > 30 ? '...' : '');
          const newSession = { id: newId, title, date: new Date().toISOString(), messages };
          updated.unshift(newSession);
          setCurrentSessionId(newId);
        }
        localStorage.setItem('cta_chat_sessions', JSON.stringify(updated));
        return updated;
      });
    }
  }, [messages, memoryEnabled, currentSessionId]);

  const loadSession = (id: string) => {
    const s = sessions.find(s => s.id === id);
    if (s) {
      setMessages(s.messages);
      setCurrentSessionId(id);
    }
  };

  const newChat = () => {
    setMessages([]);
    setCurrentSessionId(null);
  };

  const deleteSession = (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    setSessions((prev) => {
      const updated = prev.filter(s => s.id !== id);
      localStorage.setItem('cta_chat_sessions', JSON.stringify(updated));
      return updated;
    });
    if (currentSessionId === id) {
      newChat();
    }
  };

  const handleSend = async () => {
    if (!input.trim()) return;

    const newMsgs = [...messages, { role: 'user', content: input }];
    setMessages(newMsgs);
    setInput('');
    setLoading(true);
    setError('');

    try {
      const res = await fetch(`${API_BASE}/v1/chat/completions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer test' // Router uses internal keys if available
        },
        body: JSON.stringify({
          model: model,
          messages: newMsgs,
          stream: false
        })
      });

      if (!res.ok) {
        const errText = await res.text();
        throw new Error(`API Error: ${res.status} - ${errText}`);
      }

      const data = await res.json();
      if (data.choices && data.choices.length > 0) {
        setMessages([...newMsgs, data.choices[0].message]);
      } else {
        throw new Error('Invalid response format');
      }
    } catch (err: any) {
      console.error(err);
      setError(err.message || 'An error occurred generating response.');
      setMessages([...newMsgs]); // Keep user message
    } finally {
      setLoading(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div style={{ display: 'flex', height: '100%', width: '100%' }}>
      {/* Sidebar for Chat Memory */}
      {memoryEnabled && (
        <div style={{ 
          width: '260px', borderRight: '1px solid var(--border-color)', background: 'var(--bg-secondary)', 
          display: 'flex', flexDirection: 'column', padding: '1rem', overflowY: 'auto'
        }}>
          <button 
            onClick={newChat}
            style={{ 
              width: '100%', padding: '0.75rem', background: 'var(--bg-primary)', border: '1px solid var(--border-color)', 
              borderRadius: '8px', cursor: 'pointer', fontWeight: 600, color: 'var(--text-primary)', marginBottom: '1rem',
              display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '8px',
              transition: 'background 0.2s'
            }}
            onMouseEnter={e => e.currentTarget.style.background = 'var(--bg-tertiary)'}
            onMouseLeave={e => e.currentTarget.style.background = 'var(--bg-primary)'}
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
            New Chat
          </button>

          <div style={{ fontSize: '0.75rem', textTransform: 'uppercase', color: 'var(--text-muted)', fontWeight: 600, letterSpacing: '0.05em', marginBottom: '0.5rem' }}>
            History
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            {sessions.length === 0 && (
              <div style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>No recent chats.</div>
            )}
            {sessions.map(s => (
              <div 
                key={s.id}
                onClick={() => loadSession(s.id)}
                style={{ 
                  padding: '0.75rem', borderRadius: '8px', cursor: 'pointer', fontSize: '0.9rem',
                  background: currentSessionId === s.id ? 'var(--bg-tertiary)' : 'transparent',
                  color: currentSessionId === s.id ? 'var(--text-primary)' : 'var(--text-secondary)',
                  border: '1px solid transparent',
                  borderColor: currentSessionId === s.id ? 'var(--border-color)' : 'transparent',
                  transition: 'background 0.1s',
                  display: 'flex', justifyContent: 'space-between', alignItems: 'center'
                }}
                onMouseEnter={e => { if (currentSessionId !== s.id) e.currentTarget.style.background = 'var(--bg-primary)' }}
                onMouseLeave={e => { if (currentSessionId !== s.id) e.currentTarget.style.background = 'transparent' }}
              >
                <div style={{ whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis', flex: 1, paddingRight: '10px' }}>
                  {s.title}
                </div>
                <button
                  onClick={(e) => deleteSession(s.id, e)}
                  style={{
                    background: 'transparent', border: 'none', color: 'var(--text-muted)', cursor: 'pointer', padding: '4px',
                    display: 'flex', alignItems: 'center', justifyContent: 'center', borderRadius: '4px'
                  }}
                  onMouseEnter={e => { e.currentTarget.style.color = '#e74c3c'; e.currentTarget.style.background = 'rgba(231,76,60,0.1)'; }}
                  onMouseLeave={e => { e.currentTarget.style.color = 'var(--text-muted)'; e.currentTarget.style.background = 'transparent'; }}
                  title="Delete Chat"
                >
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Main Chat Area */}
      <div className="fade-in" style={{ padding: '2rem', flex: 1, display: 'flex', flexDirection: 'column', position: 'relative' }}>
        
        {messages.length === 0 ? (
          <div style={{ flex: 1, position: 'relative', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
            
            <div style={{ position: 'absolute', top: 0, right: 0, display: 'flex', alignItems: 'center', gap: '8px', zIndex: 20 }}>
              <span style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>Memory</span>
              <button 
                onClick={() => setMemoryEnabled(!memoryEnabled)}
                style={{
                  width: '40px', height: '22px', borderRadius: '11px', background: memoryEnabled ? 'var(--accent-primary)' : 'var(--bg-tertiary)',
                  border: 'none', cursor: 'pointer', position: 'relative', transition: 'background 0.2s'
                }}
              >
                <div style={{
                  width: '18px', height: '18px', borderRadius: '50%', background: '#fff',
                  position: 'absolute', top: '2px', left: memoryEnabled ? '20px' : '2px',
                  transition: 'left 0.2s', boxShadow: '0 1px 3px rgba(0,0,0,0.2)'
                }} />
              </button>
            </div>

            <div style={{ 
              width: '100%', maxWidth: '800px', 
              background: 'var(--bg-primary)', borderRadius: '24px', 
              border: '1px solid var(--border-color)', padding: '1.5rem',
              boxShadow: '0 10px 30px rgba(0,0,0,0.05)',
              zIndex: 10, position: 'relative'
            }}>
              <textarea 
                placeholder="Message AI Nexus Router..."
                value={input}
                onChange={e => setInput(e.target.value)}
                onKeyDown={handleKeyDown}
                style={{
                  width: '100%', background: 'transparent', border: 'none', color: 'var(--text-primary)',
                  fontSize: '1.15rem', resize: 'none', minHeight: '60px', outline: 'none',
                  fontFamily: 'inherit'
                }}
              />
              
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: '1rem' }}>
                <div style={{ display: 'flex', gap: '0.75rem', position: 'relative' }}>
                  <button style={{ 
                    background: 'transparent', border: '1px solid var(--border-color)', borderRadius: '50%', 
                    width: '36px', height: '36px', color: 'var(--text-secondary)', fontSize: '1.5rem', 
                    display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer',
                    transition: 'background 0.2s'
                  }}
                  onMouseEnter={e => e.currentTarget.style.background = 'var(--bg-secondary)'}
                  onMouseLeave={e => e.currentTarget.style.background = 'transparent'}
                  >+</button>
                  
                  <button 
                    onClick={() => setShowModelDropdown(!showModelDropdown)}
                    style={{
                      background: 'transparent', border: '1px solid var(--border-color)', borderRadius: '18px',
                      padding: '0 16px', color: 'var(--text-secondary)', fontSize: '0.95rem', display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer',
                      transition: 'background 0.2s', fontWeight: 500, whiteSpace: 'nowrap', maxWidth: '250px', overflow: 'hidden', textOverflow: 'ellipsis'
                    }}
                    onMouseEnter={e => e.currentTarget.style.background = 'var(--bg-secondary)'}
                    onMouseLeave={e => e.currentTarget.style.background = 'transparent'}
                  >
                    {modelNameMap[model] || model}
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ flexShrink: 0 }}><polyline points="6 9 12 15 18 9"></polyline></svg>
                  </button>

                  {showModelDropdown && (
                    <div style={{
                      position: 'absolute', top: '100%', left: '44px', marginTop: '12px',
                      background: 'var(--bg-primary)', border: '1px solid var(--border-color)', borderRadius: '12px',
                      width: '340px', padding: '0.5rem', boxShadow: '0 10px 30px rgba(0,0,0,0.1)', zIndex: 100,
                      maxHeight: '400px', overflowY: 'auto'
                    }}>
                      
                      <div 
                        onClick={() => { setModel('cta-ai-nexus'); setShowModelDropdown(false); }}
                        style={{
                          padding: '0.875rem', borderRadius: '8px', cursor: 'pointer', display: 'flex', gap: '14px', alignItems: 'flex-start',
                          background: model === 'cta-ai-nexus' ? 'var(--bg-secondary)' : 'transparent', transition: 'background 0.1s', marginBottom: '0.5rem'
                        }}
                        onMouseEnter={e => e.currentTarget.style.background = 'var(--bg-tertiary)'}
                        onMouseLeave={e => e.currentTarget.style.background = model === 'cta-ai-nexus' ? 'var(--bg-secondary)' : 'transparent'}
                      >
                        <div style={{ fontSize: '1.4rem' }}>✨</div>
                        <div>
                          <div style={{ color: 'var(--text-primary)', fontWeight: 600, fontSize: '0.95rem' }}>Smart Routing</div>
                          <div style={{ color: 'var(--text-muted)', fontSize: '0.85rem', marginTop: '2px' }}>Auto-routes based on prompt</div>
                        </div>
                      </div>

                      {providerGroups.map(pg => (
                        <div key={pg.provider_id} style={{ marginBottom: '0.5rem' }}>
                          <div style={{ padding: '0.25rem 0.875rem', fontSize: '0.75rem', textTransform: 'uppercase', color: 'var(--text-muted)', fontWeight: 600, letterSpacing: '0.05em' }}>
                            {pg.provider_name}
                          </div>
                          {pg.models.map(m => (
                            <div 
                              key={m}
                              onClick={() => { setModel(m); setShowModelDropdown(false); }}
                              style={{
                                padding: '0.6rem 0.875rem', borderRadius: '8px', cursor: 'pointer', display: 'flex', alignItems: 'center',
                                background: model === m ? 'var(--bg-secondary)' : 'transparent', transition: 'background 0.1s',
                                color: 'var(--text-secondary)', fontSize: '0.9rem'
                              }}
                              onMouseEnter={e => e.currentTarget.style.background = 'var(--bg-tertiary)'}
                              onMouseLeave={e => e.currentTarget.style.background = model === m ? 'var(--bg-secondary)' : 'transparent'}
                            >
                              {m}
                            </div>
                          ))}
                        </div>
                      ))}

                      {providerGroups.length === 0 && (
                        <div style={{ padding: '1rem', color: 'var(--text-muted)', textAlign: 'center', fontSize: '0.9rem' }}>
                          Add API keys in Providers tab to see models here.
                        </div>
                      )}

                    </div>
                  )}
                </div>

                {input.trim() && (
                  <button onClick={handleSend} style={{
                    background: 'var(--accent-primary)', color: '#fff', border: 'none', borderRadius: '18px',
                    padding: '8px 20px', fontWeight: 600, cursor: 'pointer', transition: 'transform 0.1s, background 0.2s',
                    boxShadow: '0 4px 10px rgba(234, 105, 71, 0.3)'
                  }}
                  onMouseEnter={e => e.currentTarget.style.transform = 'scale(1.05)'}
                  onMouseLeave={e => e.currentTarget.style.transform = 'scale(1)'}
                  >
                    Send
                  </button>
                )}
              </div>
            </div>

            <div style={{ position: 'absolute', top: 'calc(50% + 100px)', display: 'flex', flexWrap: 'wrap', justifyContent: 'center', gap: '14px', maxWidth: '900px', zIndex: 1 }}>
              {["Create something new", "Take a quiz", "Learn a new skill", "Brainpower", "Practice a language", "Code a python script"].map((txt, i) => (
                <div 
                  key={i} 
                  onClick={() => setInput(txt)}
                  style={{
                    background: 'var(--bg-primary)', border: '1px solid var(--border-color)',
                    color: 'var(--text-secondary)', padding: '10px 20px', borderRadius: '24px', fontSize: '0.9rem',
                    cursor: 'pointer', transition: 'all 0.2s',
                    fontWeight: 500, boxShadow: '0 2px 4px rgba(0,0,0,0.02)'
                  }}
                  onMouseEnter={e => { e.currentTarget.style.background = 'var(--bg-secondary)'; e.currentTarget.style.color = 'var(--text-primary)'; e.currentTarget.style.borderColor = '#d1d5db'; }}
                  onMouseLeave={e => { e.currentTarget.style.background = 'var(--bg-primary)'; e.currentTarget.style.color = 'var(--text-secondary)'; e.currentTarget.style.borderColor = 'var(--border-color)'; }}
                >
                  {txt}
                </div>
              ))}
            </div>
            
          </div>
        ) : (
          <>
            <div style={{ marginBottom: '1rem', display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
              <div>
                <h2>Playground</h2>
                <p style={{ color: 'var(--text-muted)' }}>Testing AI Nexus Router proxy locally.</p>
              </div>
              <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <span style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>Memory</span>
                  <button 
                    onClick={() => setMemoryEnabled(!memoryEnabled)}
                    style={{
                      width: '40px', height: '22px', borderRadius: '11px', background: memoryEnabled ? 'var(--accent-primary)' : 'var(--bg-tertiary)',
                      border: 'none', cursor: 'pointer', position: 'relative', transition: 'background 0.2s'
                    }}
                  >
                    <div style={{
                      width: '18px', height: '18px', borderRadius: '50%', background: '#fff',
                      position: 'absolute', top: '2px', left: memoryEnabled ? '20px' : '2px',
                      transition: 'left 0.2s', boxShadow: '0 1px 3px rgba(0,0,0,0.2)'
                    }} />
                  </button>
                </div>
                
                <select 
                  value={model}
                  onChange={(e) => setModel(e.target.value)}
                  style={{ 
                    padding: '0.5rem 1rem', background: 'var(--bg-primary)', border: '1px solid var(--border-color)', 
                    borderRadius: '8px', fontSize: '0.9rem', color: 'var(--text-primary)', cursor: 'pointer', outline: 'none',
                    fontFamily: 'inherit', maxWidth: '300px'
                  }}
                >
                  <option value="cta-ai-nexus">✨ Smart Routing</option>
                  {providerGroups.map(pg => (
                    <optgroup key={pg.provider_id} label={pg.provider_name}>
                      {pg.models.map(m => (
                        <option key={m} value={m}>{m}</option>
                      ))}
                    </optgroup>
                  ))}
                </select>
              </div>
            </div>

            <div style={{ flex: 1, background: 'var(--bg-primary)', borderRadius: '12px', border: '1px solid var(--border-color)', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
              
              <div style={{ flex: 1, overflowY: 'auto', padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
                {messages.map((m, i) => (
                  <div key={i} style={{ 
                    alignSelf: m.role === 'user' ? 'flex-end' : 'flex-start',
                    background: m.role === 'user' ? 'var(--accent-primary)' : 'var(--bg-tertiary)',
                    padding: m.role === 'user' ? '1rem 1.5rem' : '1.5rem',
                    borderRadius: '16px',
                    maxWidth: '85%',
                    border: m.role === 'user' ? 'none' : '1px solid var(--border-color)',
                    color: m.role === 'user' ? '#fff' : 'var(--text-primary)',
                    boxShadow: m.role === 'user' ? '0 4px 10px rgba(234, 105, 71, 0.2)' : '0 2px 4px rgba(0,0,0,0.02)',
                    lineHeight: '1.6'
                  }}>
                    <div style={{ fontSize: '0.75rem', opacity: 0.8, marginBottom: '0.75rem', textTransform: 'uppercase', fontWeight: 600, letterSpacing: '0.05em' }}>
                      {m.role === 'user' ? 'You' : 'Assistant'}
                    </div>
                    {m.role === 'user' ? (
                      <div style={{ whiteSpace: 'pre-wrap' }}>{m.content}</div>
                    ) : (
                      renderMessageContent(m.content)
                    )}
                  </div>
                ))}
                {loading && (
                  <div style={{ alignSelf: 'flex-start', padding: '1rem', color: 'var(--text-muted)' }}>
                    Routing request...
                  </div>
                )}
                {error && (
                  <div style={{ alignSelf: 'center', padding: '1rem', color: '#e74c3c', background: 'rgba(231,76,60,0.1)', borderRadius: '8px', marginTop: '1rem' }}>
                    {error}
                  </div>
                )}
                <div ref={messagesEndRef} />
              </div>

              <div style={{ padding: '1rem', borderTop: '1px solid var(--border-color)', background: 'var(--bg-secondary)', display: 'flex', gap: '1rem' }}>
                <input 
                  type="text" 
                  className="text-input" 
                  placeholder="Message AI Nexus Router..." 
                  value={input}
                  onChange={e => setInput(e.target.value)}
                  onKeyDown={handleKeyDown}
                  style={{ flex: 1 }}
                  disabled={loading}
                />
                <button className="btn-primary" onClick={handleSend} disabled={loading || !input.trim()} style={{ width: 'auto', padding: '0 2rem' }}>
                  Send
                </button>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
