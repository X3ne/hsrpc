import { useEffect, useRef, useState } from 'react'
import { Setting2 } from 'iconsax-reactjs'
import { DialogTitle } from '@radix-ui/react-dialog'

import { Config } from '@/types'
import { Dialog, DialogContent, DialogTrigger } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { SettingsSidebar } from '@/components/modals/settings/sidebar'
import { GeneralSettings } from '@/components/modals/settings/pages/general'
import { AboutSettings } from '@/components/modals/settings/pages/about'
import { GameSettings } from '@/components/modals/settings/pages/game'
import { DiscordSettings } from '@/components/modals/settings/pages/discord'
import { invoke } from '@tauri-apps/api/core'
import { toast } from 'sonner'

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

const SettingsModal = () => {
  const [activePage, setActivePage] = useState<SettingsPageKey>('general')
  const [config, setConfig] = useState<Config | null>(null)
  const [isDialogOpen, setIsDialogOpen] = useState(false)

  const isInitialLoad = useRef(true)

  useEffect(() => {
    if (isDialogOpen) {
      const loadConfig = async () => {
        try {
          const loadedConfig = await invoke<Config>('load_config_command')
          setConfig(loadedConfig)
          console.log('Config loaded:', loadedConfig)
        } catch (err) {
          console.error('Failed to load config:', err)
        }
      }
      loadConfig()
    } else {
      setConfig(null)
      setActivePage('general')
    }
  }, [isDialogOpen])

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
        console.log('Config autosaved:', config)
      } catch (err) {
        console.error('Autosave failed:', err)
        toast.error('Failed to autosave settings. Please report this issue.')
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
    <Dialog onOpenChange={setIsDialogOpen} open={isDialogOpen}>
      <DialogTrigger>
        <Button variant='outline' size='icon'>
          <Setting2 size={18} />
        </Button>
      </DialogTrigger>
      <DialogContent className='p-0 h-9/12 w-3/4'>
        <DialogTitle className='sr-only'>Settings</DialogTitle>
        <SidebarProvider className='min-h-0'>
          <SettingsSidebar
            activePage={activePage}
            onPageChange={setActivePage}
            className='h-full'
          />
          <SidebarInset className='h-full rounded-lg'>{renderActivePageContent()}</SidebarInset>
        </SidebarProvider>
      </DialogContent>
    </Dialog>
  )
}

export { SettingsModal, type SettingsPageKey }
