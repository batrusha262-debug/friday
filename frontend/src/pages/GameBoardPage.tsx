import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate, useParams } from 'react-router-dom'
import { getGame, getBoard, listRounds, listCategories, listQuestions, startGame } from '../api'
import { useGameEvents } from '../hooks/useGameEvents'
import type { Game, GameBoard, GameTeam, Question } from '../api/types'

interface BoardCell {
  question: Question | null
  answered: boolean
  answeredBy: string | null
}

function useBoardData(gameId: string) {
  const [liveGame, setLiveGame] = useState<Game | null>(null)
  const [liveBoard, setLiveBoard] = useState<GameBoard | null>(null)

  useGameEvents(gameId, (state) => {
    setLiveGame(state.game)
    setLiveBoard(state.board)
  })

  const gameQuery = useQuery({ queryKey: ['game', gameId], queryFn: () => getGame(gameId) })
  const boardQuery = useQuery({ queryKey: ['board', gameId], queryFn: () => getBoard(gameId) })

  const effectiveGame = liveGame ?? gameQuery.data ?? null
  const effectiveTeams: GameTeam[] = (liveBoard ?? boardQuery.data)?.teams ?? []
  const effectiveStates = (liveBoard ?? boardQuery.data)?.states ?? []

  const packId = effectiveGame?.pack_id
  const rounds = useQuery({
    queryKey: ['rounds', packId],
    queryFn: () => listRounds(packId!),
    enabled: !!packId,
  })

  const roundId = rounds.data?.[0]?.id
  const categories = useQuery({
    queryKey: ['categories', roundId],
    queryFn: () => listCategories(roundId!),
    enabled: !!roundId,
  })

  const catList = categories.data ?? []
  const questionsQueries = useQuery({
    queryKey: ['questions-all', catList.map(c => c.id).join(',')],
    queryFn: async () => {
      const results = await Promise.all(catList.map(c => listQuestions(c.id)))

      return Object.fromEntries(catList.map((c, i) => [c.id, results[i]]))
    },
    enabled: catList.length > 0,
  })

  const prices = JSON.parse(localStorage.getItem(`game:${gameId}:scale`) ?? '[100,200,300,400,500]') as number[]
  const answeredIds = new Set(effectiveStates.map(s => s.question_id))
  const stateByQuestion = Object.fromEntries(effectiveStates.map(s => [s.question_id, s]))

  const allQuestions = questionsQueries.data ?? {}

  const grid: BoardCell[][] = prices.map(price =>
    catList.map(cat => {
      const qs = allQuestions[cat.id] ?? []
      const q = qs.find(x => x.price === price) ?? null
      const answered = q ? answeredIds.has(q.id) : false
      const answeredBy = q ? (stateByQuestion[q.id]?.answered_by ?? null) : null

      return { question: q, answered, answeredBy }
    }),
  )

  return {
    loading: gameQuery.isLoading || boardQuery.isLoading || rounds.isLoading || categories.isLoading,
    game: effectiveGame,
    teams: effectiveTeams,
    categories: catList,
    prices,
    grid,
  }
}

export default function GameBoardPage() {
  const { gameId } = useParams<{ gameId: string }>()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { loading, game, teams, categories, prices, grid } = useBoardData(gameId!)

  const { mutate: doStart, isPending: isStarting } = useMutation({
    mutationFn: () => startGame(gameId!),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['game', gameId] })
    },
  })

  if (loading) {
    return (
      <div className="page">
        <div className="tgh"><span className="tgh-title">Загрузка…</span></div>
        <div className="center"><div className="spinner" /></div>
      </div>
    )
  }

  // Lobby: game is waiting for start
  if (game?.status === 'waiting') {
    return (
      <div className="page">
        <div className="tgh">
          <button className="tgh-back" onClick={() => navigate('/')}>
            <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="2">
              <path d="M12 5l-7 5 7 5" />
            </svg>
          </button>
          <span className="tgh-title">Ожидание игроков</span>
        </div>
        <div className="page-body" style={{ padding: 16 }}>
          {teams.length > 0 && (
            <>
              <div className="text-sm text-mid mb-8">Команды ({teams.length}):</div>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 4, marginBottom: 16 }}>
                {teams.map(t => (
                  <div
                    key={t.id}
                    style={{
                      background: '#f5f5f5',
                      borderRadius: 8,
                      padding: '8px 12px',
                      fontSize: 14,
                    }}
                  >
                    {t.name}
                  </div>
                ))}
              </div>
            </>
          )}
          <button className="tbtn" onClick={() => doStart()} disabled={isStarting || teams.length === 0}>
            {isStarting ? 'Запускаем…' : 'Начать игру'}
          </button>
          {teams.length === 0 && (
            <div style={{ fontSize: 12, color: '#999', marginTop: 8, textAlign: 'center' }}>
              Добавьте хотя бы одну команду
            </div>
          )}
        </div>
      </div>
    )
  }

  const sorted = [...teams].sort((a, b) => b.score - a.score)
  const maxScore = sorted[0]?.score ?? 0
  const currentPicker = teams.find(t => t.id === game?.current_picker_id)
  const colCount = categories.length || 1

  return (
    <div className="page">
      <div className="tgh">
        <button className="tgh-back" onClick={() => navigate('/')}>
          <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="2">
            <path d="M12 5l-7 5 7 5" />
          </svg>
        </button>
        <span className="tgh-title">{game?.pack_id ? '' : 'Игра'}</span>
        <button
          className="tgh-action"
          onClick={() => navigate(`/game/${gameId}/question/add`)}
        >
          <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="2">
            <path d="M10 4v12M4 10h12" />
          </svg>
        </button>
      </div>

      {currentPicker && (
        <div
          style={{
            background: '#1a1a1a',
            color: '#fff',
            textAlign: 'center',
            padding: '6px 12px',
            fontSize: 13,
          }}
        >
          Выбирает: <strong>{currentPicker.name}</strong>
        </div>
      )}

      <div className="page-body">
        <div
          className="board-grid"
          style={{ gridTemplateColumns: `repeat(${colCount}, minmax(0, 1fr))` }}
        >
          {categories.map(cat => (
            <div key={cat.id} className="qcell cat">
              {cat.name}
            </div>
          ))}

          {prices.map((price, pi) =>
            categories.map((cat, ci) => {
              const cell = grid[pi]?.[ci]

              if (!cell) return <div key={`${pi}-${ci}`} className="qcell empty">+</div>

              if (cell.answered) {
                return <div key={`${pi}-${ci}`} className="qcell used">{price}</div>
              }

              if (!cell.question) {
                return (
                  <button
                    key={`${pi}-${ci}`}
                    className="qcell empty"
                    onClick={() => navigate(`/game/${gameId}/question/add`, { state: { categoryId: cat.id, price } })}
                  >
                    +
                  </button>
                )
              }

              return (
                <button
                  key={`${pi}-${ci}`}
                  className="qcell"
                  onClick={() => navigate(`/game/${gameId}/question/${cell.question!.id}`)}
                >
                  {price}
                </button>
              )
            }),
          )}
        </div>

        {teams.length > 0 && (
          <div className="score-bar">
            <div className="text-sm text-mid mb-8">Счёт</div>
            <div style={{ display: 'flex', gap: 5 }}>
              {sorted.map(team => (
                <div
                  key={team.id}
                  className={`score-card ${team.score === maxScore && maxScore > 0 ? 'leading' : 'other'}`}
                >
                  <div
                    style={{
                      fontSize: 10,
                      color: team.score === maxScore && maxScore > 0 ? '#aaa' : '#999',
                    }}
                  >
                    {team.name}
                  </div>
                  <div
                    style={{
                      fontSize: 16,
                      fontWeight: 500,
                      color: team.score === maxScore && maxScore > 0 ? '#fff' : '#333',
                    }}
                  >
                    {team.score}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
