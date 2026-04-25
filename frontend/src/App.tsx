import { Routes, Route, Navigate } from 'react-router-dom'
import SetupPage from './pages/SetupPage'
import PackListPage from './pages/PackListPage'
import CreateGamePage from './pages/CreateGamePage'
import GameBoardPage from './pages/GameBoardPage'
import QuestionPage from './pages/QuestionPage'
import AddQuestionPage from './pages/AddQuestionPage'

const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i

function useUserId(): string | null {
  const id = localStorage.getItem('userId')

  if (id && !UUID_RE.test(id)) {
    localStorage.removeItem('userId')
    localStorage.removeItem('userName')
    return null
  }

  return id
}

function RequireUser({ children }: { children: React.ReactNode }) {
  const userId = useUserId()

  if (!userId) return <Navigate to="/setup" replace />

  return <>{children}</>
}

export default function App() {
  return (
    <Routes>
      <Route path="/setup" element={<SetupPage />} />
      <Route
        path="/"
        element={
          <RequireUser>
            <PackListPage />
          </RequireUser>
        }
      />
      <Route
        path="/game/create"
        element={
          <RequireUser>
            <CreateGamePage />
          </RequireUser>
        }
      />
      <Route
        path="/game/:gameId"
        element={
          <RequireUser>
            <GameBoardPage />
          </RequireUser>
        }
      />
      <Route
        path="/game/:gameId/question/:questionId"
        element={
          <RequireUser>
            <QuestionPage />
          </RequireUser>
        }
      />
      <Route
        path="/game/:gameId/question/add"
        element={
          <RequireUser>
            <AddQuestionPage />
          </RequireUser>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
