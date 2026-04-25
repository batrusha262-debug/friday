import { useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createPack, createRound, createCategory, createGame, addTeam, startGame } from '../api'

const PRICE_SCALES = {
  '100–500': [100, 200, 300, 400, 500],
  '100–1000': [100, 200, 300, 400, 500, 600, 700, 800, 900, 1000],
} as const

type ScaleKey = keyof typeof PRICE_SCALES

export default function CreateGamePage() {
  const navigate = useNavigate()
  const location = useLocation()
  const qc = useQueryClient()

  const state = location.state as { packTitle?: string } | null
  const [title, setTitle] = useState(state?.packTitle ?? '')
  const [catCount, setCatCount] = useState(5)
  const [catNames, setCatNames] = useState<string[]>(
    Array.from({ length: 5 }, (_, i) => `Категория ${i + 1}`),
  )
  const [scale, setScale] = useState<ScaleKey>('100–500')
  const [playerInput, setPlayerInput] = useState('')
  const [players, setPlayers] = useState<string[]>([])
  const [error, setError] = useState('')

  function handleCatCount(n: number) {
    setCatCount(n)
    setCatNames(prev => {
      const next = [...prev]

      while (next.length < n) next.push(`Категория ${next.length + 1}`)

      return next.slice(0, n)
    })
  }

  function addPlayer() {
    const p = playerInput.trim()

    if (!p || players.includes(p)) return

    setPlayers(prev => [...prev, p])
    setPlayerInput('')
  }

  function removePlayer(name: string) {
    setPlayers(prev => prev.filter(p => p !== name))
  }

  const { mutate, isPending } = useMutation({
    mutationFn: async () => {
      const userId = localStorage.getItem('userId')!

      // 1. Create pack
      const pack = await createPack(title.trim(), userId)

      // 2. Create round
      const round = await createRound(pack.id, 'Раунд 1', 'standard')

      // 3. Create categories
      for (let i = 0; i < catCount; i++) {
        await createCategory(round.id, catNames[i] || `Категория ${i + 1}`)
      }

      // 4. Create game
      const game = await createGame(pack.id, userId)

      // 5. Add teams
      for (const name of players) {
        await addTeam(game.id, name)
      }

      // 6. Start game
      await startGame(game.id)

      // 7. Persist game reference
      localStorage.setItem(`pack:${pack.id}:gameId`, game.id)
      localStorage.setItem(`game:${game.id}:packId`, pack.id)
      localStorage.setItem(`game:${game.id}:status`, 'active')
      localStorage.setItem(`game:${game.id}:scale`, JSON.stringify(PRICE_SCALES[scale]))

      return game
    },
    onSuccess: game => {
      qc.invalidateQueries({ queryKey: ['packs'] })
      navigate(`/game/${game.id}`, { replace: true })
    },
    onError: err => {
      setError(err instanceof Error ? err.message : 'Ошибка создания игры')
    },
  })

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()

    if (!title.trim()) { setError('Введите название'); return }
    if (players.length === 0) { setError('Добавьте хотя бы одного игрока'); return }

    setError('')
    mutate()
  }

  return (
    <div className="page">
      <div className="tgh">
        <button className="tgh-back" onClick={() => navigate(-1)}>
          <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="2">
            <path d="M12 5l-7 5 7 5" />
          </svg>
        </button>
        <span className="tgh-title">Новая игра</span>
      </div>

      <div className="page-body">
        <form onSubmit={handleSubmit} style={{ padding: '10px 0 20px' }}>
          <div className="wcard" style={{ margin: '0 8px 8px' }}>
            <span className="tlabel">Название игры</span>
            <input
              className="tinput"
              placeholder="Например: «Итоги квартала»"
              value={title}
              onChange={e => setTitle(e.target.value)}
            />

            <span className="tlabel">Категории ({catCount})</span>
            <div className="pill-row" style={{ gridTemplateColumns: 'repeat(5, minmax(0, 1fr))', marginBottom: 8 }}>
              {[3, 4, 5, 6, 7].map(n => (
                <button
                  key={n}
                  type="button"
                  className={`ppill${catCount === n ? ' on' : ''}`}
                  onClick={() => handleCatCount(n)}
                >
                  {n}
                </button>
              ))}
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: 4, marginBottom: 8 }}>
              {catNames.map((name, i) => (
                <input
                  key={i}
                  className="tinput"
                  style={{ marginBottom: 0 }}
                  placeholder={`Категория ${i + 1}`}
                  value={name}
                  onChange={e => {
                    const next = [...catNames]

                    next[i] = e.target.value
                    setCatNames(next)
                  }}
                />
              ))}
            </div>

            <span className="tlabel">Стоимость вопросов</span>
            <div className="pill-row" style={{ gridTemplateColumns: 'repeat(2, minmax(0, 1fr))' }}>
              {(Object.keys(PRICE_SCALES) as ScaleKey[]).map(k => (
                <button
                  key={k}
                  type="button"
                  className={`ppill${scale === k ? ' on' : ''}`}
                  onClick={() => setScale(k)}
                >
                  {k}
                </button>
              ))}
            </div>

            <span className="tlabel">Игроки</span>
            <div style={{ display: 'flex', gap: 6, marginBottom: 8 }}>
              <input
                className="tinput"
                style={{ margin: 0, flex: 1 }}
                placeholder="Имя участника"
                value={playerInput}
                onChange={e => setPlayerInput(e.target.value)}
                onKeyDown={e => e.key === 'Enter' && (e.preventDefault(), addPlayer())}
              />
              <button
                type="button"
                onClick={addPlayer}
                style={{
                  background: '#1a1a1a',
                  color: '#fff',
                  border: 'none',
                  borderRadius: 8,
                  padding: '9px 14px',
                  fontSize: 14,
                  cursor: 'pointer',
                }}
              >
                +
              </button>
            </div>

            {players.length > 0 && (
              <div style={{ display: 'flex', gap: 5, flexWrap: 'wrap' }}>
                {players.map(p => (
                  <span key={p} className="chip">
                    {p}
                    <button className="chip-remove" type="button" onClick={() => removePlayer(p)}>×</button>
                  </span>
                ))}
              </div>
            )}

            {error && (
              <div style={{ fontSize: 12, color: '#c00', marginTop: 8 }}>{error}</div>
            )}
          </div>

          <div style={{ padding: '0 8px' }}>
            <button className="tbtn" type="submit" disabled={isPending}>
              {isPending ? 'Создаём…' : 'Создать и добавить вопросы →'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
