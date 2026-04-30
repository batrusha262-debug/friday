import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate, useParams } from 'react-router-dom'
import {
  getGame,
  getBoard,
  listRounds,
  listCategories,
  listQuestions,
  startGame,
  deleteGame,
  deletePack,
  validateClaim,
  setGameOpen,
} from '../api'
import { useGameEvents } from '../hooks/useGameEvents'
import { useAuth } from '../App'
import type { AnswerClaim, Game, GameBoard, GameTeam, Question } from '../api/types'

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
  const pendingClaims: AnswerClaim[] = (liveBoard ?? boardQuery.data)?.pending_claims ?? []

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

  // flat map of questionId -> Question for claim lookups
  const questionById: Record<string, Question> = {}
  for (const qs of Object.values(allQuestions)) {
    for (const q of qs) {
      questionById[q.id] = q
    }
  }

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
    pendingClaims,
    questionById,
  }
}

function TeamsPanel({ teams, currentPickerId }: { teams: GameTeam[]; currentPickerId?: string }) {
  const sorted = [...teams].sort((a, b) => b.score - a.score)
  const maxScore = sorted[0]?.score ?? 0

  if (teams.length === 0) return null

  return (
    <div className="score-bar">
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 8,
        }}
      >
        <span className="text-sm text-mid">Команды</span>
        {currentPickerId && (
          <span className="text-sm text-mid">
            Ход: <strong style={{ color: '#1a1a1a' }}>
              {teams.find(t => t.id === currentPickerId)?.name}
            </strong>
          </span>
        )}
      </div>
      <div style={{ display: 'flex', gap: 5 }}>
        {sorted.map(team => {
          const isPicker = team.id === currentPickerId
          const isLeading = team.score === maxScore && maxScore > 0

          return (
            <div
              key={team.id}
              className={`score-card ${isLeading ? 'leading' : 'other'}`}
              style={{ outline: isPicker ? '2px solid #f0a500' : undefined, outlineOffset: 2 }}
            >
              <div style={{ fontSize: 10, color: isLeading ? '#aaa' : '#999' }}>
                {team.name}
              </div>
              <div style={{ fontSize: 16, fontWeight: 500, color: isLeading ? '#fff' : '#333' }}>
                {team.score}
              </div>
              {isPicker && (
                <div style={{ fontSize: 9, color: '#f0a500', marginTop: 1 }}>▶ ход</div>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}

export default function GameBoardPage() {
  const { gameId } = useParams<{ gameId: string }>()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { role } = useAuth()
  const isAdmin = role === 'admin'
  const { loading, game, teams, categories, prices, grid, pendingClaims, questionById } = useBoardData(gameId!)

  const { mutate: doStart, isPending: isStarting } = useMutation({
    mutationFn: () => startGame(gameId!),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['game', gameId] })
    },
  })

  const { mutateAsync: doDelete, isPending: isDeleting } = useMutation({
    mutationFn: () => deleteGame(gameId!),
  })

  const { mutate: doToggleOpen, isPending: isTogglingOpen } = useMutation({
    mutationFn: (open: boolean) => setGameOpen(gameId!, open),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['game', gameId] })
    },
  })

  const { mutate: doValidate } = useMutation({
    mutationFn: ({ claimId, approved }: { claimId: string; approved: boolean }) =>
      validateClaim(claimId, approved),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['board', gameId] })
    },
  })

  async function handleDelete() {
    if (!window.confirm('Удалить игру? Это действие нельзя отменить.')) return
    const packId = game?.pack_id
    try {
      await doDelete()
      if (packId) await deletePack(packId)
    } catch {
      return
    }
    if (packId) localStorage.removeItem(`pack:${packId}:gameId`)
    localStorage.removeItem(`game:${gameId}:status`)
    localStorage.removeItem(`game:${gameId}:scale`)
    qc.invalidateQueries({ queryKey: ['packs'] })
    navigate('/')
  }

  if (loading) {
    return (
      <div className="page">
        <div className="tgh"><span className="tgh-title">Загрузка…</span></div>
        <div className="center"><div className="spinner" /></div>
      </div>
    )
  }

  // Pending claims banner (admin only)
  const claimsBanner = isAdmin && pendingClaims.length > 0 && (
    <div style={{ background: '#fffbe6', borderBottom: '1px solid #ffe58f', padding: '8px 12px' }}>
      <div style={{ fontSize: 11, color: '#996600', textTransform: 'uppercase', letterSpacing: 0.5, marginBottom: 6 }}>
        Ожидают подтверждения ({pendingClaims.length})
      </div>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
        {pendingClaims.map(c => {
          const team = teams.find(t => t.id === c.team_id)
          const question = questionById[c.question_id]

          return (
            <div
              key={c.id}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: 8,
                background: '#fff',
                borderRadius: 8,
                padding: '6px 10px',
                border: '0.5px solid #e0e0e0',
              }}
            >
              <div style={{ flex: 1, fontSize: 13 }}>
                <strong>{team?.name ?? '…'}</strong>
                {question && (
                  <span style={{ color: '#999', marginLeft: 6 }}>— {question.price} очков</span>
                )}
              </div>
              <button
                onClick={() => doValidate({ claimId: c.id, approved: true })}
                style={{
                  background: '#1a1a1a',
                  color: '#fff',
                  border: 'none',
                  borderRadius: 6,
                  padding: '4px 10px',
                  fontSize: 12,
                  cursor: 'pointer',
                  fontFamily: 'inherit',
                }}
              >
                Засчитать
              </button>
              <button
                onClick={() => doValidate({ claimId: c.id, approved: false })}
                style={{
                  background: '#f5f5f5',
                  color: '#999',
                  border: '0.5px solid #e0e0e0',
                  borderRadius: 6,
                  padding: '4px 10px',
                  fontSize: 12,
                  cursor: 'pointer',
                  fontFamily: 'inherit',
                }}
              >
                Отклонить
              </button>
            </div>
          )
        })}
      </div>
    </div>
  )

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
          {isAdmin && (
            <button className="tgh-action" onClick={handleDelete} disabled={isDeleting} title="Удалить игру">
              <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="1.8">
                <path d="M5 7h10l-1 9H6L5 7z" />
                <path d="M3 7h14M8 7V5h4v2" />
              </svg>
            </button>
          )}
        </div>
        {isAdmin && (
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              padding: '10px 16px',
              borderBottom: '0.5px solid #e0e0e0',
              background: '#fafafa',
            }}
          >
            <span style={{ fontSize: 13, color: '#333' }}>
              {game?.is_open ? 'Игра открыта для гостей' : 'Игра закрыта (черновик)'}
            </span>
            <button
              disabled={isTogglingOpen}
              onClick={() => doToggleOpen(!game?.is_open)}
              style={{
                position: 'relative',
                width: 44,
                height: 24,
                borderRadius: 12,
                border: 'none',
                background: game?.is_open ? '#1a1a1a' : '#ccc',
                cursor: 'pointer',
                transition: 'background 0.2s',
                flexShrink: 0,
              }}
            >
              <span
                style={{
                  position: 'absolute',
                  top: 3,
                  left: game?.is_open ? 23 : 3,
                  width: 18,
                  height: 18,
                  borderRadius: '50%',
                  background: '#fff',
                  transition: 'left 0.2s',
                  boxShadow: '0 1px 3px rgba(0,0,0,0.2)',
                }}
              />
            </button>
          </div>
        )}
        <div className="page-body" style={{ padding: 16 }}>
          <div
            style={{
              background: '#f5f5f5',
              borderRadius: 10,
              padding: '12px 16px',
              marginBottom: 16,
              textAlign: 'center',
            }}
          >
            <div style={{ fontSize: 11, color: '#999', marginBottom: 4, textTransform: 'uppercase', letterSpacing: 1 }}>
              Код для входа
            </div>
            <div style={{ fontSize: 26, fontWeight: 700, letterSpacing: 5, fontFamily: 'monospace', color: '#1a1a1a' }}>
              {gameId?.slice(0, 8).toUpperCase()}
            </div>
          </div>

          <TeamsPanel teams={teams} />

          {isAdmin && (
            <button className="tbtn" onClick={() => doStart()} disabled={isStarting}>
              {isStarting ? 'Запускаем…' : 'Начать игру'}
            </button>
          )}
        </div>
      </div>
    )
  }

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
        {isAdmin && (
          <div style={{ display: 'flex', gap: 4 }}>
            <button
              className="tgh-action"
              onClick={() => navigate(`/game/${gameId}/question/add`)}
            >
              <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="2">
                <path d="M10 4v12M4 10h12" />
              </svg>
            </button>
            <button className="tgh-action" onClick={handleDelete} disabled={isDeleting} title="Удалить игру">
              <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="1.8">
                <path d="M5 7h10l-1 9H6L5 7z" />
                <path d="M3 7h14M8 7V5h4v2" />
              </svg>
            </button>
          </div>
        )}
      </div>

      {isAdmin && (
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            padding: '10px 16px',
            borderBottom: '0.5px solid #e0e0e0',
            background: '#fafafa',
          }}
        >
          <span style={{ fontSize: 13, color: '#333' }}>
            {game?.is_open ? 'Игра открыта для гостей' : 'Игра закрыта (черновик)'}
          </span>
          <button
            disabled={isTogglingOpen}
            onClick={() => doToggleOpen(!game?.is_open)}
            style={{
              position: 'relative',
              width: 44,
              height: 24,
              borderRadius: 12,
              border: 'none',
              background: game?.is_open ? '#1a1a1a' : '#ccc',
              cursor: 'pointer',
              transition: 'background 0.2s',
              flexShrink: 0,
            }}
          >
            <span
              style={{
                position: 'absolute',
                top: 3,
                left: game?.is_open ? 23 : 3,
                width: 18,
                height: 18,
                borderRadius: '50%',
                background: '#fff',
                transition: 'left 0.2s',
                boxShadow: '0 1px 3px rgba(0,0,0,0.2)',
              }}
            />
          </button>
        </div>
      )}

      {claimsBanner}

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
                if (!isAdmin) {
                  return <div key={`${pi}-${ci}`} className="qcell empty" />
                }

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
                  onClick={() => {
                    if (isAdmin) {
                      navigate(`/game/${gameId}/question/add`, {
                        state: { categoryId: cat.id, price, questionId: cell.question!.id },
                      })
                    } else {
                      navigate(`/game/${gameId}/question/${cell.question!.id}`)
                    }
                  }}
                >
                  {price}
                </button>
              )
            }),
          )}
        </div>

        <TeamsPanel teams={teams} currentPickerId={game?.current_picker_id} />
      </div>
    </div>
  )
}
