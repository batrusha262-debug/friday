import { useState, useEffect, useRef } from 'react'
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
  joinGame,
  claimMiniGame,
  resetGame,
} from '../api'
import { useGameEvents } from '../hooks/useGameEvents'
import { useAuth } from '../App'
import type { AnswerClaim, Game, GameBoard, GameTeam, MiniGame, Question } from '../api/types'

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

  const gameQuery = useQuery({
    queryKey: ['game', gameId],
    queryFn: () => getGame(gameId),
    refetchInterval: 2000,
  })
  const boardQuery = useQuery({
    queryKey: ['board', gameId],
    queryFn: () => getBoard(gameId),
    refetchInterval: 2000,
  })

  const effectiveGame = liveGame ?? gameQuery.data ?? null
  const effectiveTeams: GameTeam[] = (liveBoard ?? boardQuery.data)?.teams ?? []
  const effectiveStates = (liveBoard ?? boardQuery.data)?.states ?? []
  const pendingClaims: AnswerClaim[] = (liveBoard ?? boardQuery.data)?.pending_claims ?? []
  const miniGame = (liveBoard ?? boardQuery.data)?.mini_game ?? null

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
  const answeredIds = new Set(effectiveStates.filter(s => !!s.answered_by).map(s => s.question_id))
  const stateByQuestion = Object.fromEntries(effectiveStates.map(s => [s.question_id, s]))

  const allQuestions = questionsQueries.data ?? {}

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

  const totalExpected = prices.length * catList.length
  const filledCount = grid.reduce(
    (acc, row) => acc + row.reduce((a, c) => a + (c.question ? 1 : 0), 0),
    0,
  )
  const allFilled = totalExpected > 0 && filledCount === totalExpected

  return {
    loading: gameQuery.isLoading || boardQuery.isLoading || rounds.isLoading || categories.isLoading,
    game: effectiveGame,
    teams: effectiveTeams,
    categories: catList,
    prices,
    grid,
    pendingClaims,
    questionById,
    miniGame,
    allFilled,
    filledCount,
    totalExpected,
  }
}

function OnlineTab({ teams, currentPickerId }: { teams: GameTeam[]; currentPickerId?: string }) {
  const sorted = [...teams].sort((a, b) => b.score - a.score)

  if (teams.length === 0) {
    return (
      <div className="empty">
        <div className="empty-icon">👥</div>
        Команд пока нет
      </div>
    )
  }

  return (
    <div style={{ padding: '10px 12px', display: 'flex', flexDirection: 'column', gap: 6 }}>
      {sorted.map((team, i) => {
        const isPicker = team.id === currentPickerId
        const isFirst = i === 0

        return (
          <div
            key={team.id}
            style={{
              background: isFirst ? '#1a1a1a' : '#fff',
              borderRadius: 12,
              padding: '14px 16px',
              display: 'flex',
              alignItems: 'center',
              gap: 12,
              outline: isPicker ? '2px solid #f0a500' : undefined,
              outlineOffset: isPicker ? 2 : undefined,
            }}
          >
            <div
              style={{
                fontSize: 20,
                fontWeight: 700,
                color: isFirst ? 'rgba(255,255,255,0.4)' : '#ccc',
                width: 28,
                textAlign: 'center',
              }}
            >
              {i + 1}
            </div>
            <div style={{ flex: 1 }}>
              <div style={{ fontSize: 16, fontWeight: 500, color: isFirst ? '#fff' : '#1a1a1a' }}>
                {team.name}
              </div>
              {isPicker && (
                <div style={{ fontSize: 12, color: '#f0a500', marginTop: 2 }}>▶ сейчас ходит</div>
              )}
            </div>
            <div style={{ fontSize: 24, fontWeight: 700, color: isFirst ? '#fff' : '#1a1a1a' }}>
              {team.score}
            </div>
          </div>
        )
      })}
    </div>
  )
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
              <div style={{ fontSize: 11, color: isLeading ? '#aaa' : '#999' }}>
                {team.name}
              </div>
              <div style={{ fontSize: 18, fontWeight: 500, color: isLeading ? '#fff' : '#333' }}>
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

function MiniGameOverlay({
  miniGame,
  gameId,
  myTeamId,
  teams,
}: {
  miniGame: MiniGame
  gameId: string
  myTeamId: string | null
  teams: GameTeam[]
}) {
  const qc = useQueryClient()
  const navigate = useNavigate()
  const [now, setNow] = useState(() => Date.now())
  const [outcome, setOutcome] = useState<'win' | 'lose' | null>(null)

  useEffect(() => {
    const id = setInterval(() => setNow(Date.now()), 100)

    return () => clearInterval(id)
  }, [])

  useEffect(() => {
    if (outcome !== 'win') return
    const t = setTimeout(() => {
      navigate(`/game/${gameId}/question/${miniGame.question_id}`)
    }, 1200)

    return () => clearTimeout(t)
  }, [outcome, gameId, miniGame.question_id, navigate])

  const { mutate: claim, isPending } = useMutation({
    mutationFn: () => claimMiniGame(gameId, miniGame.id, myTeamId!),
    onSuccess: (mg) => {
      setOutcome(mg.winner_team_id === myTeamId ? 'win' : 'lose')
      qc.invalidateQueries({ queryKey: ['board', gameId] })
    },
    onError: () => {
      setOutcome('lose')
    },
  })

  const appearsAt = new Date(miniGame.appears_at).getTime()
  const remainingMs = Math.max(0, appearsAt - now)
  const remainingSec = Math.ceil(remainingMs / 1000)
  const buttonVisible = remainingMs <= 0

  const isExcluded = myTeamId !== null && miniGame.excluded_team_id === myTeamId
  const isWinner = miniGame.winner_team_id !== undefined && miniGame.winner_team_id !== null

  if (isWinner && !outcome) {
    const winTeam = teams.find(t => t.id === miniGame.winner_team_id)

    return (
      <div style={overlayBg}>
        <div style={overlayCard}>
          <div style={{ fontSize: 36, marginBottom: 8 }}>🏁</div>
          <div style={{ fontSize: 17, fontWeight: 600, color: '#1a1a1a' }}>
            Ход переходит к
          </div>
          <div style={{ fontSize: 20, fontWeight: 700, marginTop: 4 }}>
            {winTeam?.name ?? '—'}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div style={overlayBg}>
      {!buttonVisible && (
        <div style={overlayCard}>
          <div style={{ fontSize: 36, marginBottom: 8 }}>⚡</div>
          <div style={{ fontSize: 16, fontWeight: 600, color: '#1a1a1a' }}>
            {isExcluded ? 'Вы выбыли из гонки' : 'Приготовьтесь!'}
          </div>
          <div style={{ fontSize: 14, color: '#666', marginTop: 6 }}>
            Кнопка появится через {remainingSec}с — кто первый, тот ходит
          </div>
        </div>
      )}

      {buttonVisible && !isExcluded && !outcome && myTeamId && (
        <button
          onClick={() => !isPending && claim()}
          disabled={isPending}
          style={{
            position: 'fixed',
            left: `${miniGame.pos_x}%`,
            top: `${miniGame.pos_y}%`,
            transform: 'translate(-50%, -50%)',
            background: '#f0a500',
            color: '#fff',
            border: 'none',
            borderRadius: '50%',
            width: 72,
            height: 72,
            fontSize: 28,
            fontWeight: 700,
            cursor: 'pointer',
            boxShadow: '0 6px 20px rgba(240,165,0,0.6)',
            zIndex: 1001,
            fontFamily: 'inherit',
          }}
        >
          ⚡
        </button>
      )}

      {buttonVisible && isExcluded && (
        <div style={overlayCard}>
          <div style={{ fontSize: 36, marginBottom: 8 }}>🚫</div>
          <div style={{ fontSize: 16, fontWeight: 600, color: '#1a1a1a' }}>
            Вы выбыли из гонки
          </div>
          <div style={{ fontSize: 14, color: '#666', marginTop: 6 }}>
            Ждём пока другие команды нажмут кнопку…
          </div>
        </div>
      )}

      {buttonVisible && !isExcluded && !myTeamId && (
        <div style={overlayCard}>
          <div style={{ fontSize: 16, color: '#666' }}>Идёт мини-игра между командами</div>
        </div>
      )}

      {outcome === 'win' && (
        <div style={overlayCard}>
          <div style={{ fontSize: 36, marginBottom: 8 }}>🏆</div>
          <div style={{ fontSize: 17, fontWeight: 600, color: '#1a6a1a' }}>
            Ход ваш!
          </div>
        </div>
      )}

      {outcome === 'lose' && (
        <div style={overlayCard}>
          <div style={{ fontSize: 32, marginBottom: 8 }}>⏱</div>
          <div style={{ fontSize: 16, color: '#666' }}>Опоздали</div>
        </div>
      )}
    </div>
  )
}

const overlayBg: React.CSSProperties = {
  position: 'fixed',
  inset: 0,
  background: 'rgba(0,0,0,0.55)',
  zIndex: 1000,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  padding: 20,
}

const overlayCard: React.CSSProperties = {
  background: '#fff',
  borderRadius: 16,
  padding: '20px 24px',
  textAlign: 'center',
  boxShadow: '0 10px 30px rgba(0,0,0,0.25)',
  maxWidth: 320,
}

export default function GameBoardPage() {
  const { gameId } = useParams<{ gameId: string }>()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { role } = useAuth()
  const isAdmin = role === 'admin'
  const [tab, setTab] = useState<'board' | 'online'>('board')
  const [joinError, setJoinError] = useState('')
  const [resetOpen, setResetOpen] = useState(false)
  const [resetText, setResetText] = useState('')
  const [resetError, setResetError] = useState('')
  const { loading, game, teams, categories, prices, grid, pendingClaims, questionById, miniGame, allFilled, filledCount, totalExpected } = useBoardData(gameId!)

  const autoJoinedRef = useRef(false)

  useEffect(() => {
    if (!isAdmin || !gameId) return
    const qid = game?.current_question_id
    const key = `game:${gameId}:lastFollowedQuestion`
    if (!qid) {
      sessionStorage.removeItem(key)
      return
    }
    if (sessionStorage.getItem(key) === qid) return
    sessionStorage.setItem(key, qid)
    navigate(`/game/${gameId}/question/${qid}`)
  }, [isAdmin, gameId, game?.current_question_id, navigate])

  // Guest auto-joins as a team on entering an open waiting lobby
  useEffect(() => {
    if (isAdmin || !gameId || !game) return
    if (game.status !== 'waiting' || !game.is_open) return
    if (autoJoinedRef.current) return
    if (localStorage.getItem(`game:${gameId}:teamId`)) return

    const myName = localStorage.getItem('userName')?.trim()

    if (!myName) return

    autoJoinedRef.current = true
    joinGame(gameId, myName)
      .then(team => {
        localStorage.setItem(`game:${gameId}:teamId`, team.id)
        qc.invalidateQueries({ queryKey: ['board', gameId] })
      })
      .catch(err => {
        autoJoinedRef.current = false
        setJoinError(err instanceof Error ? err.message : 'Не удалось присоединиться')
      })
  }, [isAdmin, gameId, game, qc])

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

  const RESET_PHRASE = '67 крутой мем'
  const { mutate: doReset, isPending: isResetting } = useMutation({
    mutationFn: () => resetGame(gameId!, resetText),
    onSuccess: () => {
      setResetOpen(false)
      setResetText('')
      setResetError('')
      qc.invalidateQueries({ queryKey: ['game', gameId] })
      qc.invalidateQueries({ queryKey: ['board', gameId] })
    },
    onError: (err: unknown) => {
      setResetError(err instanceof Error ? err.message : 'Не удалось сбросить игру')
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
    localStorage.removeItem(`game:${gameId}:teamId`)
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

  // Pending claim modal (admin only) — shows top-most pending claim
  const topClaim = isAdmin ? pendingClaims[0] : null
  const claimsModal = topClaim && (() => {
    const team = teams.find(t => t.id === topClaim.team_id)
    const question = questionById[topClaim.question_id]
    const extra = pendingClaims.length - 1

    return (
      <div
        style={{
          position: 'fixed',
          inset: 0,
          background: 'rgba(0,0,0,0.45)',
          zIndex: 1000,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          padding: 20,
        }}
      >
        <div
          style={{
            background: '#fff',
            borderRadius: 16,
            padding: 20,
            width: '100%',
            maxWidth: 360,
            boxShadow: '0 10px 30px rgba(0,0,0,0.25)',
          }}
        >
          <div
            style={{
              fontSize: 12,
              color: '#996600',
              textTransform: 'uppercase',
              letterSpacing: 0.5,
              marginBottom: 10,
              textAlign: 'center',
            }}
          >
            Подтверждение ответа
          </div>
          <div style={{ fontSize: 17, textAlign: 'center', marginBottom: 4 }}>
            <strong>{team?.name ?? '…'}</strong>
          </div>
          {question && (
            <div style={{ fontSize: 14, color: '#666', textAlign: 'center', marginBottom: 14 }}>
              {question.question}
            </div>
          )}
          {question && (
            <div style={{ fontSize: 13, color: '#999', textAlign: 'center', marginBottom: 16 }}>
              Ответ: <span style={{ color: '#333' }}>{question.answer}</span> · {question.price} очков
            </div>
          )}
          <div style={{ display: 'flex', gap: 8 }}>
            <button
              onClick={() => doValidate({ claimId: topClaim.id, approved: false })}
              style={{
                flex: 1,
                background: '#f5f5f5',
                color: '#333',
                border: '0.5px solid #e0e0e0',
                borderRadius: 10,
                padding: '12px 14px',
                fontSize: 14,
                cursor: 'pointer',
                fontFamily: 'inherit',
              }}
            >
              Не засчитывать
            </button>
            <button
              onClick={() => doValidate({ claimId: topClaim.id, approved: true })}
              style={{
                flex: 1,
                background: '#1a1a1a',
                color: '#fff',
                border: 'none',
                borderRadius: 10,
                padding: '12px 14px',
                fontSize: 14,
                fontWeight: 500,
                cursor: 'pointer',
                fontFamily: 'inherit',
              }}
            >
              Засчитать
            </button>
          </div>
          {extra > 0 && (
            <div style={{ fontSize: 12, color: '#999', textAlign: 'center', marginTop: 12 }}>
              Ещё {extra} в очереди
            </div>
          )}
        </div>
      </div>
    )
  })()

  const resetModal = isAdmin && resetOpen && (
    <div
      style={{
        position: 'fixed',
        inset: 0,
        background: 'rgba(0,0,0,0.45)',
        zIndex: 1100,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        padding: 20,
      }}
    >
      <div
        style={{
          background: '#fff',
          borderRadius: 16,
          padding: 20,
          width: '100%',
          maxWidth: 360,
          boxShadow: '0 10px 30px rgba(0,0,0,0.25)',
        }}
      >
        <div
          style={{
            fontSize: 12,
            color: '#c00',
            textTransform: 'uppercase',
            letterSpacing: 0.5,
            marginBottom: 10,
            textAlign: 'center',
          }}
        >
          Сброс игры
        </div>
        <div style={{ fontSize: 15, color: '#1a1a1a', textAlign: 'center', marginBottom: 6 }}>
          Все ответы, очки и текущий ход будут сброшены.
        </div>
        <div style={{ fontSize: 13, color: '#666', textAlign: 'center', marginBottom: 14 }}>
          Чтобы подтвердить, введите фразу:
          <div style={{ fontFamily: 'monospace', color: '#1a1a1a', marginTop: 4 }}>
            {RESET_PHRASE}
          </div>
        </div>
        <input
          type="text"
          value={resetText}
          onChange={e => {
            setResetText(e.target.value)
            setResetError('')
          }}
          placeholder={RESET_PHRASE}
          style={{
            width: '100%',
            padding: '12px 14px',
            borderRadius: 10,
            border: '0.5px solid #e0e0e0',
            fontSize: 15,
            marginBottom: 10,
            fontFamily: 'inherit',
            boxSizing: 'border-box',
          }}
          autoFocus
        />
        {resetError && (
          <div style={{ fontSize: 12, color: '#c00', textAlign: 'center', marginBottom: 10 }}>
            {resetError}
          </div>
        )}
        <div style={{ display: 'flex', gap: 8 }}>
          <button
            onClick={() => {
              setResetOpen(false)
              setResetText('')
              setResetError('')
            }}
            style={{
              flex: 1,
              background: '#f5f5f5',
              color: '#333',
              border: '0.5px solid #e0e0e0',
              borderRadius: 10,
              padding: '12px 14px',
              fontSize: 14,
              cursor: 'pointer',
              fontFamily: 'inherit',
            }}
          >
            Отмена
          </button>
          <button
            onClick={() => doReset()}
            disabled={isResetting || resetText !== RESET_PHRASE}
            style={{
              flex: 1,
              background: resetText === RESET_PHRASE ? '#c00' : '#e0e0e0',
              color: '#fff',
              border: 'none',
              borderRadius: 10,
              padding: '12px 14px',
              fontSize: 14,
              fontWeight: 500,
              cursor: resetText === RESET_PHRASE && !isResetting ? 'pointer' : 'not-allowed',
              fontFamily: 'inherit',
            }}
          >
            {isResetting ? 'Сбрасываем…' : 'Сбросить'}
          </button>
        </div>
      </div>
    </div>
  )

  const canOpen = allFilled
  const openToggle = isAdmin && (
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
      <div style={{ display: 'flex', flexDirection: 'column' }}>
        <span style={{ fontSize: 14, color: '#333' }}>
          {game?.is_open ? 'Игра открыта для гостей' : 'Игра закрыта (черновик)'}
        </span>
        {!canOpen && (
          <span style={{ fontSize: 11, color: '#c00', marginTop: 2 }}>
            Заполните все вопросы ({filledCount}/{totalExpected})
          </span>
        )}
      </div>
      <button
        disabled={isTogglingOpen || (!game?.is_open && !canOpen)}
        onClick={() => doToggleOpen(!game?.is_open)}
        style={{
          position: 'relative',
          width: 44,
          height: 24,
          borderRadius: 12,
          border: 'none',
          background: game?.is_open ? '#1a1a1a' : '#ccc',
          cursor: !game?.is_open && !canOpen ? 'not-allowed' : 'pointer',
          opacity: !game?.is_open && !canOpen ? 0.5 : 1,
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
  )

  // Lobby: waiting for players.
  // Admin sees the question-filling board first; the lobby appears once all
  // questions are filled.
  if (game?.status === 'waiting' && (!isAdmin || allFilled)) {
    return (
      <div className="page">
        <div className="tgh">
          <button className="tgh-back" onClick={() => navigate('/')}>
            <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="2">
              <path d="M12 5l-7 5 7 5" />
            </svg>
          </button>
          <span className="tgh-title">Комната ожидания</span>
          {isAdmin && (
            <button className="tgh-action" onClick={handleDelete} disabled={isDeleting} title="Удалить игру">
              <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="1.8">
                <path d="M5 7h10l-1 9H6L5 7z" />
                <path d="M3 7h14M8 7V5h4v2" />
              </svg>
            </button>
          )}
        </div>

        {openToggle}

        <div className="page-body" style={{ padding: 14 }}>
          {/* Join code */}
          <div
            style={{
              background: '#fff',
              borderRadius: 14,
              padding: '18px 16px',
              marginBottom: 14,
              textAlign: 'center',
              border: '0.5px solid #e0e0e0',
            }}
          >
            <div style={{ fontSize: 12, color: '#999', marginBottom: 6, textTransform: 'uppercase', letterSpacing: 1.5 }}>
              Код для входа
            </div>
            <div style={{ fontSize: 32, fontWeight: 700, letterSpacing: 6, fontFamily: 'monospace', color: '#1a1a1a' }}>
              {gameId?.slice(0, 8).toUpperCase()}
            </div>
            <div style={{ fontSize: 12, color: '#bbb', marginTop: 8 }}>
              Введите код на экране входа
            </div>
          </div>

          {/* Teams list */}
          <div
            style={{
              background: '#fff',
              borderRadius: 14,
              overflow: 'hidden',
              marginBottom: 14,
              border: '0.5px solid #e0e0e0',
            }}
          >
            <div style={{ padding: '12px 16px', borderBottom: '0.5px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span style={{ fontSize: 14, fontWeight: 500, color: '#1a1a1a' }}>Команды</span>
              <span style={{ fontSize: 14, color: '#999' }}>{teams.length}</span>
            </div>
            {teams.length === 0 ? (
              <div style={{ padding: '20px 16px', textAlign: 'center', color: '#bbb', fontSize: 14 }}>
                Никто ещё не подключился
              </div>
            ) : (
              teams.map((team, i) => (
                <div
                  key={team.id}
                  style={{
                    padding: '12px 16px',
                    display: 'flex',
                    alignItems: 'center',
                    gap: 10,
                    borderBottom: i < teams.length - 1 ? '0.5px solid #f5f5f5' : undefined,
                  }}
                >
                  <div
                    style={{
                      width: 34,
                      height: 34,
                      borderRadius: '50%',
                      background: '#1a1a1a',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      color: '#fff',
                      fontWeight: 600,
                      fontSize: 14,
                      flexShrink: 0,
                    }}
                  >
                    {team.name.charAt(0).toUpperCase()}
                  </div>
                  <span style={{ fontSize: 15, color: '#1a1a1a' }}>{team.name}</span>
                </div>
              ))
            )}
          </div>

          {isAdmin && (
            <button
              className="tbtn"
              onClick={() => doStart()}
              disabled={isStarting || teams.length === 0 || !allFilled}
            >
              {isStarting
                ? 'Запускаем…'
                : !allFilled
                  ? `Заполните вопросы (${filledCount}/${totalExpected})`
                  : `Начать игру${teams.length > 0 ? ` (${teams.length} команд)` : ''}`}
            </button>
          )}

          {!isAdmin && (
            <div
              style={{
                background: '#f5f5f5',
                borderRadius: 12,
                padding: '14px 16px',
                textAlign: 'center',
              }}
            >
              <div style={{ fontSize: 15, color: '#666' }}>Ждём начала игры…</div>
              <div style={{ fontSize: 12, color: '#bbb', marginTop: 4 }}>
                {localStorage.getItem(`game:${gameId}:teamId`)
                  ? 'Вы в команде, ведущий скоро запустит'
                  : 'Подключаемся…'}
              </div>
              {joinError && (
                <div style={{ fontSize: 12, color: '#c00', marginTop: 8 }}>{joinError}</div>
              )}
            </div>
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
            <button
              className="tgh-action"
              onClick={() => setResetOpen(true)}
              title="Сбросить игру"
            >
              <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="1.8">
                <path d="M4 10a6 6 0 1 0 1.76-4.24" />
                <path d="M4 4v3.5h3.5" strokeLinecap="round" strokeLinejoin="round" />
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

      {openToggle}
      {claimsModal}
      {resetModal}
      {miniGame && gameId && (
        <MiniGameOverlay
          miniGame={miniGame}
          gameId={gameId}
          myTeamId={localStorage.getItem(`game:${gameId}:teamId`)}
          teams={teams}
        />
      )}

      {/* Tab bar */}
      <div className="tab-bar">
        <button
          className={`tab-btn ${tab === 'board' ? 'active' : ''}`}
          onClick={() => setTab('board')}
        >
          Игровое поле
        </button>
        <button
          className={`tab-btn ${tab === 'online' ? 'active' : ''}`}
          onClick={() => setTab('online')}
        >
          Онлайн {teams.length > 0 && `(${teams.length})`}
        </button>
      </div>

      <div className="page-body">
        {tab === 'board' && (
          <>
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

                  const myTeamId = localStorage.getItem(`game:${gameId}:teamId`)
                  const isPicker = !isAdmin && myTeamId === game?.current_picker_id
                  const canClick = isAdmin || isPicker

                  return (
                    <button
                      key={`${pi}-${ci}`}
                      className="qcell"
                      disabled={!canClick}
                      style={!canClick ? { opacity: 0.6, cursor: 'not-allowed' } : undefined}
                      onClick={() => {
                        if (!canClick) return
                        if (isAdmin && game?.status === 'waiting') {
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
          </>
        )}

        {tab === 'online' && (
          <OnlineTab teams={teams} currentPickerId={game?.current_picker_id} />
        )}
      </div>
    </div>
  )
}
