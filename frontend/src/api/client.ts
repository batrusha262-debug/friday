const BASE = (import.meta.env.VITE_API_URL ?? '').replace(/\/$/, '')

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const token = localStorage.getItem('token')
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (token) headers['Authorization'] = `Bearer ${token}`

  const res = await fetch(`${BASE}${path}`, {
    headers,
    ...init,
  })

  if (res.status === 204) return undefined as T

  const text = await res.text()
  let data: unknown

  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = text
    }
  }

  if (!res.ok) {
    if (typeof data === 'string') {
      throw new Error(data)
    }

    if (typeof data === 'object' && data !== null && 'message' in data) {
      const message = (data as { message?: unknown }).message
      if (typeof message === 'string' && message) {
        throw new Error(message)
      }
    }

    throw new Error(`HTTP ${res.status}`)
  }

  return data as T
}

export const api = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, body: unknown) =>
    request<T>(path, { method: 'POST', body: JSON.stringify(body) }),
  put: <T>(path: string, body: unknown) =>
    request<T>(path, { method: 'PUT', body: JSON.stringify(body) }),
  patch: <T>(path: string, body: unknown) =>
    request<T>(path, { method: 'PATCH', body: JSON.stringify(body) }),
  delete: (path: string) => request<void>(path, { method: 'DELETE' }),
}
