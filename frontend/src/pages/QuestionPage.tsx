import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  getQuestion,
  listTeams,
  answerQuestion,
  getBoard,
  getGame,
  openQuestion,
  revealNextOption,
  startQuestionTimer,
} from '../api'
import { useAuth } from '../App'
import { useGameEvents } from '../hooks/useGameEvents'
import type { Game, GameBoard, GameQuestionState } from '../api/types'

const TIMER_SECONDS = 20

export default function QuestionPage() {
  const { gameId, questionId } = useParams<{ gameId: string; questionId: string }>()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { role } = useAuth()
  const isAdmin = role === 'admin'
  const [pickedIdx, setPickedIdx] = useState<number | null>(null)
  const [now, setNow] = useState(() => Date.now())

  const [liveGame, setLiveGame] = useState<Game | null>(null)
  const [liveBoard, setLiveBoard] = useState<GameBoard | null>(null)

  useGameEvents(gameId!, state => {
    setLiveGame(state.game)
    setLiveBoard(state.board)
  })

  const { data: question, isLoading: qLoading } = useQuery({
    queryKey: ['question', questionId],
    queryFn: () => getQuestion(questionId!),
  })

  const { data: teams } = useQuery({
    queryKey: ['teams', gameId],
    queryFn: () => listTeams(gameId!),
    enabled: !!gameId,
  })

  const gameQuery = useQuery({
    queryKey: ['game', gameId],
    queryFn: () => getGame(gameId!),
    refetchInterval: 2000,
    enabled: !!gameId,
  })

  const boardQuery = useQuery({
    queryKey: ['board', gameId],
    queryFn: () => getBoard(gameId!),
    refetchInterval: 1000,
    enabled: !!gameId,
  })

  const game = liveGame ?? gameQuery.data
  const board = liveBoard ?? boardQuery.data
  const state: GameQuestionState | undefined = board?.states.find(
    s => s.question_id === questionId,
  )

  const ensureOpen = useMutation({
    mutationFn: () => openQuestion(gameId!, questionId!),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['board', gameId] }),
  })

  useEffect(() => {
    if (!isAdmin || !state || ensureOpen.isPending) return
    if (state.revealed_count > 0 || state.timer_started_at) return
    ensureOpen.mutate()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAdmin, state?.id])

  useEffect(() => {
    if (!state?.timer_started_at) return
    const id = setInterval(() => setNow(Date.now()), 200)

    return () => clearInterval(id)
  }, [state?.timer_started_at])

  const { mutate: doReveal, isPending: isRevealing } = useMutation({
    mutationFn: () => revealNextOption(gameId!, questionId!),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['board', gameId] }),
  })

  const { mutate: doStartTimer, isPending: isStartingTimer } = useMutation({
    mutationFn: () => startQuestionTimer(gameId!, questionId!),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['board', gameId] }),
  })

  const { mutate: submit, isPending } = useMutation({
    mutationFn: ({
      teamId,
      wrongTeamId,
      optionIdx,
    }: {
      teamId: string | null
      wrongTeamId?: string | null
      optionIdx?: number | null
    }) => answerQuestion(gameId!, questionId!, teamId, wrongTeamId ?? null, optionIdx ?? null),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['board', gameId] })
      setTimeout(() => navigate(`/game/${gameId}`, { replace: true }), 1500)
    },
  })

  if (qLoading) {
    return (
      <div className="page">
        <div className="tgh"><span className="tgh-title">Загрузка…</span></div>
        <div className="center"><div className="spinner" /></div>
      </div>
    )
  }

  if (!question) {
    return (
      <div className="page">
        <div className="tgh">
          <button className="tgh-back" onClick={() => navigate(-1)}>
            <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="2">
              <path d="M12 5l-7 5 7 5" />
            </svg>
          </button>
          <span className="tgh-title">Вопрос не найден</span>
        </div>
      </div>
    )
  }

  const catLabel = `${question.price} очков`
  const opts = question.options ?? []
  const correctIdx = question.correct_option ?? 0
  const revealedCount = state?.revealed_count ?? 0
  const wrongOptions = new Set(state?.wrong_options ?? [])
  const timerStartedAt = state?.timer_started_at ? new Date(state.timer_started_at).getTime() : null
  const elapsedSec = timerStartedAt ? (now - timerStartedAt) / 1000 : 0
  const remainingSec = timerStartedAt ? Math.max(0, TIMER_SECONDS - elapsedSec) : null
  const timerExpired = remainingSec !== null && remainingSec <= 0

  const myTeamId = localStorage.getItem(`game:${gameId}:teamId`)
  const myTeam = myTeamId ? (teams ?? []).find(t => t.id === myTeamId) : null
  const isPicker = myTeamId && game?.current_picker_id === myTeamId

  function handlePick(idx: number) {
    if (pickedIdx !== null || isPending) return
    if (isAdmin) return
    if (!myTeam) return
    if (!isPicker) return
    if (wrongOptions.has(idx)) return
    if (idx >= revealedCount) return
    if (timerExpired) return

    setPickedIdx(idx)
    const correct = idx === correctIdx

    submit({
      teamId: correct ? myTeam.id : null,
      wrongTeamId: correct ? null : myTeam.id,
      optionIdx: idx,
    })
  }

  const result: 'correct' | 'wrong' | null =
    pickedIdx === null ? null : pickedIdx === correctIdx ? 'correct' : 'wrong'

  const answered = pickedIdx !== null
  const canLeave = isAdmin || answered

  return (
    <div className="page">
      <div className="tgh">
        {canLeave ? (
          <button className="tgh-back" onClick={() => navigate(`/game/${gameId}`)}>
            <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="2">
              <path d="M12 5l-7 5 7 5" />
            </svg>
          </button>
        ) : (
          <span style={{ width: 32 }} />
        )}
        <span className="tgh-title">{catLabel}</span>
        {timerStartedAt && (
          <span
            style={{
              marginLeft: 'auto',
              fontSize: 14,
              fontWeight: 700,
              color: remainingSec! <= 5 ? '#ff4444' : '#fff',
              fontFamily: 'monospace',
            }}
          >
            {Math.ceil(remainingSec ?? 0)}с
          </span>
        )}
      </div>

      <div className="page-body" style={{ padding: 12 }}>
        <div className="question-card">
          <div className="question-label">Вопрос</div>
          <div className="question-text">{question.question}</div>
        </div>

        {isAdmin && (
          <div
            style={{
              background: '#fafafa',
              border: '0.5px solid #e0e0e0',
              borderRadius: 10,
              padding: 10,
              marginBottom: 8,
              fontSize: 13,
              color: '#666',
            }}
          >
            Правильный ответ: <b style={{ color: '#1a6a1a' }}>{opts[correctIdx]}</b>
          </div>
        )}

        {!isAdmin && !myTeam && (
          <div
            style={{
              background: '#fffbe6',
              border: '0.5px solid #ffe58f',
              borderRadius: 12,
              padding: '14px 16px',
              textAlign: 'center',
              fontSize: 14,
              color: '#7a5c00',
            }}
          >
            Ваша команда не определена — вернитесь в лобби.
          </div>
        )}

        {!isAdmin && myTeam && !isPicker && (
          <div
            style={{
              background: '#fffbe6',
              border: '0.5px solid #ffe58f',
              borderRadius: 12,
              padding: '14px 16px',
              textAlign: 'center',
              fontSize: 14,
              color: '#7a5c00',
              marginBottom: 8,
            }}
          >
            Сейчас отвечает другая команда. Ждите своей очереди.
          </div>
        )}

        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          {opts.map((opt, i) => {
            const visible = i < revealedCount
            const isWrong = wrongOptions.has(i)
            const isPicked = pickedIdx === i
            const isCorrect = i === correctIdx
            const reveal = pickedIdx !== null

            if (!visible) {
              return (
                <div
                  key={i}
                  style={{
                    background: '#f5f5f5',
                    border: '0.5px dashed #ccc',
                    borderRadius: 10,
                    padding: '14px 12px',
                    fontSize: 14,
                    color: '#bbb',
                    textAlign: 'center',
                  }}
                >
                  Вариант скрыт
                </div>
              )
            }

            let bg = '#f5f5f5'
            let color = '#1a1a1a'
            let border = '0.5px solid #e0e0e0'
            let textDecoration: string | undefined

            if (isWrong) {
              bg = '#fff5f5'
              color = '#999'
              border = '0.5px solid #ffcccc'
              textDecoration = 'line-through'
            }

            if (reveal && isCorrect) {
              bg = '#f0faf0'
              color = '#1a6a1a'
              border = '0.5px solid #b2e0b2'
              textDecoration = undefined
            } else if (reveal && isPicked && !isCorrect) {
              bg = '#fff5f5'
              color = '#8a1a1a'
              border = '0.5px solid #ffb3b3'
            }

            const disabled =
              isAdmin ||
              !myTeam ||
              !isPicker ||
              pickedIdx !== null ||
              isPending ||
              isWrong ||
              timerExpired

            return (
              <button
                key={i}
                onClick={() => handlePick(i)}
                disabled={disabled}
                style={{
                  background: bg,
                  color,
                  border,
                  borderRadius: 10,
                  padding: '14px 12px',
                  fontSize: 15,
                  fontWeight: 500,
                  textAlign: 'left',
                  cursor: disabled ? 'default' : 'pointer',
                  fontFamily: 'inherit',
                  opacity: isAdmin && !isCorrect ? 0.6 : 1,
                  textDecoration,
                }}
              >
                {opt}
              </button>
            )
          })}
        </div>

        {isAdmin && (
          <div style={{ display: 'flex', gap: 8, marginTop: 12 }}>
            <button
              className="tbtn-ghost"
              style={{ flex: 1 }}
              disabled={isRevealing || revealedCount >= opts.length}
              onClick={() => doReveal()}
            >
              {revealedCount === 0
                ? 'Показать 1-й вариант'
                : revealedCount >= opts.length
                  ? 'Все варианты показаны'
                  : `Показать ${revealedCount + 1}-й вариант`}
            </button>
            <button
              className="tbtn"
              style={{ flex: 1 }}
              disabled={isStartingTimer || !!timerStartedAt || revealedCount === 0}
              onClick={() => doStartTimer()}
            >
              {timerStartedAt ? 'Таймер запущен' : `Старт ${TIMER_SECONDS}с`}
            </button>
          </div>
        )}

        {result === 'correct' && (
          <div
            style={{
              marginTop: 12,
              background: '#f0faf0',
              border: '0.5px solid #b2e0b2',
              borderRadius: 12,
              padding: '14px 16px',
              textAlign: 'center',
            }}
          >
            <div style={{ fontSize: 32, marginBottom: 6 }}>✅</div>
            <div style={{ fontSize: 16, color: '#1a6a1a', fontWeight: 600 }}>
              +{question.price} очков
            </div>
          </div>
        )}

        {result === 'wrong' && (
          <div
            style={{
              marginTop: 12,
              background: '#fff5f5',
              border: '0.5px solid #ffb3b3',
              borderRadius: 12,
              padding: '14px 16px',
              textAlign: 'center',
            }}
          >
            <div style={{ fontSize: 32, marginBottom: 6 }}>❌</div>
            <div style={{ fontSize: 15, color: '#8a1a1a', fontWeight: 600 }}>
              Неверно
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
