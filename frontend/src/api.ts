/** Dev keeps 20128; production Vite builds use 30128. */
export const API_PORT = import.meta.env.DEV ? '20128' : '30128'

export const API_BASE = `http://localhost:${API_PORT}`
