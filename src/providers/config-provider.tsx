import { invoke } from '@tauri-apps/api/core'
import React from 'react'

interface Rect {
  x: number
  y: number
  width: number
  height: number
}

interface Resolution {
  width: number
  height: number
}

/**
 * Defines the UI coordinates for various in-game elements
 */
interface UiCoordsConfig {
  esc: Rect
  menu: Rect
  sub_menu: Rect
  combat: Rect
  location: Rect
  boss: Rect
  characters: Rect[]
  characters_box: Rect[]
}

/**
 * Defines the application's behavior when the window is closed
 * Example: "exit" | "minimize"
 */
type ClosingBehavior = 'exit' | 'minimize'

/**
 * Main application configuration
 */
interface Config {
  window_name: string
  window_class: string
  resolution: Resolution
  loop_time: number
  autostart: boolean
  tray_launch: boolean
  closing_behavior: ClosingBehavior
  auto_update: boolean
  enable_status: boolean
  account_uid: string | null
  account_name: string | null
  display_name: boolean
  display_level: boolean
  preprocess_threshold: number
  discord_app_id: string
  ui_coords: UiCoordsConfig
}

interface ConfigContextType {
  config: Config | null
  setConfig: React.Dispatch<React.SetStateAction<Config | null>>
  isLoading: boolean
  isSaving: boolean
  saveConfig: () => Promise<void>
  loadConfig: () => Promise<void>
}

const ConfigContext = React.createContext<ConfigContextType>({
  config: null,
  setConfig: () => {},
  isLoading: false,
  isSaving: false,
  saveConfig: async () => {},
  loadConfig: async () => {}
})

const ConfigProvider = ({ children }: { children: React.ReactNode }) => {
  const [config, setConfig] = React.useState<Config | null>(null)
  const [isLoading, setIsLoading] = React.useState(false)
  const [isSaving, setIsSaving] = React.useState(false)

  const saveConfig = React.useCallback(async () => {
    setIsSaving(true)
    try {
      await invoke('save_config_command', { newConfig: config })
      console.debug('Config saved:', config)
    } catch (error) {
      console.error('Failed to save config:', error)
      throw error
    } finally {
      setIsSaving(false)
    }
  }, [config])

  const loadConfig = React.useCallback(async () => {
    setIsLoading(true)
    try {
      const loadedConfig = await invoke<Config>('load_config_command')
      setConfig(loadedConfig)
      console.debug('Config loaded:', loadedConfig)
    } catch (error) {
      console.error('Failed to load config:', error)
      throw error
    } finally {
      setIsLoading(false)
    }
  }, [])

  return (
    <ConfigContext.Provider
      value={{
        config,
        setConfig,
        isLoading,
        isSaving,
        saveConfig,
        loadConfig
      }}
    >
      {children}
    </ConfigContext.Provider>
  )
}

export {
  ConfigContext,
  ConfigProvider,
  type Config,
  type Rect,
  type Resolution,
  type UiCoordsConfig,
  type ClosingBehavior
}
