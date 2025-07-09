use std::sync::{Arc, Mutex};

use discord_rich_presence::{activity, DiscordIpc};
use serde::Serialize;
use tauri::{Emitter, Manager};
use tokio::time::{sleep, Duration};

use crate::config::{get_gui_coords, Config};
use crate::constants::LOOP_RETRY_TIMEOUT;
use crate::game::data::{Data, GameData};
use crate::ocr::windows::WindowsOcrManager;
use crate::ocr::{GameOcrJob, Lang, OcrManager, PreprocessOptions};
use crate::utils::{find_closest_correspondence, find_current_character};
use crate::AppState;

pub struct App {
    pub config: Arc<Mutex<Config>>,
    pub ocr_manager: WindowsOcrManager,
    pub state: State,
    pub app_started: chrono::DateTime<chrono::Utc>,
    pub game_data: GameData,
    pub app_handle: tauri::AppHandle,
}

#[derive(Debug, Eq, PartialEq, Clone, Serialize)]
#[serde(rename_all = "snake_case")]
pub enum State {
    Idle,
    Location {
        location: Data,
        character: Data,
    },
    Menu {
        menu: Data,
    },
    Combat {
        started: chrono::DateTime<chrono::Utc>,
        boss: Option<Data>,
    },
}

impl App {
    pub fn new(
        config: Arc<Mutex<Config>>,
        game_data: GameData,
        tesseract_path: &str,
        tessdata_path: &str,
        app_handle: tauri::AppHandle,
    ) -> Self {
        #[cfg(target_os = "windows")]
        let ocr_manager = WindowsOcrManager::new(tesseract_path, tessdata_path, Lang::Eng);

        #[cfg(not(target_os = "windows"))]
        panic!("Unsupported OS");

        App {
            config,
            ocr_manager,
            state: State::Idle,
            app_started: chrono::Utc::now(),
            game_data,
            app_handle,
        }
    }

    pub async fn start_loop(&mut self) {
        loop {
            let config;
            {
                let config_guard = self
                    .config
                    .lock()
                    .map_err(|e| {
                        log::error!("Failed to lock config: {}", e);
                    })
                    .unwrap();
                config = config_guard.clone();
            }

            if !config.enable_status {
                sleep(Duration::from_millis(config.loop_time)).await;
                continue;
            }

            if self
                .ocr_manager
                .set_game_window(&config.window_name)
                .is_err()
            {
                self.handle_window_closed();
                sleep(Duration::from_millis(LOOP_RETRY_TIMEOUT)).await;
                continue;
            }

            if let Some(ref account_name) = config.account_name {
                self.game_data.characters.insert(
                    "char_trailblazer".to_owned(),
                    Data {
                        asset_id: "char_trailblazer".to_string(),
                        value: account_name.clone(),
                        ..Default::default()
                    },
                );
            }

            match self.ocr_manager.is_initialized() {
                false => {
                    log::info!("Game window found: '{}'", config.window_name);
                    self.app_started = chrono::Utc::now();
                }
                true => {
                    self.set_game_resolution();

                    let old_state = self.state.clone();
                    if let Some(new_state) = self.capture_game_data(&config) {
                        self.state = new_state;
                    }

                    if old_state != self.state {
                        log::info!(
                            "App state changed or Discord client connected: {:?} -> {:?}",
                            old_state,
                            self.state
                        );
                        self.app_handle
                            .emit("app-state", &self.state)
                            .unwrap_or_else(|e| log::error!("Failed to emit app-state: {}", e));

                        if let Err(err) = self.update_discord_presence() {
                            log::error!("Failed to update Discord presence: {}", err);
                        }
                    } else {
                        log::trace!("No state change. Skipping presence update");
                    }

                    tokio::time::sleep(Duration::from_millis(config.loop_time)).await;
                }
            }
        }
    }

    fn handle_window_closed(&mut self) {
        if self.ocr_manager.is_initialized() {
            self.ocr_manager.pause_ocr();
            self.state = State::Idle;

            let state = self.app_handle.state::<AppState>();
            let mut discord_ipc_state = state.discord_ipc_state.lock().unwrap();

            discord_ipc_state
                .ipc_client
                .clear_activity()
                .unwrap_or_else(|e| {
                    log::error!("Failed to clear Discord activity: {}", e);
                });

            log::info!("Game window closed");
        }
    }

    fn set_game_resolution(&mut self) {
        let res = match self.ocr_manager.get_window_size() {
            Ok(size) => size,
            Err(e) => {
                log::error!("Failed to get window size: {}", e);
                self.handle_window_closed();
                return;
            }
        };

        let mut config_guard = self
            .config
            .lock()
            .map_err(|e| {
                log::error!("Failed to lock config: {}", e);
            })
            .unwrap();

        if config_guard.resolution.width == res.width
            && config_guard.resolution.height == res.height
        {
            return;
        }

        config_guard.resolution.width = res.width;
        config_guard.resolution.height = res.height;

        config_guard.ui_coords = get_gui_coords(config_guard.resolution.clone(), 0, 0);

        if let Err(e) = config_guard.save() {
            log::error!("Failed to save config: {}", e);
        }
    }

    fn capture_game_data(&mut self, config: &Config) -> Option<State> {
        if let Some(location_state) = self.capture_location(config) {
            return Some(location_state);
        }

        if let Some(menu_state) = self.capture_game_menu(config) {
            return Some(menu_state);
        }

        if let Some(combat_state) = self.capture_combat(config) {
            return Some(combat_state);
        }

        None
    }

    fn capture_location(&self, config: &Config) -> Option<State> {
        match self
            .ocr_manager
            .game_ocr(config.ui_coords.location, GameOcrJob::Location, None)
        {
            Ok(text) => {
                log::debug!("Location OCR raw result: '{}'", text);

                if text.is_empty() {
                    return None;
                }

                let location =
                    match find_closest_correspondence(&text, self.game_data.locations.values()) {
                        Some(loc) => {
                            log::debug!(
                                "Location OCR matched: '{}' (distance: {})",
                                loc.value,
                                strsim::levenshtein(&text, &loc.value)
                            );
                            loc
                        }
                        None => {
                            log::debug!(
                                "No sufficiently close correspondence found for location: '{}'",
                                text
                            );
                            return None;
                        }
                    };

                let mut character_data = self.capture_character(config);
                if character_data.is_none() && location.sub_region == "Astral Express" {
                    let char_name = match config.account_name.clone() {
                        Some(name) if config.display_name => name,
                        _ => "Trailblazer".to_string(),
                    };
                    character_data = Some(Data {
                        asset_id: "char_trailblazer".to_string(),
                        value: char_name,
                        ..Default::default()
                    });
                }

                let character = character_data?;

                let state = State::Location {
                    location,
                    character,
                };

                Some(state)
            }
            Err(e) => {
                log::error!("Failed to perform OCR for location: {}", e);
                None
            }
        }
    }

    fn capture_character(&self, config: &Config) -> Option<Data> {
        let current_character_index =
            find_current_character(&self.ocr_manager, &config.ui_coords.characters_box);

        log::debug!("Current character index: {}", current_character_index);

        if current_character_index == -1 {
            log::debug!("No current character detected");
            return None;
        }

        let char_rect = config
            .ui_coords
            .characters
            .get(current_character_index as usize)?;

        let text = match self
            .ocr_manager
            .game_ocr(*char_rect, GameOcrJob::Character, None)
        {
            Ok(t) => t,
            Err(e) => {
                log::error!("Failed to perform OCR for character: {}", e);
                return None;
            }
        };

        if text.is_empty() {
            return None;
        }

        match find_closest_correspondence(&text, self.game_data.characters.values()) {
            Some(mut char_data) => {
                log::debug!(
                    "Character OCR matched: '{}' (distance: {})",
                    char_data.value,
                    strsim::levenshtein(&text, &char_data.value)
                );

                if char_data.asset_id == "char_trailblazer" {
                    if let Some(account_name) = &config.account_name {
                        char_data.value = if config.display_name {
                            account_name.clone()
                        } else {
                            "Trailblazer".to_string()
                        };
                    }
                }

                Some(char_data)
            }
            None => {
                log::debug!(
                    "No sufficiently close correspondence found for character: '{}'",
                    text
                );
                None
            }
        }
    }

    fn capture_game_menu(&self, config: &Config) -> Option<State> {
        let esc_text_result =
            self.ocr_manager
                .game_ocr(config.ui_coords.esc, GameOcrJob::Menu, None);
        let menu_text_result = self.ocr_manager.game_ocr(
            config.ui_coords.menu,
            GameOcrJob::Menu,
            Some(PreprocessOptions {
                threshold: config.preprocess_threshold,
            }),
        );

        let mut menu_data: Option<Data> = None;

        if let Ok(esc_text) = esc_text_result {
            log::debug!("ESC Menu OCR raw result: '{}'", esc_text);
            if !esc_text.is_empty() {
                if let Some(prediction) =
                    find_closest_correspondence(&esc_text, self.game_data.menus.values())
                {
                    log::debug!(
                        "ESC Menu OCR matched: '{}' (distance: {})",
                        prediction.value,
                        strsim::levenshtein(&esc_text, &prediction.value)
                    );
                    if prediction.value == "Trailblaze Level" {
                        menu_data = Some(prediction);
                    }
                } else {
                    log::debug!(
                        "No sufficiently close correspondence found for ESC Menu: '{}'",
                        esc_text
                    );
                }
            }
        }

        if menu_data.is_none() {
            if let Ok(menu_text) = menu_text_result {
                log::debug!("Main Menu OCR raw result: '{}'", menu_text);
                if !menu_text.is_empty() {
                    if let Some(prediction) =
                        find_closest_correspondence(&menu_text, self.game_data.menus.values())
                    {
                        log::debug!(
                            "Main Menu OCR matched: '{}' (distance: {})",
                            prediction.value,
                            strsim::levenshtein(&menu_text, &prediction.value)
                        );
                        menu_data = Some(prediction);
                    } else {
                        log::debug!(
                            "No sufficiently close correspondence found for Main Menu: '{}'",
                            menu_text
                        );
                    }
                }
            }
        }

        if let Some(mut menu) = menu_data {
            let sub_menu_text_result = self.ocr_manager.game_ocr(
                config.ui_coords.sub_menu,
                GameOcrJob::SubMenu,
                Some(PreprocessOptions {
                    threshold: config.preprocess_threshold,
                }),
            );
            if let Ok(sub_menu_text) = sub_menu_text_result {
                log::debug!("Sub-Menu OCR raw result: '{}'", sub_menu_text);
                if !sub_menu_text.is_empty() {
                    if let Some(sub_menu_prediction) = find_closest_correspondence(
                        &sub_menu_text,
                        self.game_data.sub_menus.values(),
                    ) {
                        log::debug!(
                            "Sub-Menu OCR matched: '{}' (distance: {})",
                            sub_menu_prediction.value,
                            strsim::levenshtein(&sub_menu_text, &sub_menu_prediction.value)
                        );
                        menu.sub_region = sub_menu_prediction.value;
                    } else {
                        log::debug!(
                            "No sufficiently close correspondence found for Sub-Menu: '{}'",
                            sub_menu_text
                        );
                    }
                }
            }
            Some(State::Menu { menu })
        } else {
            None
        }
    }

    fn capture_combat(&self, config: &Config) -> Option<State> {
        let combat_text_result =
            self.ocr_manager
                .game_ocr(config.ui_coords.combat, GameOcrJob::Combat, None);

        let started = match &self.state {
            State::Combat { started, .. } => *started,
            _ => chrono::Utc::now(),
        };

        if let Ok(combat_text) = combat_text_result {
            log::debug!("Combat OCR raw result: '{}'", combat_text);
            if !combat_text.is_empty() && combat_text.replace(" ", "").len() > 5 {
                return Some(State::Combat {
                    started,
                    boss: None,
                });
            } else {
                let boss_text_result = self.ocr_manager.game_ocr(
                    config.ui_coords.boss,
                    GameOcrJob::Boss,
                    Some(PreprocessOptions {
                        threshold: config.preprocess_threshold,
                    }),
                );

                if let Ok(boss_text) = boss_text_result {
                    log::debug!("Boss OCR raw result: '{}'", boss_text);
                    if !boss_text.is_empty() {
                        let boss_data =
                            find_closest_correspondence(&boss_text, self.game_data.bosses.values());
                        if let Some(boss) = boss_data {
                            log::debug!(
                                "Boss OCR matched: '{}' (distance: {})",
                                boss.value,
                                strsim::levenshtein(&boss_text, &boss.value)
                            );

                            return Some(State::Combat {
                                started,
                                boss: Some(boss),
                            });
                        } else {
                            log::debug!(
                                "No sufficiently close correspondence found for Boss: '{}'",
                                boss_text
                            );
                        }
                    }
                } else {
                    log::error!(
                        "Failed to perform OCR for boss: {:?}",
                        boss_text_result.err()
                    );
                }

                return None;
            }
        } else {
            log::error!(
                "Failed to perform OCR for combat: {:?}",
                combat_text_result.err()
            );
        }
        None
    }

    fn update_discord_presence(&mut self) -> Result<(), Box<dyn std::error::Error>> {
        let state = self.app_handle.state::<AppState>();
        let mut discord_ipc_state = state.discord_ipc_state.lock().unwrap();

        if !discord_ipc_state.connected {
            log::info!("Discord IPC not connected, trying to connect...");
            discord_ipc_state.ipc_client.connect()?;
            discord_ipc_state.connected = true;
            log::info!("Discord IPC connected successfully");
        }

        match &self.state {
            State::Idle => {
                log::debug!("Discord presence Idle");
            }
            State::Location {
                location,
                character,
            } => {
                let activity = activity::Activity::new()
                    .state(&location.sub_region)
                    .details(&location.value)
                    .assets(
                        activity::Assets::new()
                            .large_image(&location.asset_id)
                            .small_image(&character.asset_id)
                            .small_text(&character.value),
                    )
                    .timestamps(activity::Timestamps::new().start(self.app_started.timestamp()));
                discord_ipc_state.ipc_client.set_activity(activity)?;
                log::debug!("Discord presence updated: Location - {}", location.value);
            }
            State::Menu { menu } => {
                let mut activity = activity::Activity::new()
                    .state(&menu.message)
                    .details(&menu.value)
                    .assets(activity::Assets::new().large_image(&menu.asset_id))
                    .timestamps(activity::Timestamps::new().start(self.app_started.timestamp()));

                if menu.value == "Trailblaze Level" {
                    activity = activity.details("")
                }

                if !menu.sub_region.is_empty() {
                    activity = activity.details(&menu.sub_region);
                }

                discord_ipc_state.ipc_client.set_activity(activity)?;
                log::debug!("Discord presence updated: Menu - {}", menu.value);
            }
            State::Combat { started, boss } => {
                let mut activity = activity::Activity::new()
                    .state("In Combat")
                    .timestamps(activity::Timestamps::new().start(started.timestamp()));

                if let Some(boss_data) = boss {
                    activity = activity.details(&boss_data.value);
                    activity = activity.assets(
                        activity::Assets::new()
                            .large_image("menu_combat")
                            .small_image(&boss_data.asset_id),
                    );
                } else {
                    activity = activity.assets(activity::Assets::new().large_image("menu_combat"));
                }
                discord_ipc_state.ipc_client.set_activity(activity)?;
                log::debug!("Discord presence updated: Combat");
            }
        }

        Ok(())
    }
}
