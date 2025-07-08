/**
 * Represents a rectangle with integer coordinates and dimensions.
 */
export interface Rect {
  x: number
  y: number
  width: number
  height: number
}

export interface Resolution {
  width: number
  height: number
}

/**
 * Defines the UI coordinates for various in-game elements
 */
export interface UiCoordsConfig {
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
 * Example: "quit" | "minimize"
 */
export type ClosingBehavior = 'exit' | 'minimize'

/**
 * Main application configuration
 */
export interface Config {
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
