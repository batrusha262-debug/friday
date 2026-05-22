import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getQuestion, listTeams, answerQuestion, claimAnswer, getBoard } from '../api'
import { useAuth } from '../App'

export default function QuestionPage() {
  const { gameId, questionId } = useParams<{ gameId: string; questionId: string }>()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { role } = useAuth()
  const isAdmin = role === 'admin'
  const [revealed, setRevealed] = useState(false)
  const [awarded, setAwarded] = useState(false)
  const [claimedTeamId, setClaimedTeamId] = useState<string | null>(null)
  const [claimConfirmed, setClaimConfirmed] = useState(false)
  const [claimResult, setClaimResult] = useState<'approved' | 'rejected' | null>(null)

  const { data: question, isLoading: qLoading } = useQuery({
    queryKey: ['question', questionId],
    queryFn: () => getQuestion(questionId!),
  })

  const { data: teams } = useQuery({
    queryKey: ['teams', gameId],
    queryFn: () => listTeams(gameId!),
    enabled: !!gameId,
  })

  // Poll board to detect when admin validates our claim
  const { data: board } = useQuery({
    queryKey: ['board-claim', gameId],
    queryFn: () => getBoard(gameId!),
    enabled: claimConfirmed && claimResult === null,
    refetchInterval: claimConfirmed && claimResult === null ? 2000 : false,
  })

  useEffect(() => {
    if (!claimConfirmed || !board || claimResult !== null) return

    const stillPending = board.pending_claims.some(
      c => c.question_id === questionId && c.team_id === claimedTeamId,
    )

    if (!stillPending) {
      const answered = board.states.some(s => s.question_id === questionId)
      setClaimResult(answered ? 'approved' : 'rejected')
    }
  }, [board, claimConfirmed, claimResult, questionId, claimedTeamId])

  useEffect(() => {
    if (claimResult === null) return
    const timer = setTimeout(
      () => navigate(`/game/${gameId}`, { replace: true }),
      2000,
    )

    return () => clearTimeout(timer)
  }, [claimResult, navigate, gameId])

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
    onSuccess: (_, teamId) => {
      setClaimedTeamId(teamId)
      setClaimConfirmed(true)
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

        {!revealed && !awarded && !claimConfirmed && (
          <button
            className="tbtn"
            style={{ marginBottom: 8 }}
            disabled={isClaiming}
            onClick={() => {
              setRevealed(true)
              if (!isAdmin) {
                const myTeamId = localStorage.getItem(`game:${gameId}:teamId`)
                if (myTeamId && (teams ?? []).some(t => t.id === myTeamId)) {
                  claim(myTeamId)
                }
              }
            }}
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
                    fontSize: 13,
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
                <div style={{ display: 'flex', gap: 8 }}>
                  {(teams ?? []).map((team, i) => (
                    <button
                      key={team.id}
                      onClick={() => award(team.id)}
                      style={{
                        flex: 1,
                        background: i === 0 ? '#1a1a1a' : '#f5f5f5',
                        color: i === 0 ? '#fff' : '#333',
                        border: i === 0 ? 'none' : '0.5px solid #e0e0e0',
                        borderRadius: 10,
                        padding: 12,
                        fontSize: 14,
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
                      borderRadius: 10,
                      padding: 12,
                      fontSize: 14,
                      cursor: 'pointer',
                      fontFamily: 'inherit',
                    }}
                  >
                    Никто
                  </button>
                </div>
              </div>
            )}

            {!isAdmin && !claimConfirmed && !isClaiming && (() => {
              const myTeamId = localStorage.getItem(`game:${gameId}:teamId`)
              const myTeam = myTeamId ? (teams ?? []).find(t => t.id === myTeamId) : null

              if (myTeam) return null

              return (
                <div>
                  <div className="text-sm text-mid text-center mb-8">
                    Ваша команда ответила правильно?
                  </div>
                  <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
                    {(teams ?? []).map((team, i) => (
                      <button
                        key={team.id}
                        onClick={() => claim(team.id)}
                        disabled={isClaiming}
                        style={{
                          flex: '1 1 calc(50% - 4px)',
                          background: i === 0 ? '#1a1a1a' : '#f5f5f5',
                          color: i === 0 ? '#fff' : '#333',
                          border: i === 0 ? 'none' : '0.5px solid #e0e0e0',
                          borderRadius: 10,
                          padding: 12,
                          fontSize: 14,
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
              )
            })()}

            {claimConfirmed && claimResult === null && (
              <div
                style={{
                  background: '#fffbe6',
                  border: '0.5px solid #ffe58f',
                  borderRadius: 12,
                  padding: '18px 16px',
                  textAlign: 'center',
                }}
              >
                <div style={{ fontSize: 28, marginBottom: 10 }}>⏳</div>
                <div style={{ fontSize: 16, color: '#7a5c00', fontWeight: 500 }}>
                  Ждём подтверждения ведущего…
                </div>
                <div style={{ fontSize: 13, color: '#b08000', marginTop: 6 }}>
                  Ведущий рассмотрит ваш ответ
                </div>
              </div>
            )}

            {claimResult === 'approved' && (
              <div
                style={{
                  background: '#f0faf0',
                  border: '0.5px solid #b2e0b2',
                  borderRadius: 12,
                  padding: '18px 16px',
                  textAlign: 'center',
                }}
              >
                <div style={{ fontSize: 36, marginBottom: 8 }}>✅</div>
                <div style={{ fontSize: 17, color: '#1a6a1a', fontWeight: 600 }}>
                  Засчитано!
                </div>
                <div style={{ fontSize: 13, color: '#2a8a2a', marginTop: 6 }}>
                  +{question.price} очков команде
                </div>
              </div>
            )}

            {claimResult === 'rejected' && (
              <div
                style={{
                  background: '#fff5f5',
                  border: '0.5px solid #ffb3b3',
                  borderRadius: 12,
                  padding: '18px 16px',
                  textAlign: 'center',
                }}
              >
                <div style={{ fontSize: 36, marginBottom: 8 }}>❌</div>
                <div style={{ fontSize: 17, color: '#8a1a1a', fontWeight: 600 }}>
                  Ответ не принят
                </div>
                <div style={{ fontSize: 13, color: '#aa3333', marginTop: 6 }}>
                  Возвращаемся к игре…
                </div>
              </div>
            )}

            {awarded && (
              <div className="center" style={{ padding: 16 }}>
                <div style={{ fontSize: 15, color: '#999' }}>Записываем…</div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}
