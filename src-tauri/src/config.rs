use std::fs;
use std::path::PathBuf;

use figment::{
    providers::{Format, Serialized, Toml},
    Figment,
};
use serde::{Deserialize, Serialize};
use tauri::AppHandle;
use tauri_plugin_autostart::ManagerExt;

use crate::error::Error;
use crate::utils::{adjust_size, Rect, Resolution};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct UiCoordsConfig {
    pub esc: Rect,
    pub menu: Rect,
    pub sub_menu: Rect,
    pub combat: Rect,
    pub location: Rect,
    pub boss: Rect,
    pub characters: Vec<Rect>,
    pub characters_box: Vec<Rect>,
}

pub fn get_gui_coords(
    game_resolution: Resolution,
    x_adjustment: i32,
    y_adjustment: i32,
) -> UiCoordsConfig {
    let reference_resolution = Resolution {
        width: 2560,
        height: 1080,
    };
    let scale_x = game_resolution.width as f64 / reference_resolution.width as f64;
    let scale_y = game_resolution.height as f64 / reference_resolution.height as f64;

    let w_adjust = if game_resolution.width <= 1920 { 50 } else { 0 };

    UiCoordsConfig {
        esc: Rect {
            x: adjust_size(1925 + (x_adjustment * 2), scale_x),
            y: adjust_size(250 + y_adjustment, scale_y),
            width: adjust_size(180, scale_x) + 50,
            height: adjust_size(30, scale_y),
        },
        menu: Rect {
            x: adjust_size(100 + -(x_adjustment / 3), scale_x),
            y: adjust_size(35 + y_adjustment, scale_y),
            width: adjust_size(300, scale_x) + 50,
            height: adjust_size(40, scale_y),
        },
        sub_menu: Rect {
            x: adjust_size(100, scale_x),
            y: adjust_size(65 + y_adjustment, scale_y),
            width: adjust_size(370, scale_x) + 50,
            height: adjust_size(25, scale_y),
        },
        combat: Rect {
            x: adjust_size(2100 + x_adjustment, scale_x),
            y: adjust_size(25 + y_adjustment, scale_y),
            width: adjust_size(85, scale_x) + 50,
            height: adjust_size(40, scale_y),
        },
        location: Rect {
            x: adjust_size(55, scale_x),
            y: adjust_size(15 + y_adjustment, scale_y),
            width: adjust_size(320, scale_x) + 50,
            height: adjust_size(25, scale_y),
        },
        boss: Rect {
            x: adjust_size(850, scale_x),
            y: adjust_size(y_adjustment, scale_y),
            width: adjust_size(745, scale_x) + 50,
            height: adjust_size(50, scale_y),
        },
        characters: vec![
            Rect {
                x: adjust_size(2250 + x_adjustment, scale_x),
                y: adjust_size(305 + y_adjustment, scale_y),
                width: adjust_size(170, scale_x) + w_adjust,
                height: adjust_size(30, scale_y),
            },
            Rect {
                x: adjust_size(2250 + x_adjustment, scale_x),
                y: adjust_size(400 + y_adjustment, scale_y),
                width: adjust_size(170, scale_x) + w_adjust,
                height: adjust_size(30, scale_y),
            },
            Rect {
                x: adjust_size(2250 + x_adjustment, scale_x),
                y: adjust_size(495 + y_adjustment, scale_y),
                width: adjust_size(170, scale_x) + w_adjust,
                height: adjust_size(30, scale_y),
            },
            Rect {
                x: adjust_size(2250 + x_adjustment, scale_x),
                y: adjust_size(585 + y_adjustment, scale_y),
                width: adjust_size(170, scale_x) + w_adjust,
                height: adjust_size(30, scale_y),
            },
        ],
        characters_box: vec![
            Rect {
                x: adjust_size(2400 + x_adjustment, scale_x),
                y: adjust_size(351 + y_adjustment, scale_y),
                width: 1,
                height: 1,
            },
            Rect {
                x: adjust_size(2400 + x_adjustment, scale_x),
                y: adjust_size(445 + y_adjustment, scale_y),
                width: 1,
                height: 1,
            },
            Rect {
                x: adjust_size(2400 + x_adjustment, scale_x),
                y: adjust_size(538 + y_adjustment, scale_y),
                width: 1,
                height: 1,
            },
            Rect {
                x: adjust_size(2400 + x_adjustment, scale_x),
                y: adjust_size(632 + y_adjustment, scale_y),
                width: 1,
                height: 1,
            },
        ],
    }
}

impl Default for UiCoordsConfig {
    fn default() -> Self {
        get_gui_coords(Resolution::new(1920, 1080), 0, 0)
    }
}

#[derive(Debug, Serialize, Deserialize, Clone)]
#[serde(rename_all = "snake_case")]
pub enum ClosingBehavior {
    Exit,
    Minimize,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Config {
    pub window_name: String,
    pub window_class: String,
    pub resolution: Resolution,
    pub loop_time: u64,
    pub autostart: bool,
    pub tray_launch: bool,
    pub closing_behavior: ClosingBehavior,
    pub auto_update: bool,
    pub enable_status: bool,
    pub account_uid: Option<String>,
    pub account_name: Option<String>,
    pub display_name: bool,
    pub display_level: bool,
    pub preprocess_threshold: u32,
    pub discord_app_id: String,
    pub ui_coords: UiCoordsConfig,
    pub path: PathBuf,
}

impl Default for Config {
    fn default() -> Self {
        Self {
            window_name: "Star Rail".to_string(),
            window_class: "UnityWndClass".to_string(),
            resolution: Resolution::new(1920, 1080),
            loop_time: 2000,
            autostart: false,
            tray_launch: false,
            closing_behavior: ClosingBehavior::Minimize,
            auto_update: false,
            enable_status: true,
            account_uid: None,
            account_name: None,
            display_name: true,
            display_level: true,
            preprocess_threshold: 135,
            discord_app_id: "1208212792574869544".to_string(),
            ui_coords: UiCoordsConfig::default(),
            path: PathBuf::new(),
        }
    }
}

impl Config {
    pub fn load(path_str: &str, app: &AppHandle) -> Result<Self, Error> {
        let config_path = PathBuf::from(path_str);
        let autostart_manager = app.autolaunch();

        let auto_start_enabled = autostart_manager.is_enabled().unwrap_or(false);

        if !config_path.exists() {
            let mut default_config = Config {
                autostart: auto_start_enabled,
                ..Config::default()
            };

            let serialized_config = toml::to_string_pretty(&default_config)
                .map_err(|e| Error::Custom(format!("Failed to serialize default config: {}", e)))?;
            fs::write(&config_path, serialized_config).map_err(|e| {
                Error::Custom(format!(
                    "Failed to write default config to {}: {}",
                    config_path.display(),
                    e
                ))
            })?;
            default_config.path = config_path.clone();
            log::info!("Created default config file at: {}", config_path.display());
            return Ok(default_config);
        }

        let config = Figment::new()
            .merge(Serialized::defaults(Config::default()))
            .merge(Toml::file(&config_path));

        let mut config: Config = config
            .extract()
            .map_err(|e| Error::Custom(format!("Failed to load config: {}", e)))?;

        config.autostart = auto_start_enabled;
        config.path = config_path;

        config
            .save()
            .map_err(|e| Error::Custom(format!("Failed to save config: {}", e)))?;

        Ok(config)
    }

    pub fn save(&self) -> Result<(), Error> {
        let toml_string = toml::to_string_pretty(self)
            .map_err(|e| Error::Custom(format!("Failed to serialize config: {}", e)))?;
        fs::write(&self.path, toml_string).map_err(|e| {
            Error::Custom(format!(
                "Failed to write config to {}: {}",
                self.path.display(),
                e
            ))
        })?;
        Ok(())
    }
}
