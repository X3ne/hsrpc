use xcap::image;

#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error("Window not found")]
    WindowNotFound,
    #[error("Failed to create app folders: {0}")]
    AppFoldersCreationFailed(String),
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
