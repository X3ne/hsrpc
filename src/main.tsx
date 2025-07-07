import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { getCurrentWindow } from '@tauri-apps/api/window'
import { ThemeProvider } from '@/components/theme-provider'
import './index.css'

import { routeTree } from './routeTree.gen'

const router = createRouter({ routeTree })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

const rootElement = document.getElementById('root')!
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement)
  root.render(
    <StrictMode>
      <ThemeProvider defaultTheme='dark' storageKey='ui-theme'>
        <RouterProvider router={router} />
      </ThemeProvider>
    </StrictMode>
  )

  setTimeout(async () => {
    const window = getCurrentWindow()
    await window.show()
    await window.setFocus()
  }, 0)
}
