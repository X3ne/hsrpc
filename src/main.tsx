import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

import { ConfigProvider } from '@/providers/config-provider'
import { ThemeProvider } from '@/components/theme-provider'
import './index.css'

import { routeTree } from './routeTree.gen'

const router = createRouter({ routeTree })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

const queryClient = new QueryClient()

document.addEventListener('keydown', function (e) {
  if (e.key === 'Tab') {
    e.preventDefault()
  }
})

const rootElement = document.getElementById('root')!
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement)
  root.render(
    <StrictMode>
      <QueryClientProvider client={queryClient}>
        <ThemeProvider defaultTheme='dark' storageKey='ui-theme'>
          <ConfigProvider>
            <RouterProvider router={router} />
          </ConfigProvider>
        </ThemeProvider>
      </QueryClientProvider>
    </StrictMode>
  )
}
