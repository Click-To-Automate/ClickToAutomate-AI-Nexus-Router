import { useState, useEffect } from 'react';

interface ProviderKey {
  provider_id: string;
  masked_key: string;
  api_key: string;
}

const LOGO = (domain: string) => `https://www.google.com/s2/favicons?domain=${domain}&sz=64`;

const AUTH_BADGE: Record<string, { label: string; color: string }> = {
  api_key:    { label: 'API Key',    color: '#3b82f6' },
  oauth:      { label: 'OAuth',      color: '#8b5cf6' },
  web_cookie: { label: 'Web Cookie', color: '#f59e0b' },
  local:      { label: 'Local',      color: '#22c55e' },
  bearer_token:{ label: 'Token',     color: '#06b6d4' },
};

const HARDCODED_META: Record<string, any> = {
  // ── API Key ───────────────────────────────────────────────────────────────
  'openai':        { name: 'OpenAI',         url: 'https://platform.openai.com',           logo: LOGO('openai.com'),            color: '#10a37f', tier: 'Paid',  models: 14,      description: 'Industry-leading GPT-4o, o1, and o3 models with multimodal reasoning.',                                              integration: 'api_key',    pricing: 'paid',  account: 'needed',     keyLabel: 'API Key',                  keyHint: 'Get at platform.openai.com/api-keys' },
  'anthropic':     { name: 'Anthropic',      url: 'https://console.anthropic.com',          logo: LOGO('anthropic.com'),         color: '#d35400', tier: 'Paid',  models: 6,       description: 'Claude 3.5 family — 200K context, safety-first with exceptional code generation.',                                    integration: 'api_key',    pricing: 'paid',  account: 'needed',     keyLabel: 'API Key',                  keyHint: 'Get at console.anthropic.com/settings/keys' },
  'gemini':        { name: 'Google Gemini',  url: 'https://aistudio.google.com',           logo: LOGO('aistudio.google.com'),   color: '#4285F4', tier: 'Free',  models: 8,       description: "Natively multimodal with 2M-token context and massive free tiers.",                                                   integration: 'api_key',    pricing: 'free',  account: 'needed',     keyLabel: 'API Key',                  keyHint: 'Get at aistudio.google.com/apikey' },
  'groq':          { name: 'Groq',           url: 'https://console.groq.com',              logo: LOGO('groq.com'),              color: '#F55036', tier: 'Free',  models: 12,      description: 'LPU inference — hundreds of tokens/sec on LLaMA, Gemma, and Mixtral.',                                               integration: 'api_key',    pricing: 'free',  account: 'needed',     keyLabel: 'API Key',                  keyHint: 'Get at console.groq.com/keys' },
  'openrouter':    { name: 'OpenRouter',     url: 'https://openrouter.ai',                 logo: LOGO('openrouter.ai'),         color: '#9b59b6', tier: 'Free',  models: 250,     description: 'Unified API: 250+ models from OpenAI, Anthropic, Meta, Google and open-source.',                                     integration: 'api_key',    pricing: 'free',  account: 'needed',     keyLabel: 'API Key',                  keyHint: 'Get at openrouter.ai/keys' },
  'mistral':       { name: 'Mistral',        url: 'https://console.mistral.ai',            logo: LOGO('mistral.ai'),            color: '#f39c12', tier: 'Free',  models: 10,      description: "Europe's leading AI lab — efficient open-weight models with vision via Pixtral.",                                        integration: 'api_key',    pricing: 'free',  account: 'needed',     keyLabel: 'API Key',                  keyHint: 'Get at console.mistral.ai/api-keys' },
  'cerebras':      { name: 'Cerebras',       url: 'https://cloud.cerebras.ai',             logo: LOGO('cerebras.ai'),           color: '#e74c3c', tier: 'Free',  models: 4,       description: 'Wafer-scale silicon for sub-100ms full-response latency on LLaMA and Qwen.',                                          integration: 'api_key',    pricing: 'free',  account: 'needed',     keyLabel: 'API Key',                  keyHint: 'Get at cloud.cerebras.ai' },
  'deepseek':      { name: 'DeepSeek',       url: 'https://platform.deepseek.com',         logo: LOGO('deepseek.com'),          color: '#2980b9', tier: 'Free',  models: 4,       description: 'MoE architecture with state-of-the-art reasoning at ultra-low cost.',                                                 integration: 'api_key',    pricing: 'free',  account: 'needed',     keyLabel: 'API Key',                  keyHint: 'Get at platform.deepseek.com' },
  'together':      { name: 'Together AI',    url: 'https://api.together.xyz',              logo: LOGO('together.ai'),           color: '#34495e', tier: 'Paid',  models: 100,     description: 'Cloud platform for running, fine-tuning and deploying 100+ open-source models.',                                     integration: 'api_key',    pricing: 'paid',  account: 'needed',     keyLabel: 'API Key',                  keyHint: 'Get at api.together.xyz/settings/api-keys' },
  'sambanova':     { name: 'SambaNova',      url: 'https://cloud.sambanova.ai',            logo: LOGO('sambanova.ai'),          color: '#e67e22', tier: 'Free',  models: 6,       description: 'SN40L chips powering LLaMA 3.1 405B at full precision and industry-leading speed.',                                   integration: 'api_key',    pricing: 'free',  account: 'needed',     keyLabel: 'API Key',                  keyHint: 'Get at cloud.sambanova.ai' },
  'llm7':          { name: 'LLM7',           url: 'https://llm7.io',                       logo: LOGO('llm7.io'),               color: '#171717', tier: 'Free',  models: 4,       description: 'Completely free inference — no API key required. GPT-4o-class models at zero cost.',                                  integration: 'api_key',    pricing: 'free',  account: 'not_needed', keyLabel: 'API Key (optional)',       keyHint: 'Leave blank — LLM7 requires no key' },
  'github-models': { name: 'GitHub Models',  url: 'https://github.com/marketplace/models', logo: LOGO('github.com'),            color: '#24292e', tier: 'Free',  models: 15,      description: 'Free OpenAI-compatible API for top models via your GitHub Personal Access Token.',                                   integration: 'api_key',    pricing: 'free',  account: 'needed',     keyLabel: 'GitHub Token (PAT)',       keyHint: 'Generate at github.com/settings/tokens → select models scope' },
  // ── Local ─────────────────────────────────────────────────────────────────
  'ollama':        { name: 'Ollama',         url: 'https://ollama.com',                    logo: LOGO('ollama.com'),            color: '#2c3e50', tier: 'Free',  models: 'Local', description: 'Run any open model locally — one-command install, no cloud required.',                                                 integration: 'local',      pricing: 'free',  account: 'not_needed', keyLabel: 'Not required',             keyHint: 'Ollama runs on http://localhost:11434 — no key needed' },
  // ── Web Cookie ────────────────────────────────────────────────────────────
  'deepseek-web':  { name: 'DeepSeek Web',   url: 'https://chat.deepseek.com',             logo: LOGO('chat.deepseek.com'),     color: '#2980b9', tier: 'Free',  models: 2,       description: 'Use your DeepSeek web session by pasting the userToken from chat.deepseek.com localStorage.',                         integration: 'web_cookie', pricing: 'free',  account: 'needed',     keyLabel: 'Session Token',            keyHint: 'DevTools → Application → Local Storage → https://chat.deepseek.com → userToken' },
  'gemini-web':    { name: 'Gemini Web',     url: 'https://gemini.google.com',             logo: LOGO('gemini.google.com'),     color: '#4285F4', tier: 'Free',  models: 3,       description: 'Use your Gemini web session by pasting __Secure-1PSID from gemini.google.com cookies.',                               integration: 'web_cookie', pricing: 'free',  account: 'needed',     keyLabel: '__Secure-1PSID Cookie',    keyHint: 'DevTools → Application → Cookies → gemini.google.com → __Secure-1PSID' },
  'claude-web':    { name: 'Claude Web',     url: 'https://claude.ai',                     logo: LOGO('claude.ai'),             color: '#d35400', tier: 'Free',  models: 4,       description: 'Use your Claude.ai web session by pasting the session cookie from DevTools.',                                          integration: 'web_cookie', pricing: 'free',  account: 'needed',     keyLabel: 'Session Cookie',           keyHint: 'DevTools → Application → Cookies → claude.ai → copy the full Cookie header' },
  'qwen-web':      { name: 'Qwen Web',       url: 'https://chat.qwen.ai',                  logo: LOGO('chat.qwen.ai'),          color: '#8e44ad', tier: 'Free',  models: 5,       description: 'Use your Qwen web session by pasting the auth token from chat.qwen.ai Local Storage.',                                 integration: 'web_cookie', pricing: 'free',  account: 'needed',     keyLabel: 'Auth Token',               keyHint: 'DevTools → Application → Local Storage → https://chat.qwen.ai → token' },
  // ── OAuth ─────────────────────────────────────────────────────────────────
  'github':        { name: 'GitHub Copilot', url: 'https://github.com/settings/tokens',    logo: LOGO('github.com'),            color: '#24292e', tier: 'Paid',  models: 8,       description: 'GitHub Copilot API — paste your OAuth token or PAT with the copilot scope enabled.',                                  integration: 'oauth',      pricing: 'paid',  account: 'needed',     keyLabel: 'OAuth / PAT Token',        keyHint: 'github.com/settings/tokens → New token → select copilot scope' },
  'cursor':        { name: 'Cursor IDE',     url: 'https://cursor.com',                    logo: LOGO('cursor.com'),            color: '#1a1a2e', tier: 'Paid',  models: 6,       description: 'Use Cursor IDE models directly — paste the access token from your Cursor session.',                                    integration: 'oauth',      pricing: 'paid',  account: 'needed',     keyLabel: 'Cursor Access Token',      keyHint: 'Cursor IDE → Command Palette → "Cursor: Provide Auth Token"' },
};

const AUTH_INSTRUCTIONS: Record<string, { title: string; steps: string[] }> = {
  api_key: {
    title: 'How to get your API Key',
    steps: ['Sign in to the provider dashboard', 'Navigate to API Keys / Settings', 'Create a new key and copy it', 'Paste it below'],
  },
  web_cookie: {
    title: 'How to extract your session cookie',
    steps: ['Open the provider website and sign in', 'Press F12 to open DevTools', 'Go to Application → Cookies (or Local Storage)', 'Find and copy the session token/cookie value', 'Paste it below — the router will use it as-is'],
  },
  oauth: {
    title: 'How to get your OAuth token',
    steps: ['Sign in to the provider dashboard', 'Generate a Personal Access Token or OAuth token', 'Select the required scopes (see hint below)', 'Copy and paste the token below'],
  },
  local: {
    title: 'Local provider — no key required',
    steps: ['Make sure Ollama is installed and running (ollama serve)', 'Pull a model: ollama pull llama3', 'No key needed — the router connects to localhost:11434 automatically'],
  },
  bearer_token: {
    title: 'How to get your session token',
    steps: ['Sign in to the provider website', 'Open DevTools (F12) → Application → Local Storage', 'Find the token key and copy its value', 'Paste it below'],
  },
};

export function ProvidersSettings() {
  const [savedKeys, setSavedKeys] = useState<ProviderKey[]>([]);
  const [providersList, setProvidersList] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [searchQuery, setSearchQuery] = useState('');
  const [filterIntegration, setFilterIntegration] = useState('all');
  const [filterPricing, setFilterPricing] = useState('all');
  const [filterAccount, setFilterAccount] = useState('all');
  const [filterStatus, setFilterStatus] = useState('all');

  // Modal State
  const [selectedProvider, setSelectedProvider] = useState<string | null>(null);
  const [newKeyInput, setNewKeyInput] = useState('');

  useEffect(() => { fetchData(); }, []);

  const fetchData = async () => {
    try {
      const [keysRes, provRes] = await Promise.all([
        fetch('http://localhost:20128/v1/keys'),
        fetch('http://localhost:20128/v1/providers'),
      ]);
      const keysData = await keysRes.json();
      const provData = await provRes.json();
      setSavedKeys(keysData.keys || []);

      const dynProviders = (provData.providers || []).map((p: any) => {
        const meta = HARDCODED_META[p.id] || {};
        return {
          id: p.id,
          name: meta.name || p.name || p.id,
          url: meta.url || p.base_url || 'https://example.com',
          logo: meta.logo || LOGO(new URL(p.base_url || 'https://example.com').hostname),
          color: meta.color || '#6366f1',
          tier: meta.tier || 'API Key',
          models: meta.models || 'Dynamic',
          description: meta.description || `OpenAI-compatible endpoint for ${meta.name || p.name || p.id}.`,
          integration: meta.integration || 'api_key',
          pricing: meta.pricing || 'paid',
          account: meta.account || 'needed',
          keyLabel: meta.keyLabel || 'API Key',
          keyHint: meta.keyHint || `Get your key from ${meta.url || p.base_url}`,
          authType: p.auth_type || 'bearer',
        };
      });
      setProvidersList(dynProviders);
    } catch (err) {
      console.error('Failed to fetch data', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchKeys = async () => {
    try {
      const res = await fetch('http://localhost:20128/v1/keys');
      const data = await res.json();
      setSavedKeys(data.keys || []);
    } catch (err) {
      console.error('Failed to fetch keys', err);
    }
  };

  const handleAddKey = async () => {
    if (!selectedProvider || !newKeyInput) return;
    try {
      await fetch('http://localhost:20128/v1/keys', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ provider_id: selectedProvider, api_key: newKeyInput, action: 'add' }),
      });
      setNewKeyInput('');
      fetchKeys();
    } catch (err) {
      console.error('Failed to save key', err);
    }
  };

  const handleDeleteKey = async (apiKey: string) => {
    if (!selectedProvider) return;
    try {
      await fetch('http://localhost:20128/v1/keys', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ provider_id: selectedProvider, api_key: apiKey, action: 'delete' }),
      });
      fetchKeys();
    } catch (err) {
      console.error('Failed to delete key', err);
    }
  };

  const closeModal = () => { setSelectedProvider(null); setNewKeyInput(''); };

  if (loading) return <div style={{ padding: '2rem' }}>Loading providers...</div>;

  const getProviderKeys = (providerId: string) => savedKeys.filter(k => k.provider_id === providerId);

  const filteredProviders = providersList.filter(p => {
    const matchesSearch = p.name.toLowerCase().includes(searchQuery.toLowerCase()) || p.id.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesIntegration = filterIntegration === 'all' || p.integration === filterIntegration;
    const matchesPricing = filterPricing === 'all' || p.pricing === filterPricing;
    const matchesAccount = filterAccount === 'all' || p.account === filterAccount;
    const keysCount = getProviderKeys(p.id).length;
    const matchesStatus = filterStatus === 'all' || (filterStatus === 'connected' && keysCount > 0) || (filterStatus === 'setup' && keysCount === 0);
    return matchesSearch && matchesIntegration && matchesPricing && matchesAccount && matchesStatus;
  });

  const activeProvider = providersList.find(p => p.id === selectedProvider);
  const activeKeys = selectedProvider ? getProviderKeys(selectedProvider) : [];
  const activeMeta = HARDCODED_META[selectedProvider || ''] || {};
  const activeAuthType = activeProvider?.integration || 'api_key';
  const instructions = AUTH_INSTRUCTIONS[activeAuthType] || AUTH_INSTRUCTIONS['api_key'];

  const ProviderLogo = ({ p, size = 36 }: { p: any; size?: number }) => (
    <img
      src={p.logo}
      alt={p.name}
      width={size}
      height={size}
      style={{ borderRadius: size * 0.22, objectFit: 'contain', background: '#fff', padding: '2px' }}
      onError={(e) => {
        e.currentTarget.onerror = null;
        e.currentTarget.src = `data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><rect width='100' height='100' fill='${encodeURIComponent(p.color)}' rx='18'/><text x='50' y='62' font-family='Arial' font-size='46' font-weight='bold' fill='white' text-anchor='middle'>${p.name.charAt(0).toUpperCase()}</text></svg>`;
      }}
    />
  );

  const PILL_STYLE = (active: boolean) => ({
    display: 'flex' as const, alignItems: 'center' as const, gap: '0.3rem',
    padding: '0.3rem 0.75rem', borderRadius: '999px', fontSize: '0.78rem', fontWeight: 500,
    cursor: 'pointer' as const, transition: 'all 0.15s ease',
    border: active ? '1.5px solid var(--accent-primary)' : '1.5px solid var(--border-color)',
    background: active ? 'var(--accent-primary)' : 'var(--bg-primary)',
    color: active ? '#fff' : 'var(--text-secondary)',
    boxShadow: active ? '0 2px 8px var(--accent-glow)' : 'none',
  });

  return (
    <div className="fade-in" style={{ paddingBottom: '2rem' }}>
      <div style={{ marginBottom: '3rem', display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div>
          <h2>Provider Configurations</h2>
          <p style={{ color: 'var(--text-muted)' }}>
            Click a provider to add, manage, and delete API keys. The router automatically load-balances requests across all configured keys for a provider!
          </p>
          <div style={{ marginTop: '1rem', display: 'flex', flexWrap: 'wrap', gap: '0.75rem', alignItems: 'center' }}>
            <input
              type="text" className="text-input" placeholder="Search providers..."
              value={searchQuery} onChange={e => setSearchQuery(e.target.value)}
              style={{ width: '220px' }}
            />

            {/* Integration pills */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.35rem', flexWrap: 'wrap' }}>
              {([
                { value: 'all',        label: 'All',        icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg> },
                { value: 'api_key',    label: 'API Key',    icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg> },
                { value: 'oauth',      label: 'OAuth',      icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg> },
                { value: 'web_cookie', label: 'Web Cookie', icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg> },
                { value: 'local',      label: 'Local',      icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><rect x="2" y="3" width="20" height="14" rx="2"/><path d="M8 21h8M12 17v4"/></svg> },
              ] as const).map(opt => (
                <button key={opt.value} onClick={() => setFilterIntegration(opt.value)} style={PILL_STYLE(filterIntegration === opt.value)}>
                  {opt.icon}{opt.label}
                </button>
              ))}
            </div>

            <div style={{ width: '1px', height: '24px', background: 'var(--border-color)' }} />

            {/* Pricing pills */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.35rem' }}>
              {([
                { value: 'all',  label: 'All Pricing', icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg> },
                { value: 'free', label: 'Free',        icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg> },
                { value: 'paid', label: 'Paid',        icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><rect x="1" y="4" width="22" height="16" rx="2" ry="2"/><line x1="1" y1="10" x2="23" y2="10"/></svg> },
              ] as const).map(opt => (
                <button key={opt.value} onClick={() => setFilterPricing(opt.value)} style={PILL_STYLE(filterPricing === opt.value)}>
                  {opt.icon}{opt.label}
                </button>
              ))}
            </div>

            <div style={{ width: '1px', height: '24px', background: 'var(--border-color)' }} />

            {/* Account pills */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.35rem' }}>
              {([
                { value: 'all',        label: 'All',             icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/></svg> },
                { value: 'needed',     label: 'Account Needed',  icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg> },
                { value: 'not_needed', label: 'No Account',      icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="23" y1="11" x2="17" y2="11"/></svg> },
              ] as const).map(opt => (
                <button key={opt.value} onClick={() => setFilterAccount(opt.value)} style={PILL_STYLE(filterAccount === opt.value)}>
                  {opt.icon}{opt.label}
                </button>
              ))}
            </div>

            <div style={{ width: '1px', height: '24px', background: 'var(--border-color)' }} />

            {/* Status pills */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.35rem' }}>
              {([
                { value: 'all',       label: 'All',          icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/></svg> },
                { value: 'connected', label: 'Connected',    icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg> },
                { value: 'setup',     label: 'Needs Setup',  icon: <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg> },
              ] as const).map(opt => (
                <button key={opt.value} onClick={() => setFilterStatus(opt.value)} style={PILL_STYLE(filterStatus === opt.value)}>
                  {opt.icon}{opt.label}
                </button>
              ))}
            </div>
          </div>
        </div>
        <div className="view-toggle">
          <button className={`view-toggle-btn ${viewMode === 'grid' ? 'active' : ''}`} onClick={() => setViewMode('grid')}>Grid</button>
          <button className={`view-toggle-btn ${viewMode === 'list' ? 'active' : ''}`} onClick={() => setViewMode('list')}>List</button>
        </div>
      </div>

      <div className={viewMode === 'grid' ? 'provider-grid' : 'provider-list'}>
        {filteredProviders.map(p => {
          const keysCount = getProviderKeys(p.id).length;
          const isGrid = viewMode === 'grid';
          const badge = AUTH_BADGE[p.integration] || AUTH_BADGE['api_key'];

          return (
            <div
              key={p.id}
              className={isGrid ? 'provider-card-v2' : 'provider-card-list'}
              onClick={() => setSelectedProvider(p.id)}
              style={{
                borderColor: keysCount > 0 ? 'rgba(34, 197, 94, 0.5)' : 'var(--border-color)',
                boxShadow: keysCount > 0 ? '0 0 15px rgba(34, 197, 94, 0.1), 0 4px 6px -1px rgba(0,0,0,0.02)' : '',
              }}
            >
              <div className="provider-card-top-bar" style={{ background: `linear-gradient(90deg, ${p.color}, transparent)` }} />

              <div className="provider-card-header">
                <div className="provider-card-avatar">
                  <ProviderLogo p={p} size={36} />
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '0.25rem', alignItems: 'flex-end' }}>
                  {/* Auth type badge */}
                  <span style={{
                    fontSize: '0.62rem', fontWeight: 700, letterSpacing: '0.04em',
                    padding: '0.15rem 0.5rem', borderRadius: '999px', textTransform: 'uppercase',
                    background: `${badge.color}22`, color: badge.color, border: `1px solid ${badge.color}44`,
                  }}>{badge.label}</span>
                  {/* Active/New badge */}
                  {keysCount > 0
                    ? <span className="badge badge-active">ACTIVE</span>
                    : <span className="badge badge-inactive">NEW</span>
                  }
                </div>
              </div>

              <div className="provider-card-body">
                <h3 className="provider-card-title">{p.name}</h3>
                <p className="provider-card-subtitle">
                  <a href={p.url} target="_blank" rel="noreferrer"
                    style={{ color: 'var(--accent-primary)', textDecoration: 'none' }}
                    onClick={e => e.stopPropagation()}>
                    {p.integration === 'web_cookie' ? 'Open Website ↗' : p.integration === 'oauth' ? 'Get Token ↗' : 'Get API Key ↗'}
                  </a>
                </p>
                <p className="provider-card-desc">{p.description}</p>
              </div>

              <div className="provider-card-footer">
                {!isGrid && (
                  <div className="stat-block" style={{ textAlign: 'left', alignItems: 'flex-start' }}>
                    {keysCount > 0 ? <span className="badge badge-active">ACTIVE</span> : <span className="badge badge-inactive">NEW</span>}
                  </div>
                )}
                <div className="stat-block">
                  <span className="stat-label">TIER</span>
                  <span className="stat-value">{p.tier}</span>
                </div>
                <div className="stat-block">
                  <span className="stat-label">MODELS</span>
                  <span className="stat-value">{p.models}</span>
                </div>
                <div className="stat-block">
                  <span className="stat-label">KEYS</span>
                  <span className="stat-value">{keysCount}</span>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Modal */}
      {selectedProvider && activeProvider && (
        <div className="modal-overlay" onClick={closeModal}>
          <div className="modal-content fade-in" onClick={e => e.stopPropagation()} style={{ maxWidth: '560px' }}>

            {/* Header */}
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.75rem' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.9rem' }}>
                <ProviderLogo p={activeProvider} size={40} />
                <div>
                  <h2 style={{ margin: 0, fontFamily: "'Inter Tight', sans-serif", fontSize: '1.25rem' }}>{activeProvider.name}</h2>
                  <span style={{
                    fontSize: '0.68rem', fontWeight: 700, letterSpacing: '0.05em',
                    padding: '0.15rem 0.55rem', borderRadius: '999px', textTransform: 'uppercase',
                    background: `${(AUTH_BADGE[activeAuthType] || AUTH_BADGE['api_key']).color}22`,
                    color: (AUTH_BADGE[activeAuthType] || AUTH_BADGE['api_key']).color,
                    border: `1px solid ${(AUTH_BADGE[activeAuthType] || AUTH_BADGE['api_key']).color}44`,
                  }}>{(AUTH_BADGE[activeAuthType] || AUTH_BADGE['api_key']).label}</span>
                </div>
              </div>
              <button onClick={closeModal} style={{ fontSize: '1.5rem', background: 'transparent', border: 'none', cursor: 'pointer', color: 'var(--text-muted)', lineHeight: 1 }}>×</button>
            </div>

            {/* How-to steps */}
            <div style={{ background: 'var(--bg-tertiary)', borderRadius: '10px', padding: '1rem 1.25rem', marginBottom: '1.75rem', border: '1px solid var(--border-color)' }}>
              <p style={{ margin: '0 0 0.6rem', fontSize: '0.78rem', fontWeight: 700, color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.06em' }}>
                {instructions.title}
              </p>
              <ol style={{ margin: 0, paddingLeft: '1.2rem', display: 'flex', flexDirection: 'column', gap: '0.3rem' }}>
                {instructions.steps.map((step, i) => (
                  <li key={i} style={{ fontSize: '0.82rem', color: 'var(--text-muted)', lineHeight: 1.5 }}>{step}</li>
                ))}
              </ol>
            </div>

            {/* Active Keys */}
            <div style={{ marginBottom: '1.75rem' }}>
              <h4 style={{ marginBottom: '0.75rem', color: 'var(--text-secondary)', fontSize: '0.85rem', textTransform: 'uppercase', letterSpacing: '0.06em' }}>
                Configured {activeMeta.keyLabel || 'Keys'} ({activeKeys.length})
              </h4>
              {activeKeys.length === 0 ? (
                <div style={{ padding: '1.25rem', background: 'var(--bg-tertiary)', borderRadius: '8px', textAlign: 'center', color: 'var(--text-muted)', fontSize: '0.85rem' }}>
                  No credentials configured yet.
                </div>
              ) : (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                  {activeKeys.map((k, idx) => (
                    <div key={idx} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '0.85rem 1rem', border: '1px solid var(--border-color)', borderRadius: '8px', background: 'var(--bg-primary)' }}>
                      <span style={{ fontFamily: 'monospace', fontSize: '0.95rem', color: 'var(--text-primary)' }}>{k.masked_key}</span>
                      <button onClick={() => handleDeleteKey(k.api_key)} style={{ background: 'transparent', border: 'none', color: '#ef4444', cursor: 'pointer', fontWeight: 600, fontSize: '0.82rem' }}>
                        Remove
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Add Key */}
            <div style={{ borderTop: '1px solid var(--border-color)', paddingTop: '1.5rem' }}>
              <h4 style={{ marginBottom: '0.4rem', color: 'var(--text-secondary)', fontSize: '0.85rem', textTransform: 'uppercase', letterSpacing: '0.06em' }}>
                Add {activeMeta.keyLabel || 'API Key'}
              </h4>
              {activeMeta.keyHint && (
                <p style={{ margin: '0 0 0.75rem', fontSize: '0.75rem', color: 'var(--text-muted)', display: 'flex', alignItems: 'center', gap: '0.35rem' }}>
                  <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
                  {activeMeta.keyHint}
                </p>
              )}
              <div style={{ display: 'flex', gap: '0.75rem' }}>
                <input
                  type={activeAuthType === 'web_cookie' ? 'text' : 'password'}
                  className="text-input"
                  placeholder={`Paste your ${activeMeta.keyLabel || 'API Key'}...`}
                  value={newKeyInput}
                  onChange={e => setNewKeyInput(e.target.value)}
                  style={{ flex: 1 }}
                  onKeyDown={e => { if (e.key === 'Enter') handleAddKey(); }}
                />
                <button className="btn-secondary" onClick={handleAddKey} disabled={!newKeyInput}>
                  Save
                </button>
              </div>
              {activeProvider.url && (
                <p style={{ marginTop: '0.6rem', fontSize: '0.76rem', color: 'var(--text-muted)' }}>
                  <a href={activeProvider.url} target="_blank" rel="noreferrer" style={{ color: 'var(--accent-primary)', textDecoration: 'none' }}>
                    {activeAuthType === 'web_cookie' ? 'Open provider site ↗' : activeAuthType === 'oauth' ? 'Generate token ↗' : 'Get API key ↗'}
                  </a>
                </p>
              )}
            </div>

          </div>
        </div>
      )}
    </div>
  );
}
