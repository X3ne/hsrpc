use tauri::{path::BaseDirectory, AppHandle, Manager};

#[tauri::command]
pub async fn open_log_file(app_handle: AppHandle) -> Result<(), String> {
    let config_dir = app_handle
        .path()
        .resolve("", BaseDirectory::AppLog)
        .map_err(|e| format!("Failed to resolve config directory: {}", e))?;

    if !config_dir.exists() {
        return Err(format!("Log dir not found at: {}", config_dir.display()));
    }

    tauri_plugin_opener::open_path(config_dir.to_string_lossy().into_owned(), None::<&str>)
        .map_err(|e| format!("Failed to open log file: {}", e))?;

    Ok(())
}
