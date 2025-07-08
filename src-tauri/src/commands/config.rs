use tauri::Manager;

use crate::{config::Config, AppState};

#[tauri::command]
pub async fn load_config_command(app_handle: tauri::AppHandle) -> Result<Config, String> {
    let state = app_handle.state::<AppState>();
    let config_guard = state
        .config
        .lock()
        .map_err(|e| format!("Failed to lock config: {}", e))?;
    Ok(config_guard.clone())
}

#[tauri::command]
pub async fn save_config_command(
    app_handle: tauri::AppHandle,
    new_config: Config,
) -> Result<(), String> {
    let state = app_handle.state::<AppState>();
    let mut config_guard = state
        .config
        .lock()
        .map_err(|e| format!("Failed to lock config: {}", e))?;

    *config_guard = new_config.clone();

    config_guard
        .save()
        .map_err(|e| format!("Failed to save config: {}", e))?;

    Ok(())
}
