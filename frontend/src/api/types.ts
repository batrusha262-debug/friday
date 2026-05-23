export interface User {
  id: string
  username: string
  email?: string
  role: 'admin' | 'guest'
  created_at: string
}

export interface Session {
  token: string
  user: User
}

export interface Pack {
  id: string
  title: string
  author_id: string
  created_at: string
}

export interface Round {
  id: string
  pack_id: string
  name: string
  type: 'standard' | 'double' | 'final'
  order_num: number
}

export interface Category {
  id: string
  round_id: string
  name: string
  order_num: number
}

export interface Question {
  id: string
  category_id: string
  price: number
  type: 'standard' | 'auction' | 'cat_in_bag' | 'no_risk'
  question: string
  answer: string
  comment?: string
  media_url?: string
  order_num: number
  options: string[]
  correct_option: number
}

export interface Game {
  id: string
  pack_id: string
  host_id: string
  status: 'waiting' | 'active' | 'finished'
  is_open: boolean
  created_at: string
  started_at?: string
  finished_at?: string
  current_picker_id?: string
}

export interface GameTeam {
  id: string
  game_id: string
  name: string
  score: number
  order_num: number
}

export interface GameQuestionState {
  id: string
  game_id: string
  question_id: string
  answered_by?: string
  answered_at?: string
  revealed_count: number
  timer_started_at?: string
  wrong_options: number[]
}

export interface AnswerClaim {
  id: string
  game_id: string
  question_id: string
  team_id: string
  claimed_at: string
  status: 'pending' | 'approved' | 'rejected'
  reviewed_at?: string
}

export interface MiniGame {
  id: string
  game_id: string
  question_id: string
  excluded_team_id?: string
  pos_x: number
  pos_y: number
  started_at: string
  appears_at: string
  winner_team_id?: string
  finished_at?: string
}

export interface Lobby {
  game_id: string
  pack_id: string
  pack_title: string
  team_count: number
  is_open: boolean
}

export interface GameBoard {
  teams: GameTeam[]
  states: GameQuestionState[]
  pending_claims: AnswerClaim[]
  mini_game?: MiniGame
}
