import React from 'react'
import { createRootRoute, Outlet } from '@tanstack/react-router'
import { error } from '@tauri-apps/plugin-log'
import { getCurrentWindow } from '@tauri-apps/api/window'

import { Layout } from '@/components/layout/layout'
import { useConfig } from '@/hooks/use-config'
import { toast } from 'sonner'
import { listen } from '@tauri-apps/api/event'
import { invoke } from '@tauri-apps/api/core'
import { useUpdate } from '@/hooks/use-update'
import { Update } from '@/providers/update-provider'

export const Route = createRootRoute({
  component: Root
})

function Root() {
  const { config, loadConfig } = useConfig()
  const { setUpdate } = useUpdate()
  const ready = React.useRef(false)

  const showWindow = async () => {
    const window = getCurrentWindow()
    await window.show()
    await window.setFocus()
  }

  listen<Update>('update-available', e => {
    console.log('Update available event received', e.payload)
    showWindow()
    setUpdate(e.payload)
  })

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
    if (config === null || ready.current) {
      return
    }

    invoke('ready')

    if (!config?.tray_launch) {
      showWindow()
    }

    ready.current = true
  }, [config])

  return (
    <>
      <Layout>
        <Outlet />
      </Layout>
    </>
  )
}
