import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import App from './App'
import './styles.css'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { retry: 1, staleTime: 30_000 },
  },
})

type SafeAreaInset = { top?: number; right?: number; bottom?: number; left?: number }
type TgWebApp = {
  ready?: () => void
  expand?: () => void
  requestFullscreen?: () => void
  disableVerticalSwipes?: () => void
  safeAreaInset?: SafeAreaInset
  contentSafeAreaInset?: SafeAreaInset
  onEvent?: (event: string, handler: () => void) => void
  viewportHeight?: number
  viewportStableHeight?: number
}

declare global {
  interface Window {
    Telegram?: { WebApp?: TgWebApp }
  }
}

function applySafeArea(tg: TgWebApp) {
  const sa = tg.safeAreaInset ?? {}
  const csa = tg.contentSafeAreaInset ?? {}
  const root = document.documentElement.style
  root.setProperty('--tg-safe-top', `${sa.top ?? 0}px`)
  root.setProperty('--tg-safe-bottom', `${sa.bottom ?? 0}px`)
  root.setProperty('--tg-content-safe-top', `${csa.top ?? 0}px`)
  root.setProperty('--tg-content-safe-bottom', `${csa.bottom ?? 0}px`)
}

const tg = window.Telegram?.WebApp
if (tg) {
  tg.ready?.()
  tg.expand?.()
  tg.requestFullscreen?.()
  tg.disableVerticalSwipes?.()
  applySafeArea(tg)
  tg.onEvent?.('safeAreaChanged', () => applySafeArea(tg))
  tg.onEvent?.('contentSafeAreaChanged', () => applySafeArea(tg))
  tg.onEvent?.('fullscreenChanged', () => applySafeArea(tg))
  tg.onEvent?.('viewportChanged', () => applySafeArea(tg))
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>,
)
