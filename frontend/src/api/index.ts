import { api } from './client'
import type {
  AnswerClaim,
  Category,
  Game,
  GameBoard,
  GameQuestionState,
  GameTeam,
  Lobby,
  MiniGame,
  Pack,
  Question,
  Round,
  Session,
  User,
} from './types'

function ensureSession(data: Session): Session {
  if (
    !data ||
    typeof data.token !== 'string' ||
    !data.token ||
    !data.user ||
    typeof data.user.id !== 'string'
  ) {
    throw new Error('Некорректный ответ сервера авторизации')
  }

  return data
}

// Auth
export const requestCode = (email: string) =>
  api.post<{ ok: boolean }>('/auth/request-code', { email })

export const verifyCode = (email: string, code: string) =>
  api.post<Session>('/auth/verify-code', { email, code }).then(ensureSession)

export const guestLogin = (name: string) =>
  api.post<Session>('/auth/guest', { name }).then(ensureSession)

export const logout = () =>
  api.post<void>('/auth/logout', {})

// Users
export const createUser = (username: string) =>
  api.post<User>('/admin/users', { username })

export const listUsers = () =>
  api.get<User[]>('/admin/users')

// Packs
export const createPack = (title: string, author_id: string) =>
  api.post<Pack>('/admin/packs', { title, author_id })

export const listPacks = () =>
  api.get<Pack[]>('/admin/packs')

export const getPack = (packId: string) =>
  api.get<Pack>(`/admin/packs/${packId}`)

export const deletePack = (packId: string) =>
  api.delete(`/admin/packs/${packId}`)

// Rounds
export const createRound = (packId: string, name: string, type = 'standard') =>
  api.post<Round>(`/admin/packs/${packId}/rounds`, { name, type })

export const listRounds = (packId: string) =>
  api.get<Round[]>(`/admin/packs/${packId}/rounds`)

// Categories
export const createCategory = (roundId: string, name: string) =>
  api.post<Category>(`/admin/rounds/${roundId}/categories`, { name })

export const listCategories = (roundId: string) =>
  api.get<Category[]>(`/admin/rounds/${roundId}/categories`)

// Questions
export const createQuestion = (
  categoryId: string,
  data: {
    price: number
    type: string
    question: string
    answer: string
    comment?: string
    order_num: number
    options: string[]
    correct_option: number
  },
) => api.post<Question>(`/admin/categories/${categoryId}/questions`, data)

export const listQuestions = (categoryId: string) =>
  api.get<Question[]>(`/admin/categories/${categoryId}/questions`)

export const getQuestion = (questionId: string) =>
  api.get<Question>(`/admin/questions/${questionId}`)

export const updateQuestion = (
  questionId: string,
  data: Partial<Question>,
) => api.put<Question>(`/admin/questions/${questionId}`, data)

export const deleteQuestion = (questionId: string) =>
  api.delete(`/admin/questions/${questionId}`)

// Games
export const createGame = (pack_id: string, host_id: string) =>
  api.post<Game>('/admin/games', { pack_id, host_id })

export const getGame = (gameId: string) =>
  api.get<Game>(`/admin/games/${gameId}`)

export const getGameByPack = (packId: string) =>
  api.get<Game>(`/admin/packs/${packId}/game`)

export const findGameByCode = (code: string) =>
  api.get<Game>(`/admin/games/join/${code.toLowerCase()}`)

export const deleteGame = (gameId: string) =>
  api.delete(`/admin/games/${gameId}`)

export const startGame = (gameId: string) =>
  api.post<Game>(`/admin/games/${gameId}/start`, {})

export const finishGame = (gameId: string) =>
  api.post<Game>(`/admin/games/${gameId}/finish`, {})

export const setGameOpen = (gameId: string, open: boolean) =>
  api.patch<Game>(`/admin/games/${gameId}/open`, { open })

// Teams
export const addTeam = (gameId: string, name: string) =>
  api.post<GameTeam>(`/admin/games/${gameId}/teams`, { name })

export const joinGame = (gameId: string, name: string) =>
  api.post<GameTeam>(`/admin/games/${gameId}/join`, { name })

export const listTeams = (gameId: string) =>
  api.get<GameTeam[]>(`/admin/games/${gameId}/teams`)

export const removeTeam = (teamId: string) =>
  api.delete(`/admin/teams/${teamId}`)

// Board
export const getBoard = (gameId: string) =>
  api.get<GameBoard>(`/admin/games/${gameId}/board`)

export const answerQuestion = (
  gameId: string,
  questionId: string,
  teamId: string | null,
  wrongTeamId?: string | null,
  optionIdx?: number | null,
) =>
  api.post<GameQuestionState>(
    `/admin/games/${gameId}/questions/${questionId}/answer`,
    {
      team_id: teamId,
      wrong_team_id: wrongTeamId ?? null,
      option_idx: optionIdx ?? null,
    },
  )

export const openQuestion = (gameId: string, questionId: string) =>
  api.post<GameQuestionState>(
    `/admin/games/${gameId}/questions/${questionId}/open`,
    {},
  )

export const revealNextOption = (gameId: string, questionId: string) =>
  api.post<GameQuestionState>(
    `/admin/games/${gameId}/questions/${questionId}/reveal`,
    {},
  )

export const startQuestionTimer = (gameId: string, questionId: string) =>
  api.post<GameQuestionState>(
    `/admin/games/${gameId}/questions/${questionId}/timer`,
    {},
  )

export const claimMiniGame = (
  gameId: string,
  miniGameId: string,
  teamId: string,
) =>
  api.post<MiniGame>(
    `/admin/games/${gameId}/mini-games/${miniGameId}/claim`,
    { team_id: teamId },
  )

export const claimAnswer = (gameId: string, questionId: string, teamId: string) =>
  api.post<AnswerClaim>(
    `/admin/games/${gameId}/questions/${questionId}/claim`,
    { team_id: teamId },
  )

export const validateClaim = (claimId: string, approved: boolean) =>
  api.post<AnswerClaim>(`/admin/claims/${claimId}/validate`, { approved })

// Lobbies
export const listActiveLobbies = () =>
  api.get<Lobby[]>('/admin/lobbies')
