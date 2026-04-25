export interface User {
  id: string
  username: string
  created_at: string
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
}

export interface Game {
  id: string
  pack_id: string
  host_id: string
  status: 'waiting' | 'active' | 'finished'
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
}

export interface GameBoard {
  teams: GameTeam[]
  states: GameQuestionState[]
}
