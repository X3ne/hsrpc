import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useRef, useState } from 'react'
import { toast } from 'sonner'
import { invoke } from '@tauri-apps/api/core'
import { error } from '@tauri-apps/plugin-log'

import { Config } from '@/types'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { SettingsSidebar } from '@/components/modals/settings/sidebar'
import { GeneralSettings } from '@/components/modals/settings/pages/general'
import { AboutSettings } from '@/components/modals/settings/pages/about'
import { GameSettings } from '@/components/modals/settings/pages/game'
import { DiscordSettings } from '@/components/modals/settings/pages/discord'

type SettingsPageKey = 'general' | 'game' | 'discord' | 'about'

const settingsPageComponents: Record<
  SettingsPageKey,
  React.ComponentType<{
    config: Config
    onConfigChange: React.Dispatch<React.SetStateAction<Config | null>>
  }>
> = {
  general: GeneralSettings,
  game: GameSettings,
  discord: DiscordSettings,
  about: AboutSettings
}

export const Route = createFileRoute('/')({
  component: RouteComponent
})

function RouteComponent() {
  const [activePage, setActivePage] = useState<SettingsPageKey>('general')
  const [config, setConfig] = useState<Config | null>(null)

  const isInitialLoad = useRef(true)

  useEffect(() => {
    const loadConfig = async () => {
      try {
        const loadedConfig = await invoke<Config>('load_config_command')
        setConfig(loadedConfig)
        console.debug('Config loaded:', loadedConfig)
      } catch (err) {
        error('Failed to load config:', {
          keyValues: {
            error: err instanceof Error ? err.message : String(err)
          }
        })
      }
    }
    loadConfig()
  }, [])

  useEffect(() => {
    if (config === null) {
      return
    }

    if (isInitialLoad.current) {
      isInitialLoad.current = false
      return
    }

    const handler = setTimeout(async () => {
      try {
        await invoke('save_config_command', { newConfig: config })
        console.debug('Config autosaved:', config)
      } catch (e) {
        error('Autosave failed:', {
          keyValues: {
            error: e instanceof Error ? e.message : String(e)
          }
        })
        toast.error('Failed to autosave settings. Please check logs and report this issue.')
      }
    }, 1000)

    return () => {
      clearTimeout(handler)
    }
  }, [config])

  const renderActivePageContent = () => {
    const ComponentToRender = settingsPageComponents[activePage]

    // TODO: refactor this
    if (!config) {
      return <div className='flex items-center justify-center h-full'>Loading...</div>
    }

    return <ComponentToRender config={config} onConfigChange={setConfig} />
  }

  return (
    <div className='p-0 h-full w-full'>
      <SidebarProvider className='min-h-0'>
        <SettingsSidebar activePage={activePage} onPageChange={setActivePage} className='h-full' />
        <SidebarInset className='h-screen'>{renderActivePageContent()}</SidebarInset>
      </SidebarProvider>
    </div>
  )
}
