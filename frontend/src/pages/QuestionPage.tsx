import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getQuestion, listTeams, answerQuestion } from '../api'
import { useAuth } from '../App'

export default function QuestionPage() {
  const { gameId, questionId } = useParams<{ gameId: string; questionId: string }>()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { role } = useAuth()
  const isAdmin = role === 'admin'
  const [selected, setSelected] = useState<number | null>(null)
  const [result, setResult] = useState<'correct' | 'wrong' | null>(null)

  const { data: question, isLoading: qLoading } = useQuery({
    queryKey: ['question', questionId],
    queryFn: () => getQuestion(questionId!),
  })

  const { data: teams } = useQuery({
    queryKey: ['teams', gameId],
    queryFn: () => listTeams(gameId!),
    enabled: !!gameId,
  })

  const { mutate: submit, isPending } = useMutation({
    mutationFn: ({ teamId, wrongTeamId }: { teamId: string | null; wrongTeamId?: string | null }) =>
      answerQuestion(gameId!, questionId!, teamId, wrongTeamId ?? null),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['board', gameId] })
      setTimeout(() => navigate(`/game/${gameId}`, { replace: true }), 1200)
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

  const myTeamId = localStorage.getItem(`game:${gameId}:teamId`)
  const myTeam = myTeamId ? (teams ?? []).find(t => t.id === myTeamId) : null

  function handlePick(idx: number) {
    if (selected !== null || isPending) return
    if (isAdmin) return
    if (!myTeam) return

    setSelected(idx)
    const correct = idx === correctIdx

    setResult(correct ? 'correct' : 'wrong')
    submit({
      teamId: correct ? myTeam.id : null,
      wrongTeamId: correct ? null : myTeam.id,
    })
  }

  return (
    <div className="page">
      <div className="tgh">
        <button className="tgh-back" onClick={() => navigate(`/game/${gameId}`)}>
          <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="2">
            <path d="M12 5l-7 5 7 5" />
          </svg>
        </button>
        <span className="tgh-title">{catLabel}</span>
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

        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          {opts.map((opt, i) => {
            const isPicked = selected === i
            const isCorrect = i === correctIdx
            const reveal = selected !== null
            let bg = '#f5f5f5'
            let color = '#1a1a1a'
            let border = '0.5px solid #e0e0e0'

            if (reveal && isCorrect) {
              bg = '#f0faf0'
              color = '#1a6a1a'
              border = '0.5px solid #b2e0b2'
            } else if (reveal && isPicked && !isCorrect) {
              bg = '#fff5f5'
              color = '#8a1a1a'
              border = '0.5px solid #ffb3b3'
            }

            return (
              <button
                key={i}
                onClick={() => handlePick(i)}
                disabled={isAdmin || !myTeam || selected !== null || isPending}
                style={{
                  background: bg,
                  color,
                  border,
                  borderRadius: 10,
                  padding: '14px 12px',
                  fontSize: 15,
                  fontWeight: 500,
                  textAlign: 'left',
                  cursor: isAdmin || selected !== null ? 'default' : 'pointer',
                  fontFamily: 'inherit',
                  opacity: isAdmin && !isCorrect ? 0.6 : 1,
                }}
              >
                {opt}
              </button>
            )
          })}
        </div>

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
