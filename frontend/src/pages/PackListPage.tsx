import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { listPacks, logout, findGameByCode, getGameByPack } from '../api'
import { useAuth } from '../App'
import type { Pack } from '../api/types'

function packBadge(pack: Pack, isAdmin: boolean) {
  if (!isAdmin) {
    return { label: 'Открыта', cls: 'badge-ready' }
  }

  const stored = localStorage.getItem(`pack:${pack.id}:gameId`)

  if (stored) {
    const status = localStorage.getItem(`game:${stored}:status`)

    if (status === 'active') return { label: 'В игре', cls: 'badge-active' }
    if (status === 'finished') return { label: 'Завершена', cls: 'badge-ready' }
  }

  return { label: 'Черновик', cls: 'badge-draft' }
}

const ICONS = ['#1a1a1a', '#333', '#555', '#666', '#777', '#888']

export default function PackListPage() {
  const navigate = useNavigate()
  const { role } = useAuth()
  const isAdmin = role === 'admin'
  const [joinCode, setJoinCode] = useState('')
  const [joinError, setJoinError] = useState('')
  const [joining, setJoining] = useState(false)
  const [joinModal, setJoinModal] = useState(false)

  async function handleJoin(e: React.FormEvent) {
    e.preventDefault()
    const code = joinCode.trim()
    if (!code) return
    setJoining(true)
    setJoinError('')
    try {
      const game = await findGameByCode(code)
      setJoinModal(false)
      setJoinCode('')
      navigate(`/game/${game.id}`)
    } catch {
      setJoinError('Игра не найдена. Проверьте код.')
    } finally {
      setJoining(false)
    }
  }

  async function handleLogout() {
    await logout().catch(() => {})
    localStorage.clear()
    navigate('/login', { replace: true })
  }
  const { data: packs, isLoading } = useQuery({
    queryKey: ['packs'],
    queryFn: listPacks,
  })

  async function openPack(pack: Pack) {
    const cached = localStorage.getItem(`pack:${pack.id}:gameId`)

    if (cached) {
      navigate(`/game/${cached}`)

      return
    }

    if (isAdmin) {
      try {
        const game = await getGameByPack(pack.id)
        localStorage.setItem(`pack:${pack.id}:gameId`, game.id)
        localStorage.setItem(`game:${game.id}:status`, game.status)
        navigate(`/game/${game.id}`)
      } catch {
        // no game found — create a new one
        navigate(`/game/create`, { state: { packId: pack.id, packTitle: pack.title } })
      }
    } else {
      setJoinCode('')
      setJoinError('')
      setJoinModal(true)
    }
  }

  return (
    <div className="page">
      <div className="tgh">
        <span style={{ position: 'absolute', left: 0, right: 0, textAlign: 'center', letterSpacing: 3, fontSize: 12, fontWeight: 400, color: '#fff', textTransform: 'uppercase', pointerEvents: 'none' }}>KODE</span>
        <div style={{ display: 'flex', gap: 4 }}>
          {isAdmin && (
            <button className="tgh-action" onClick={() => navigate('/game/create')}>
              <svg width="18" height="18" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="1.8">
                <path d="M10 4v12M4 10h12" />
              </svg>
            </button>
          )}
          <button className="tgh-action" onClick={handleLogout} title="Выйти">
            <svg width="18" height="18" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="1.8">
              <path d="M13 10H4M10 7l3 3-3 3" />
              <path d="M8 5H5a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h3" />
            </svg>
          </button>
        </div>
      </div>

      <div className="page-body">
        {!isAdmin && (
          <div style={{ padding: '12px 8px 4px' }}>
            <form onSubmit={handleJoin}>
              <div className="wcard" style={{ padding: '12px 12px 8px' }}>
                <span className="tlabel">Код игры</span>
                <div style={{ display: 'flex', gap: 6 }}>
                  <input
                    className="tinput"
                    style={{ margin: 0, flex: 1, textTransform: 'uppercase', letterSpacing: 2, fontFamily: 'monospace' }}
                    placeholder="ABCD1234"
                    value={joinCode}
                    maxLength={8}
                    onChange={e => setJoinCode(e.target.value.replace(/[^a-fA-F0-9]/g, '').toUpperCase())}
                  />
                  <button
                    type="submit"
                    disabled={joining || joinCode.length < 4}
                    style={{
                      background: '#1a1a1a',
                      color: '#fff',
                      border: 'none',
                      borderRadius: 8,
                      padding: '9px 16px',
                      fontSize: 14,
                      cursor: 'pointer',
                      opacity: joinCode.length < 4 ? 0.4 : 1,
                    }}
                  >
                    {joining ? '…' : 'Войти'}
                  </button>
                </div>
                {joinError && (
                  <div style={{ fontSize: 12, color: '#c00', marginTop: 6 }}>{joinError}</div>
                )}
              </div>
            </form>
          </div>
        )}

        {isLoading && (
          <div className="center" style={{ padding: 40 }}>
            <div className="spinner" />
          </div>
        )}

        {!isLoading && (!packs || packs.length === 0) && (
          <div className="empty">
            <div className="empty-icon">🎲</div>
            <div>Игр пока нет</div>
            <div style={{ fontSize: 12, marginTop: 4 }}>Нажмите + чтобы создать</div>
          </div>
        )}

        {joinModal && (
          <div
            style={{
              position: 'fixed',
              inset: 0,
              background: 'rgba(0,0,0,0.5)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              zIndex: 100,
              padding: 16,
            }}
            onClick={() => setJoinModal(false)}
          >
            <div
              style={{ background: '#fff', borderRadius: 16, padding: 20, width: '100%', maxWidth: 340 }}
              onClick={e => e.stopPropagation()}
            >
              <div style={{ fontSize: 16, fontWeight: 600, marginBottom: 4 }}>Войти в игру</div>
              <div style={{ fontSize: 13, color: '#999', marginBottom: 16 }}>
                Введите код, который показывает ведущий
              </div>
              <form onSubmit={handleJoin}>
                <input
                  className="tinput"
                  style={{ textTransform: 'uppercase', letterSpacing: 2, fontFamily: 'monospace' }}
                  placeholder="ABCD1234"
                  value={joinCode}
                  maxLength={8}
                  autoFocus
                  onChange={e => setJoinCode(e.target.value.replace(/[^a-fA-F0-9]/g, '').toUpperCase())}
                />
                {joinError && (
                  <div style={{ fontSize: 12, color: '#c00', marginBottom: 8 }}>{joinError}</div>
                )}
                <div style={{ display: 'flex', gap: 8 }}>
                  <button
                    type="button"
                    className="tbtn-ghost"
                    style={{ flex: 1 }}
                    onClick={() => setJoinModal(false)}
                  >
                    Отмена
                  </button>
                  <button
                    type="submit"
                    className="tbtn"
                    style={{ flex: 1 }}
                    disabled={joining || joinCode.length < 4}
                  >
                    {joining ? '…' : 'Войти'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {packs && packs.length > 0 && (
          <div style={{ paddingTop: 6 }}>
            <div className="row-between px-14 py-6">
              <span className="text-sm text-mid">{packs.length} {pluralize(packs.length, 'игра', 'игры', 'игр')}</span>
            </div>

            {packs.map((pack, i) => {
              const badge = packBadge(pack, isAdmin)
              const iconBg = ICONS[i % ICONS.length]

              return (
                <div key={pack.id}>
                  <div
                    className="wcard wcard-row"
                    style={{ marginBottom: 0, cursor: 'pointer' }}
                    onClick={() => openPack(pack)}
                  >
                    <div className="pack-icon" style={{ background: iconBg }}>
                      <svg width="18" height="18" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="1.5">
                        <rect x="3" y="5" width="14" height="10" rx="2" />
                        <path d="M7 9h6M7 12h3" />
                      </svg>
                    </div>

                    <div style={{ flex: 1, minWidth: 0 }}>
                      <div style={{ fontSize: 14, fontWeight: 500, color: '#1a1a1a' }}>
                        {pack.title}
                      </div>
                      <div className="text-sm text-mid mt-4">
                        {new Date(pack.created_at).toLocaleDateString('ru', {
                          day: 'numeric',
                          month: 'short',
                          year: 'numeric',
                        })}
                      </div>
                    </div>

                    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: 4 }}>
                      <span className={`badge ${badge.cls}`}>{badge.label}</span>
                      <svg width="14" height="14" viewBox="0 0 20 20" fill="none" stroke="#ccc" strokeWidth="2">
                        <path d="M8 5l6 5-6 5" />
                      </svg>
                    </div>
                  </div>
                  {i < packs.length - 1 && <div className="divider" />}
                </div>
              )
            })}

            {isAdmin && (
              <div style={{ padding: '4px 8px 0' }}>
                <button className="tbtn" onClick={() => navigate('/game/create')}>
                  + Создать игру
                </button>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

function pluralize(n: number, one: string, few: string, many: string) {
  const mod10 = n % 10
  const mod100 = n % 100

  if (mod100 >= 11 && mod100 <= 19) return many
  if (mod10 === 1) return one
  if (mod10 >= 2 && mod10 <= 4) return few

  return many
}
