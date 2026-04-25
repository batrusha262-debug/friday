import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { createUser } from '../api'

export default function SetupPage() {
  const navigate = useNavigate()
  const [name, setName] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()

    const trimmed = name.trim()

    if (!trimmed) return

    setLoading(true)
    setError('')

    try {
      const user = await createUser(trimmed)

      localStorage.setItem('userId', user.id)
      localStorage.setItem('userName', user.username)
      navigate('/', { replace: true })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка создания профиля')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="page">
      <div className="tgh">
        <span className="tgh-title">Своя игра</span>
      </div>
      <div className="page-body center">
        <form onSubmit={handleSubmit} style={{ width: '100%', padding: '0 16px' }}>
          <div className="wcard" style={{ margin: '0 0 12px' }}>
            <div style={{ textAlign: 'center', marginBottom: 20 }}>
              <div
                style={{
                  width: 56,
                  height: 56,
                  background: '#1a1a1a',
                  borderRadius: 14,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  margin: '0 auto 12px',
                }}
              >
                <svg width="28" height="28" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="1.5">
                  <rect x="3" y="5" width="14" height="10" rx="2" />
                  <path d="M7 9h6M7 12h3" />
                </svg>
              </div>
              <div style={{ fontSize: 18, fontWeight: 600, color: '#1a1a1a', marginBottom: 4 }}>Добро пожаловать</div>
              <div style={{ fontSize: 13, color: '#999' }}>Введите ваше имя, чтобы начать</div>
            </div>
            <span className="tlabel">Ваше имя</span>
            <input
              className="tinput"
              placeholder="Например: Алекс"
              value={name}
              onChange={e => setName(e.target.value)}
              autoFocus
            />
            {error && (
              <div style={{ fontSize: 12, color: '#c00', marginBottom: 8 }}>{error}</div>
            )}
          </div>
          <button className="tbtn" type="submit" disabled={loading || !name.trim()}>
            {loading ? 'Создаём профиль…' : 'Начать →'}
          </button>
        </form>
      </div>
    </div>
  )
}
