import { useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { listPacks } from '../api'
import type { Pack } from '../api/types'

function packBadge(pack: Pack) {
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
  const { data: packs, isLoading } = useQuery({
    queryKey: ['packs'],
    queryFn: listPacks,
  })

  function openPack(pack: Pack) {
    const gameId = localStorage.getItem(`pack:${pack.id}:gameId`)

    if (gameId) {
      navigate(`/game/${gameId}`)
    } else {
      navigate(`/game/create`, { state: { packId: pack.id, packTitle: pack.title } })
    }
  }

  return (
    <div className="page">
      <div className="tgh">
        <span className="tgh-title">Своя игра</span>
        <button className="tgh-action" onClick={() => navigate('/game/create')}>
          <svg width="18" height="18" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="1.8">
            <path d="M10 4v12M4 10h12" />
          </svg>
        </button>
      </div>

      <div className="page-body">
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

        {packs && packs.length > 0 && (
          <div style={{ paddingTop: 6 }}>
            <div className="row-between px-14 py-6">
              <span className="text-sm text-mid">{packs.length} {pluralize(packs.length, 'игра', 'игры', 'игр')}</span>
              <span className="text-sm text-mid">Сортировка</span>
            </div>

            {packs.map((pack, i) => {
              const badge = packBadge(pack)
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

            <div style={{ padding: '4px 8px 0' }}>
              <button className="tbtn" onClick={() => navigate('/game/create')}>
                + Создать игру
              </button>
            </div>
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
