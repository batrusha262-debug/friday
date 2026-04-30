import { Routes, Route, Navigate } from 'react-router-dom'
import LoginPage from './pages/LoginPage'
import PackListPage from './pages/PackListPage'
import CreateGamePage from './pages/CreateGamePage'
import GameBoardPage from './pages/GameBoardPage'
import QuestionPage from './pages/QuestionPage'
import AddQuestionPage from './pages/AddQuestionPage'

const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i

export function useAuth() {
  const token = localStorage.getItem('token')
  const userId = localStorage.getItem('userId')
  const role = localStorage.getItem('userRole') as 'admin' | 'guest' | null

  if (!token || !userId || !UUID_RE.test(userId)) {
    return { userId: null, role: null, token: null }
  }

  return { userId, role: role ?? 'guest', token }
}

function RequireUser({ children }: { children: React.ReactNode }) {
  const { userId } = useAuth()
  if (!userId) return <Navigate to="/login" replace />
  return <>{children}</>
}

function RequireAdmin({ children }: { children: React.ReactNode }) {
  const { userId, role } = useAuth()
  if (!userId) return <Navigate to="/login" replace />
  if (role !== 'admin') return <Navigate to="/" replace />
  return <>{children}</>
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      {/* Legacy redirect */}
      <Route path="/setup" element={<Navigate to="/login" replace />} />
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
          <RequireAdmin>
            <CreateGamePage />
          </RequireAdmin>
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
          <RequireAdmin>
            <AddQuestionPage />
          </RequireAdmin>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
