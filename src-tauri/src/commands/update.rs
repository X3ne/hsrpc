use tauri::AppHandle;
use tauri_plugin_updater::UpdaterExt;

#[derive(Debug, serde::Serialize, serde::Deserialize)]
pub struct UpdatePayload {
    version: String,
    notes: String,
    pub_date: String,
}

#[tauri::command]
pub async fn check_for_updates(app_handle: AppHandle) -> Result<Option<UpdatePayload>, String> {
    let updater = app_handle.updater().map_err(|e| e.to_string())?;
    if let Some(update) = updater.check().await.map_err(|e| e.to_string())? {
        let update_info = serde_json::from_value::<UpdatePayload>(update.raw_json)
            .map_err(|e| format!("Failed to parse update info: {}", e))?;

        return Ok(Some(update_info));
    }

    Ok(None)
}
