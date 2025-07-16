use discord_rich_presence::DiscordIpc;
use tauri::{AppHandle, Manager};

use crate::AppState;

#[tauri::command]
pub async fn reconnect_to_discord(app_handle: AppHandle) -> Result<(), String> {
    let state = app_handle.state::<AppState>();
    let mut discord_ipc_state = state.discord_ipc_state.lock().map_err(|e| e.to_string())?;

    match discord_ipc_state.ipc_client.connect() {
        Ok(_) => {
            discord_ipc_state.connected = true;
            log::info!("Discord IPC connected successfully");
        }
        Err(e) => {
            log::error!("Failed to connect to Discord IPC: {}", e);
            discord_ipc_state.connected = false;
            return Err(format!("Failed to connect to Discord IPC: {}", e));
        }
    }
    Ok(())
}

#[tauri::command]
pub async fn is_ipc_connected(app_handle: AppHandle) -> Result<bool, String> {
    let state = app_handle.state::<AppState>();
    let discord_ipc_state = state.discord_ipc_state.lock().map_err(|e| e.to_string())?;

    Ok(discord_ipc_state.connected)
}
