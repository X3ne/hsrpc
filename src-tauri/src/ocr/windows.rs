use std::collections::HashMap;

use rusty_tesseract::{Args, Image};
use xcap::image::DynamicImage;
use xcap::Window;

use crate::error::Error;
use crate::ocr::{Lang, OcrManager};
use crate::utils::{Rect, Resolution};

#[derive(Debug)]
pub struct WindowsOcrManager {
    pub window: Option<Window>,
    pub lang: String,
    pub tesseract_path: String,
    pub tessdata_path: String,
}

impl WindowsOcrManager {
    pub fn new(tesseract_path: &str, tessdata_path: &str, lang: Lang) -> Self {
        WindowsOcrManager {
            window: None,
            lang: lang.to_string(),
            tesseract_path: tesseract_path.to_string(),
            tessdata_path: tessdata_path.to_string(),
        }
    }
}

fn backslash_path_to_forward(path: &str) -> String {
    if let Some(path) = path.strip_prefix("\\\\?\\") {
        return path.replace('\\', "/");
    }
    path.replace('\\', "/")
}

impl OcrManager for WindowsOcrManager {
    fn is_initialized(&self) -> bool {
        self.window.is_some()
    }

    fn set_game_window(&mut self, window_title: &str) -> Result<(), Error> {
        let windows = Window::all().unwrap();

        if let Some(game_window) = windows.iter().find(|w| {
            w.title()
                .to_lowercase()
                .contains(&window_title.to_lowercase())
        }) {
            self.window = Some(game_window.clone());
            return Ok(());
        }

        Err(Error::WindowNotFound)
    }

    fn refresh_current_window(&mut self) -> Result<(), Error> {
        if let Some(current_window) = &self.window {
            let windows = Window::all().unwrap();

            if let Some(game_window) = windows.iter().find(|w| {
                w.pid() == current_window.pid()
                    && w.id() == current_window.id()
                    && w.app_name() == current_window.app_name()
            }) {
                self.window = Some(game_window.clone());
                return Ok(());
            }
        }

        Err(Error::WindowNotFound)
    }

    fn is_window_focused(&self) -> bool {
        if let Some(window) = &self.window {
            return window.is_focused();
        }
        false
    }

    fn pause_ocr(&mut self) {
        self.window = None;
    }

    fn get_window_size(&self) -> Result<Resolution, Error> {
        if let Some(ref window) = self.window {
            return Ok(Resolution {
                width: window.width(),
                height: window.height(),
            });
        }

        Err(Error::WindowNotFound)
    }

    fn make_window_screenshot(&self, rect: Rect) -> Result<DynamicImage, Error> {
        if let Some(ref window) = self.window {
            let image = window.capture_image()?;

            let image = DynamicImage::from(image).crop(rect.x, rect.y, rect.width, rect.height);

            return Ok(image);
        }

        Err(Error::WindowNotFound)
    }

    fn perform_ocr(&self, image: Image) -> Result<String, Error> {
        let args = Args {
            executable: Some(self.tesseract_path.clone()),
            tessdata_dir: Some(backslash_path_to_forward(&self.tessdata_path)),
            lang: self.lang.clone(),
            config_variables: HashMap::from([(
                "tessedit_char_whitelist".into(),
                "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz, ():".into(),
            )]),
            dpi: Some(70),
            psm: Some(6),
            oem: Some(3),
        };

        let text = rusty_tesseract::image_to_string(&image, &args)?;

        Ok(text.trim().to_string())
    }
}
