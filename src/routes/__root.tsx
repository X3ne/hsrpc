import React from 'react'
import { createRootRoute, Outlet } from '@tanstack/react-router'
import { error } from '@tauri-apps/plugin-log'

import { Layout } from '@/components/layout/layout'
import { useConfig } from '@/hooks/use-config'
import { toast } from 'sonner'

export const Route = createRootRoute({
  component: Root
})

function Root() {
  const { loadConfig } = useConfig()

  React.useEffect(() => {
    const load = async () => {
      try {
        await loadConfig()
      } catch (e) {
        error('Failed to load config:', {
          keyValues: {
            error: e instanceof Error ? e.message : String(e)
          }
        })
        toast.error('Failed to load settings. Please check logs and report this issue.')
      }
    }
    load()
  }, [loadConfig])

  return (
    <>
      <Layout>
        <Outlet />
      </Layout>
    </>
  )
}
