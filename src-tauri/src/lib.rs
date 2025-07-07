use config::Config;
use tauri::menu::{Menu, MenuItem};
use tauri::path::BaseDirectory;
use tauri::tray::{MouseButton, MouseButtonState, TrayIconBuilder, TrayIconEvent};
use tauri::Manager;
use tauri_plugin_updater::UpdaterExt;

mod app;
mod config;
mod constants;
mod error;
mod game;
mod ocr;
mod utils;

async fn update(app: tauri::AppHandle) -> tauri_plugin_updater::Result<()> {
    if let Some(update) = app.updater()?.check().await? {
        let mut downloaded = 0;

        update
            .download_and_install(
                |chunk_length, content_length| {
                    downloaded += chunk_length;
                    println!("downloaded {downloaded} from {content_length:?}");
                },
                || {
                    println!("download finished");
                },
            )
            .await?;

        println!("update installed");
        app.restart();
    }

    Ok(())
}

fn start_app_loop(app_handle: tauri::AppHandle) {
    let config_dir = app_handle
        .path()
        .config_dir()
        .expect("Failed to get config directory");
    let config = Config::load(config_dir.join("config.toml").to_str().unwrap()).unwrap();

    let resources_path = app_handle
        .path()
        .resolve("game-data", BaseDirectory::Resource)
        .expect("Failed to resolve resources path");

    let game_data = game::data::GameData::load(&resources_path).expect("Failed to load game data");

    let tesseract_path = app_handle
        .path()
        .resolve("binaries/tesseract/tesseract.exe", BaseDirectory::Resource)
        .expect("failed to resolve tesseract path");

    let tessdata_path = app_handle
        .path()
        .resolve("binaries/tesseract/tessdata", BaseDirectory::Resource)
        .expect("failed to resolve tessdata path");

    // tauri::async_runtime::spawn_blocking(move || {
    //     let mut app = app::App::new(
    //         config,
    //         game_data,
    //         tesseract_path.to_str().unwrap(),
    //         tessdata_path.to_str().unwrap(),
    //         app_handle,
    //     );

    //     let rt = tokio::runtime::Builder::new_current_thread()
    //         .enable_all()
    //         .build()
    //         .expect("Failed to create Tokio runtime for App loop in blocking thread");
    //     rt.block_on(async move {
    //         app.start_loop().await;
    //     });
    // });
    std::thread::spawn(move || {
        let rt =
            tokio::runtime::Runtime::new().expect("Failed to create Tokio runtime for App loop");
        rt.block_on(async move {
            let mut app = app::App::new(
                config,
                game_data,
                tesseract_path.to_str().unwrap(),
                tessdata_path.to_str().unwrap(),
                app_handle,
            );
            app.start_loop().await;
        });
    });
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(
            tauri_plugin_log::Builder::new()
                .target(tauri_plugin_log::Target::new(
                    tauri_plugin_log::TargetKind::Stdout,
                ))
                .build(),
        )
        .plugin(tauri_plugin_updater::Builder::new().build())
        .plugin(tauri_plugin_single_instance::init(|app, args, cwd| {
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

            let handle = app.handle().clone();

            tauri::async_runtime::spawn(async move {
                update(handle).await.map_err(|e| {
                    log::error!("Failed to update: {}", e);
                }).ok();
            });

            start_app_loop(app.handle().clone());

            Ok(())
        })
        .plugin(tauri_plugin_opener::init())
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
