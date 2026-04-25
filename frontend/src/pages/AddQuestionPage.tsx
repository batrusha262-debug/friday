import { useState } from 'react'
import { useNavigate, useParams, useLocation } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getGame, listRounds, listCategories, listQuestions, createQuestion } from '../api'

const DEFAULT_PRICES = [100, 200, 300, 400, 500]

export default function AddQuestionPage() {
  const { gameId } = useParams<{ gameId: string }>()
  const navigate = useNavigate()
  const location = useLocation()
  const qc = useQueryClient()
  const state = location.state as { categoryId?: string; price?: number } | null

  const game = useQuery({ queryKey: ['game', gameId], queryFn: () => getGame(gameId!) })
  const packId = game.data?.pack_id

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
  const [catId, setCatId] = useState(state?.categoryId ?? '')
  const [price, setPrice] = useState(state?.price ?? 100)
  const [questionText, setQuestionText] = useState('')
  const [answer, setAnswer] = useState('')
  const [comment, setComment] = useState('')
  const [error, setError] = useState('')

  const effectiveCatId = catId || catList[0]?.id || ''

  const existingQuestions = useQuery({
    queryKey: ['questions', effectiveCatId],
    queryFn: () => listQuestions(effectiveCatId),
    enabled: !!effectiveCatId,
  })

  const prices = JSON.parse(
    localStorage.getItem(`game:${gameId}:scale`) ?? JSON.stringify(DEFAULT_PRICES),
  ) as number[]

  const { mutate: save, isPending } = useMutation({
    mutationFn: async () => {
      const existing = existingQuestions.data ?? []
      const orderNum = existing.length + 1

      return createQuestion(effectiveCatId, {
        price,
        type: 'standard',
        question: questionText.trim(),
        answer: answer.trim(),
        comment: comment.trim() || undefined,
        order_num: orderNum,
      })
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['questions-all'] })
      qc.invalidateQueries({ queryKey: ['questions', effectiveCatId] })
      navigate(`/game/${gameId}`, { replace: true })
    },
    onError: err => {
      setError(err instanceof Error ? err.message : 'Ошибка сохранения')
    },
  })

  const { mutate: saveAndAdd } = useMutation({
    mutationFn: async () => {
      const existing = existingQuestions.data ?? []
      const orderNum = existing.length + 1

      return createQuestion(effectiveCatId, {
        price,
        type: 'standard',
        question: questionText.trim(),
        answer: answer.trim(),
        comment: comment.trim() || undefined,
        order_num: orderNum,
      })
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['questions-all'] })
      qc.invalidateQueries({ queryKey: ['questions', effectiveCatId] })
      setQuestionText('')
      setAnswer('')
      setComment('')
      setError('')
    },
    onError: err => {
      setError(err instanceof Error ? err.message : 'Ошибка сохранения')
    },
  })

  function handleSave(e: React.FormEvent) {
    e.preventDefault()

    if (!questionText.trim()) { setError('Введите вопрос'); return }
    if (!answer.trim()) { setError('Введите ответ'); return }
    if (!effectiveCatId) { setError('Нет категорий — создайте игру заново'); return }

    setError('')
    save()
  }

  const loading = game.isLoading || rounds.isLoading || categories.isLoading

  // Progress indicator
  const totalExpected = prices.length * catList.length
  const allFilled = qc.getQueryData<Record<string, unknown[]>>(['questions-all'])
  const filledCount = allFilled
    ? Object.values(allFilled).reduce((s, qs) => s + qs.length, 0)
    : 0
  const progress = totalExpected > 0 ? filledCount / totalExpected : 0

  return (
    <div className="page">
      <div className="tgh">
        <button className="tgh-back" onClick={() => navigate(`/game/${gameId}`)}>
          <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="white" strokeWidth="2">
            <path d="M12 5l-7 5 7 5" />
          </svg>
        </button>
        <span className="tgh-title">Добавить вопрос</span>
        <div style={{ marginLeft: 'auto', display: 'flex', alignItems: 'center', gap: 6 }}>
          <span style={{ fontSize: 11, color: '#aaa' }}>{filledCount} / {totalExpected}</span>
          <div className="prog-bar">
            <div className="prog-fill" style={{ width: `${progress * 100}%` }} />
          </div>
        </div>
      </div>

      <div className="page-body">
        {loading ? (
          <div className="center" style={{ padding: 40 }}>
            <div className="spinner" />
          </div>
        ) : (
          <form onSubmit={handleSave} style={{ padding: '10px 0 20px' }}>
            <div className="wcard" style={{ margin: '0 8px 8px' }}>
              <span className="tlabel">Категория</span>
              <select
                className="tinput"
                style={{ paddingTop: 9, paddingBottom: 9, appearance: 'none', color: '#1a1a1a' }}
                value={effectiveCatId}
                onChange={e => setCatId(e.target.value)}
              >
                {catList.map(c => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>

              <span className="tlabel">Стоимость</span>
              <div
                className="pill-row"
                style={{ gridTemplateColumns: `repeat(${prices.length}, minmax(0, 1fr))` }}
              >
                {prices.map(p => (
                  <button
                    key={p}
                    type="button"
                    className={`ppill${price === p ? ' on' : ''}`}
                    onClick={() => setPrice(p)}
                  >
                    {p}
                  </button>
                ))}
              </div>

              <span className="tlabel">Текст вопроса</span>
              <textarea
                className="tinput"
                rows={3}
                placeholder="Введите вопрос…"
                value={questionText}
                onChange={e => setQuestionText(e.target.value)}
              />

              <span className="tlabel">Правильный ответ</span>
              <input
                className="tinput"
                placeholder="Введите ответ…"
                value={answer}
                onChange={e => setAnswer(e.target.value)}
              />

              <span className="tlabel">Подсказка для ведущего (необязательно)</span>
              <input
                className="tinput"
                placeholder="Доп. контекст или источник…"
                value={comment}
                onChange={e => setComment(e.target.value)}
              />

              {error && (
                <div style={{ fontSize: 12, color: '#c00', marginBottom: 8 }}>{error}</div>
              )}
            </div>

            <div style={{ display: 'flex', gap: 8, padding: '0 8px' }}>
              <button
                type="button"
                className="tbtn-ghost"
                style={{ flex: 1 }}
                disabled={isPending}
                onClick={() => {
                  if (!questionText.trim() || !answer.trim() || !effectiveCatId) return
                  setError('')
                  saveAndAdd()
                }}
              >
                + Ещё вопрос
              </button>
              <button type="submit" className="tbtn" style={{ flex: 1 }} disabled={isPending}>
                {isPending ? 'Сохраняем…' : 'Сохранить'}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  )
}
