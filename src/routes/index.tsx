import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useRef, useState } from 'react'
import { toast } from 'sonner'
import { error } from '@tauri-apps/plugin-log'

import { Config } from '@/providers/config-provider'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { SettingsSidebar } from '@/components/modals/settings/sidebar'
import { GeneralSettings } from '@/components/modals/settings/pages/general'
import { AboutSettings } from '@/components/modals/settings/pages/about'
import { GameSettings } from '@/components/modals/settings/pages/game'
import { DiscordSettings } from '@/components/modals/settings/pages/discord'
import { useConfig } from '@/hooks/use-config'

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

  const { config, isLoading, setConfig, saveConfig } = useConfig()

  const isInitialLoad = useRef(true)

  useEffect(() => {
    if (config === null || isLoading) {
      return
    }

    if (isInitialLoad.current) {
      isInitialLoad.current = false
      return
    }

    const handler = setTimeout(async () => {
      try {
        await saveConfig()
      } catch (err) {
        error('Failed to save settings:', {
          keyValues: {
            error: err instanceof Error ? err.message : String(err)
          }
        })
        toast.error('Failed to save settings. Please check logs and report this issue.')
      }
    }, 1000)

    return () => {
      clearTimeout(handler)
    }
  }, [config, isLoading, saveConfig])

  const renderActivePageContent = () => {
    const ComponentToRender = settingsPageComponents[activePage]

    if (!config || isLoading) {
      return <div className='p-4'>Loading...</div>
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
