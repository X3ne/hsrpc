use tauri::AppHandle;
use xcap::image;

#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error("Window not found")]
    WindowNotFound,
    #[error("App folders creation failed: {0}")]
    AppFolderCreation(String),
    #[error("Failed to resolve app path: {0}")]
    AppPathResolution(String),
    #[error("Config error: {0}")]
    ConfigError(String),
    #[error("Screenshot error: {0}")]
    ScreenshotError(#[from] xcap::XCapError),
    #[error("Image error: {0}")]
    ImageError(#[from] image::ImageError),
    #[error("IO error: {0}")]
    IoError(#[from] std::io::Error),
    #[error("Parse error: {0}")]
    ParseError(#[from] std::num::ParseIntError),
    #[error("Csv error: {0}")]
    CsvError(#[from] csv::Error),
    #[error("Tesseract error: {0}")]
    TesseractError(#[from] rusty_tesseract::TessError),
    #[error("Custom error: {0}")]
    Custom(String),
}

impl Error {
    pub fn report_and_exit(self, app_handle: &AppHandle, context_message: &str) {
        let error_message = format!("{}: {}", context_message, self);
        log::error!("{}", error_message);

        app_handle.exit(1)
    }
}
