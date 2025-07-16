use std::sync::{Arc, Mutex};

use config::Config;
use discord_rich_presence::DiscordIpcClient;
use tauri::menu::{Menu, MenuItem};
use tauri::path::BaseDirectory;
use tauri::tray::{MouseButton, MouseButtonState, TrayIconBuilder, TrayIconEvent};
use tauri::{Emitter, Manager};
use tauri_plugin_updater::UpdaterExt;

use crate::error::Error;

mod app;
mod commands;
mod config;
mod constants;
mod error;
mod game;
mod ocr;
mod utils;

pub struct IpcState {
    pub ipc_client: DiscordIpcClient,
    pub connected: bool,
}
pub struct AppState {
    pub config: Mutex<Config>,
    pub discord_ipc_state: Mutex<IpcState>,
}

#[derive(serde::Serialize, serde::Deserialize, Clone, Debug)]
struct Update {
    version: String,
    notes: String,
    pub_date: String,
}

async fn search_update(app: tauri::AppHandle) -> tauri_plugin_updater::Result<()> {
    if let Some(update) = app.updater()?.check().await? {
        log::info!(
            "Current version: {}, Update version: {}",
            update.current_version,
            update.version
        );
        if update.current_version == update.version {
            log::info!(
                "No updates available. Current version: {}",
                update.current_version
            );
            return Ok(());
        }

        log::info!("Update found: {}", update.version);

        let update_info = serde_json::from_value::<Update>(update.raw_json).map_err(|e| {
            log::error!("Failed to parse update info: {}", e);
            tauri_plugin_updater::Error::from(e)
        })?;

        app.emit(
            "update-available",
            Update {
                version: update_info.version,
                notes: update_info.notes,
                pub_date: update_info.pub_date,
            },
        )
        .map_err(tauri_plugin_updater::Error::Tauri)?;
    } else {
        log::info!("No updates available.");
    }

    Ok(())
}

#[tauri::command]
async fn ready(app_handle: tauri::AppHandle) -> Result<(), String> {
    tauri::async_runtime::spawn(async move {
        search_update(app_handle)
            .await
            .map_err(|e| {
                log::error!("Failed to search for update: {}", e);
            })
            .ok();
    });

    Ok(())
}

fn start_app_loop(app_handle: tauri::AppHandle) -> Result<(), Error> {
    let config_dir = app_handle
        .path()
        .resolve("", BaseDirectory::AppConfig)
        .map_err(|e| Error::Custom(format!("Failed to resolve config directory: {}", e)))?;

    if !config_dir.exists() {
        std::fs::create_dir_all(&config_dir)
            .map_err(|e| Error::AppFolderCreation(format!("Failed to create config dir: {}", e)))?;
    }

    let config = Config::load(
        config_dir.join("config.toml").to_str().unwrap(),
        &app_handle,
    )?;

    let discord_app_id = config.discord_app_id.clone();

    let discord_ipc_client = DiscordIpcClient::new(&discord_app_id)
        .map_err(|e| Error::Custom(format!("Failed to create Discord IPC client: {}", e)))?;

    app_handle.manage(AppState {
        config: Mutex::new(config),
        discord_ipc_state: Mutex::new(IpcState {
            ipc_client: discord_ipc_client,
            connected: false,
        }),
    });

    let resources_path = app_handle
        .path()
        .resolve("game-data", BaseDirectory::Resource)
        .map_err(|e| {
            Error::AppPathResolution(format!("Failed to resolve resources path: {}", e))
        })?;

    log::info!("Resources path: {}", resources_path.display());

    let game_data = game::data::GameData::load(&resources_path)
        .map_err(|e| Error::Custom(format!("Failed to load game data: {}", e)))?;

    let tesseract_path = app_handle
        .path()
        .resolve("binaries/tesseract/tesseract.exe", BaseDirectory::Resource)
        .map_err(|e| {
            Error::AppPathResolution(format!("Failed to resolve tesseract path: {}", e))
        })?;

    let tessdata_path = app_handle
        .path()
        .resolve("binaries/tesseract/tessdata", BaseDirectory::Resource)
        .map_err(|e| Error::AppPathResolution(format!("Failed to resolve tessdata path: {}", e)))?;

    std::thread::spawn(move || {
        let rt =
            tokio::runtime::Runtime::new().expect("Failed to create Tokio runtime for App loop");
        rt.block_on(async move {
            let mut app = app::App::new(
                game_data,
                tesseract_path.to_str().unwrap(),
                tessdata_path.to_str().unwrap(),
                app_handle,
            );
            app.start_loop().await;
        });
    });

    Ok(())
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(
            tauri_plugin_log::Builder::new()
                .max_file_size(50_000)
                .rotation_strategy(tauri_plugin_log::RotationStrategy::KeepOne)
                .build(),
        )
        .plugin(tauri_plugin_updater::Builder::new().build())
        .plugin(tauri_plugin_single_instance::init(|app, _, _| {
            let _ = app
                .get_webview_window("main")
                .expect("no main window")
                .set_focus();
        }))
        .setup(|app| {
            let quit_i = MenuItem::with_id(app, "quit", "Quit", true, None::<&str>)?;
            let menu = Menu::with_items(app, &[&quit_i])?;

            let _ = TrayIconBuilder::new()
                .icon(app.default_window_icon().unwrap().clone())
                .menu(&menu)
                .show_menu_on_left_click(false)
                .on_menu_event(|app, event| match event.id.as_ref() {
                    "quit" => {
                        app.exit(0);
                    }
                    _ => {
                        log::error!("menu item {:?} not handled", event.id);
                    }
                })
                .on_tray_icon_event(|tray, event| {
                    if let TrayIconEvent::Click {
                        button: MouseButton::Left,
                        button_state: MouseButtonState::Up,
                        ..
                    } = event
                    {
                        let app = tray.app_handle();
                        if let Some(window) = app.get_webview_window("main") {
                            let _ = window.show();
                            let _ = window.set_focus();
                        }
                    }
                })
                .build(app)?;

            #[cfg(desktop)]
            {
                use tauri_plugin_autostart::MacosLauncher;

                let _ = app.handle().plugin(tauri_plugin_autostart::init(
                    MacosLauncher::LaunchAgent,
                    None,
                ));
            }

            start_app_loop(app.handle().clone())
                .map_err(|e| {
                    e.report_and_exit(app.handle(), "Failed to start app");
                })
                .unwrap_or(());

            Ok(())
        })
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![
            ready,
            crate::commands::open_log_file,
            crate::commands::check_for_updates,
            crate::commands::download_and_install_update,
            crate::commands::load_config_command,
            crate::commands::save_config_command,
            crate::commands::reconnect_to_discord,
            crate::commands::is_ipc_connected,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
