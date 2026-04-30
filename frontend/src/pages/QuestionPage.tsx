import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getQuestion, listTeams, answerQuestion, claimAnswer } from '../api'
import { useAuth } from '../App'

export default function QuestionPage() {
  const { gameId, questionId } = useParams<{ gameId: string; questionId: string }>()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { role } = useAuth()
  const isAdmin = role === 'admin'
  const [revealed, setRevealed] = useState(false)
  const [awarded, setAwarded] = useState(false)
  const [claimed, setClaimed] = useState(false)

  const { data: question, isLoading: qLoading } = useQuery({
    queryKey: ['question', questionId],
    queryFn: () => getQuestion(questionId!),
  })

  const { data: teams } = useQuery({
    queryKey: ['teams', gameId],
    queryFn: () => listTeams(gameId!),
    enabled: !!gameId,
  })

  const { mutate: award } = useMutation({
    mutationFn: (teamId: string | null) => answerQuestion(gameId!, questionId!, teamId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['board', gameId] })
      setAwarded(true)
      setTimeout(() => navigate(`/game/${gameId}`, { replace: true }), 600)
    },
  })

  const { mutate: claim, isPending: isClaiming } = useMutation({
    mutationFn: (teamId: string) => claimAnswer(gameId!, questionId!, teamId),
    onSuccess: () => {
      setClaimed(true)
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

        {!revealed && !awarded && !claimed && (
          <button
            className="tbtn"
            style={{ marginBottom: 8 }}
            onClick={() => setRevealed(true)}
          >
            Показать ответ
          </button>
        )}

        {revealed && (
          <>
            <div className="answer-card">
              <div className="question-label">Ответ</div>
              <div className="answer-text">{question.answer}</div>
              {question.comment && (
                <div
                  style={{
                    marginTop: 8,
                    fontSize: 12,
                    color: '#999',
                    borderTop: '0.5px solid #f0f0f0',
                    paddingTop: 8,
                  }}
                >
                  {question.comment}
                </div>
              )}
            </div>

            {isAdmin && !awarded && (
              <div>
                <div className="text-sm text-mid text-center mb-8">
                  Кто ответил правильно?
                </div>
                <div style={{ display: 'flex', gap: 6 }}>
                  {(teams ?? []).map((team, i) => (
                    <button
                      key={team.id}
                      onClick={() => award(team.id)}
                      style={{
                        flex: 1,
                        background: i === 0 ? '#1a1a1a' : '#f5f5f5',
                        color: i === 0 ? '#fff' : '#333',
                        border: i === 0 ? 'none' : '0.5px solid #e0e0e0',
                        borderRadius: 8,
                        padding: 10,
                        fontSize: 12,
                        fontWeight: i === 0 ? 500 : 400,
                        cursor: 'pointer',
                        fontFamily: 'inherit',
                      }}
                    >
                      {team.name}
                    </button>
                  ))}
                  <button
                    onClick={() => award(null)}
                    style={{
                      flex: 1,
                      background: '#f5f5f5',
                      color: '#999',
                      border: '0.5px solid #e0e0e0',
                      borderRadius: 8,
                      padding: 10,
                      fontSize: 12,
                      cursor: 'pointer',
                      fontFamily: 'inherit',
                    }}
                  >
                    Никто
                  </button>
                </div>
              </div>
            )}

            {!isAdmin && !claimed && (
              <div>
                <div className="text-sm text-mid text-center mb-8">
                  Ваша команда ответила правильно?
                </div>
                <div style={{ display: 'flex', gap: 6 }}>
                  {(teams ?? []).map((team, i) => (
                    <button
                      key={team.id}
                      onClick={() => claim(team.id)}
                      disabled={isClaiming}
                      style={{
                        flex: 1,
                        background: i === 0 ? '#1a1a1a' : '#f5f5f5',
                        color: i === 0 ? '#fff' : '#333',
                        border: i === 0 ? 'none' : '0.5px solid #e0e0e0',
                        borderRadius: 8,
                        padding: 10,
                        fontSize: 12,
                        fontWeight: i === 0 ? 500 : 400,
                        cursor: 'pointer',
                        fontFamily: 'inherit',
                        opacity: isClaiming ? 0.5 : 1,
                      }}
                    >
                      {team.name}
                    </button>
                  ))}
                </div>
              </div>
            )}

            {claimed && (
              <div
                style={{
                  background: '#f0faf0',
                  border: '0.5px solid #b2e0b2',
                  borderRadius: 10,
                  padding: '14px 16px',
                  textAlign: 'center',
                }}
              >
                <div style={{ fontSize: 14, color: '#2a7a2a', fontWeight: 500 }}>
                  Запрос отправлен
                </div>
                <div style={{ fontSize: 12, color: '#666', marginTop: 4 }}>
                  Ждём подтверждения ведущего…
                </div>
              </div>
            )}

            {awarded && (
              <div className="center" style={{ padding: 16 }}>
                <div style={{ fontSize: 14, color: '#999' }}>Записываем…</div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}
