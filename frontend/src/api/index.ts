import { api } from './client'
import type {
  Category,
  Game,
  GameBoard,
  GameQuestionState,
  GameTeam,
  Pack,
  Question,
  Round,
  User,
} from './types'

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

export const startGame = (gameId: string) =>
  api.post<Game>(`/admin/games/${gameId}/start`, {})

export const finishGame = (gameId: string) =>
  api.post<Game>(`/admin/games/${gameId}/finish`, {})

// Teams
export const addTeam = (gameId: string, name: string) =>
  api.post<GameTeam>(`/admin/games/${gameId}/teams`, { name })

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
) =>
  api.post<GameQuestionState>(
    `/admin/games/${gameId}/questions/${questionId}/answer`,
    { team_id: teamId },
  )
