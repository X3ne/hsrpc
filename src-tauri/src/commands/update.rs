use tauri::AppHandle;
use tauri_plugin_updater::UpdaterExt;

#[derive(Debug, serde::Serialize, serde::Deserialize)]
pub struct UpdatePayload {
    version: String,
    notes: String,
    pub_date: String,
}

#[tauri::command] // TODO: use app error instead of String
pub async fn check_for_updates(app_handle: AppHandle) -> Result<Option<UpdatePayload>, String> {
    let updater = app_handle.updater().map_err(|e| e.to_string())?;
    if let Some(update) = updater.check().await.map_err(|e| {
        log::error!("Failed to check for updates: {}", e);
        e.to_string()
    })? {
        let update_info = serde_json::from_value::<UpdatePayload>(update.raw_json)
            .map_err(|e| format!("Failed to parse update info: {}", e))?;

        return Ok(Some(update_info));
    }

    Ok(None)
}

#[derive(Clone, serde::Serialize)]
#[serde(
    rename_all = "camelCase",
    rename_all_fields = "camelCase",
    tag = "event",
    content = "data"
)]
pub enum DownloadEvent {
    Started,
    Progress {
        downloaded: u64,
        content_length: Option<u64>,
    },
    Finished,
}

#[tauri::command]
pub async fn download_and_install_update(
    app_handle: AppHandle,
    on_download: tauri::ipc::Channel<DownloadEvent>,
) -> Result<(), String> {
    let updater = app_handle.updater().map_err(|e| {
        log::error!("Failed to get updater: {}", e);
        e.to_string()
    })?;

    if let Some(update) = updater.check().await.map_err(|e| {
        log::error!("Failed to check for updates: {}", e);
        e.to_string()
    })? {
        let mut downloaded: usize = 0;

        on_download
            .send(DownloadEvent::Started)
            .unwrap_or_else(|e| {
                log::error!("Failed to send download started event: {}", e);
            });

        update
            .download_and_install(
                |chunk_length, content_length| {
                    downloaded += chunk_length;

                    on_download
                        .send(DownloadEvent::Progress {
                            downloaded: downloaded as u64,
                            content_length,
                        })
                        .unwrap_or_else(|e| {
                            log::error!("Failed to send download progress: {}", e);
                        });
                },
                || {
                    on_download
                        .send(DownloadEvent::Finished)
                        .unwrap_or_else(|e| {
                            log::error!("Failed to send download finished event: {}", e);
                        });
                },
            )
            .await
            .map_err(|e| {
                log::error!("Failed to download and install update: {}", e);
                e.to_string()
            })?;

        app_handle.restart();
    }

    Ok(())
}
