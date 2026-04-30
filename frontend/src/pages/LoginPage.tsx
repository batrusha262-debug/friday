import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { requestCode, verifyCode, guestLogin } from '../api'

type Mode = 'choose' | 'email' | 'code' | 'guest'

function saveSession(token: string, user: { id: string; username: string; role: string }) {
  localStorage.setItem('token', token)
  localStorage.setItem('userId', user.id)
  localStorage.setItem('userName', user.username)
  localStorage.setItem('userRole', user.role)
}

export default function LoginPage() {
  const navigate = useNavigate()
  const [mode, setMode] = useState<Mode>('choose')
  const [email, setEmail] = useState('')
  const [code, setCode] = useState('')
  const [name, setName] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function handleRequestCode(e: React.FormEvent) {
    e.preventDefault()
    const trimmed = email.trim().toLowerCase()
    if (!trimmed) return
    setLoading(true)
    setError('')
    try {
      await requestCode(trimmed)
      setEmail(trimmed)
      setMode('code')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка отправки кода')
    } finally {
      setLoading(false)
    }
  }

  async function handleVerifyCode(e: React.FormEvent) {
    e.preventDefault()
    const trimmed = code.trim()
    if (!trimmed) return
    setLoading(true)
    setError('')
    try {
      const session = await verifyCode(email, trimmed)
      saveSession(session.token, session.user)
      navigate('/', { replace: true })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Неверный или истёкший код')
    } finally {
      setLoading(false)
    }
  }

  async function handleGuestLogin(e: React.FormEvent) {
    e.preventDefault()
    const trimmed = name.trim()
    if (!trimmed) return
    setLoading(true)
    setError('')
    try {
      const session = await guestLogin(trimmed)
      saveSession(session.token, session.user)
      navigate('/', { replace: true })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка входа')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="page">
      <div className="tgh">
        {mode !== 'choose' && (
          <button className="tgh-back" onClick={() => { setMode('choose'); setError('') }}>
            <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="2">
              <path d="M12 5l-7 5 7 5" />
            </svg>
          </button>
        )}
        <span className="tgh-title">Своя игра</span>
      </div>

      <div className="page-body center">
        <div style={{ width: '100%', padding: '0 16px' }}>

          {/* Choose mode */}
          {mode === 'choose' && (
            <div className="wcard">
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
                <div style={{ fontSize: 18, fontWeight: 600, color: '#1a1a1a', marginBottom: 4 }}>
                  Добро пожаловать
                </div>
                <div style={{ fontSize: 13, color: '#999' }}>Выберите способ входа</div>
              </div>

              <button
                className="tbtn"
                style={{ marginBottom: 8 }}
                onClick={() => setMode('email')}
              >
                Войти как администратор
              </button>
              <button
                className="tbtn-ghost"
                onClick={() => setMode('guest')}
              >
                Войти как гость
              </button>
            </div>
          )}

          {/* Admin: enter email */}
          {mode === 'email' && (
            <form onSubmit={handleRequestCode}>
              <div className="wcard" style={{ marginBottom: 12 }}>
                <div style={{ fontSize: 15, fontWeight: 600, color: '#1a1a1a', marginBottom: 4 }}>
                  Вход для администратора
                </div>
                <div style={{ fontSize: 13, color: '#999', marginBottom: 16 }}>
                  Введите email — пришлём код подтверждения
                </div>
                <span className="tlabel">Email</span>
                <input
                  className="tinput"
                  type="email"
                  placeholder="your@email.com"
                  value={email}
                  onChange={e => setEmail(e.target.value)}
                  autoFocus
                />
                {error && (
                  <div style={{ fontSize: 12, color: '#c00', marginBottom: 8 }}>{error}</div>
                )}
              </div>
              <button className="tbtn" type="submit" disabled={loading || !email.trim()}>
                {loading ? 'Отправляем…' : 'Получить код →'}
              </button>
            </form>
          )}

          {/* Admin: enter OTP code */}
          {mode === 'code' && (
            <form onSubmit={handleVerifyCode}>
              <div className="wcard" style={{ marginBottom: 12 }}>
                <div style={{ fontSize: 15, fontWeight: 600, color: '#1a1a1a', marginBottom: 4 }}>
                  Введите код
                </div>
                <div style={{ fontSize: 13, color: '#999', marginBottom: 16 }}>
                  Код отправлен на <strong>{email}</strong>
                </div>
                <span className="tlabel">Код подтверждения</span>
                <input
                  className="tinput"
                  type="text"
                  inputMode="numeric"
                  placeholder="000000"
                  maxLength={6}
                  value={code}
                  onChange={e => setCode(e.target.value.replace(/\D/g, ''))}
                  autoFocus
                />
                {error && (
                  <div style={{ fontSize: 12, color: '#c00', marginBottom: 8 }}>{error}</div>
                )}
              </div>
              <button className="tbtn" type="submit" disabled={loading || code.length < 6}>
                {loading ? 'Проверяем…' : 'Войти →'}
              </button>
              <button
                type="button"
                className="tbtn-ghost"
                style={{ marginTop: 8 }}
                disabled={loading}
                onClick={() => { setMode('email'); setCode(''); setError('') }}
              >
                Отправить код снова
              </button>
            </form>
          )}

          {/* Guest: enter name */}
          {mode === 'guest' && (
            <form onSubmit={handleGuestLogin}>
              <div className="wcard" style={{ marginBottom: 12 }}>
                <div style={{ fontSize: 15, fontWeight: 600, color: '#1a1a1a', marginBottom: 4 }}>
                  Вход как гость
                </div>
                <div style={{ fontSize: 13, color: '#999', marginBottom: 16 }}>
                  Введите имя — чтобы было понятно, кто играет
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
                {loading ? 'Входим…' : 'Войти →'}
              </button>
            </form>
          )}

        </div>
      </div>
    </div>
  )
}
