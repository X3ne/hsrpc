import React from 'react'
import { createRootRoute, Outlet } from '@tanstack/react-router'
import { error } from '@tauri-apps/plugin-log'
import { getCurrentWindow } from '@tauri-apps/api/window'

import { Layout } from '@/components/layout/layout'
import { useConfig } from '@/hooks/use-config'
import { toast } from 'sonner'

export const Route = createRootRoute({
  component: Root
})

function Root() {
  const { config, loadConfig } = useConfig()

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

  React.useEffect(() => {
    if (config === null) {
      return
    }
    const handler = setTimeout(async () => {
      if (!config?.tray_launch) {
        const window = getCurrentWindow()
        await window.show()
        await window.setFocus()
      }
    }, 0)
    return () => {
      clearTimeout(handler)
    }
  }, [config])

  return (
    <>
      <Layout>
        <Outlet />
      </Layout>
    </>
  )
}
