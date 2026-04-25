import { useEffect, useRef } from 'react'
import type { Game, GameBoard } from '../api/types'

export interface GameStateEvent {
  game: Game
  board: GameBoard
}

export function useGameEvents(
  gameId: string,
  onEvent: (state: GameStateEvent) => void,
) {
  const cbRef = useRef(onEvent)
  cbRef.current = onEvent

  useEffect(() => {
    const es = new EventSource(`/admin/games/${gameId}/events`)

    es.onmessage = (e) => {
      try {
        cbRef.current(JSON.parse(e.data) as GameStateEvent)
      } catch {
        // ignore parse errors
      }
    }

    return () => es.close()
  }, [gameId])
}
