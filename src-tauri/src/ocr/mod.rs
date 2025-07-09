use std::fmt::Display;

use core::fmt;
use rusty_tesseract::image::imageops::FilterType;
use rusty_tesseract::image::{DynamicImage, GenericImageView, ImageBuffer, Luma};
use rusty_tesseract::{image, Image};

use crate::error::Error;
use crate::utils::{Rect, Resolution};

pub mod windows;

#[derive(Debug)]
pub enum GameOcrJob {
    Character,
    Location,
    Menu,
    SubMenu,
    Boss,
    Combat,
}

impl Display for GameOcrJob {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            GameOcrJob::Character => write!(f, "character"),
            GameOcrJob::Location => write!(f, "location"),
            GameOcrJob::Menu => write!(f, "menu"),
            GameOcrJob::SubMenu => write!(f, "submenu"),
            GameOcrJob::Boss => write!(f, "boss"),
            GameOcrJob::Combat => write!(f, "combat"),
        }
    }
}

#[derive(Debug, Copy, Clone)]
pub struct PreprocessOptions {
    pub threshold: u8,
}

impl Default for PreprocessOptions {
    fn default() -> Self {
        PreprocessOptions { threshold: 135 }
    }
}

#[derive(Debug)]
pub enum Lang {
    Eng,
}

impl Display for Lang {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Lang::Eng => write!(f, "eng"),
        }
    }
}

pub trait OcrManager {
    fn is_initialized(&self) -> bool;
    fn set_game_window(&mut self, window_title: &str) -> Result<(), Error>;
    fn refresh_current_window(&mut self) -> Result<(), Error>;
    fn is_window_focused(&self) -> bool;
    fn pause_ocr(&mut self);
    fn get_window_size(&self) -> Result<Resolution, Error>;
    fn make_window_screenshot(&self, rect: Rect) -> Result<DynamicImage, Error>;

    fn preprocess_image(&self, image: DynamicImage, threshold: u8) -> DynamicImage {
        let grayscale = image.to_luma8();

        let (width, height) = image.dimensions();
        let resized =
            image::imageops::resize(&grayscale, width * 2, height * 2, FilterType::Lanczos3);

        let binarized: ImageBuffer<Luma<u8>, Vec<u8>> =
            ImageBuffer::from_fn(resized.width(), resized.height(), |x, y| {
                let pixel = resized.get_pixel(x, y);
                if pixel[0] > threshold {
                    Luma([255])
                } else {
                    Luma([0])
                }
            });

        let denoised = image::imageops::blur(&binarized, 1.0);

        DynamicImage::ImageLuma8(denoised)
    }

    fn game_ocr(
        &self,
        rect: Rect,
        job: GameOcrJob,
        preprocess_opts: Option<PreprocessOptions>,
    ) -> Result<String, Error> {
        let mut image = self.make_window_screenshot(rect)?;

        if let Some(opts) = preprocess_opts {
            image = self.preprocess_image(image, opts.threshold);
        }

        self.perform_dynamic_ocr(image)
    }

    fn perform_dynamic_ocr(&self, image: DynamicImage) -> Result<String, Error> {
        let img = Image::from_dynamic_image(&image)
            .map_err(|e| Error::Custom(format!("Failed to convert DynamicImage: {}", e)))?;

        self.perform_ocr(img)
    }

    fn perform_path_ocr(&self, path: &str) -> Result<String, Error> {
        let img = Image::from_path(path)?;

        self.perform_ocr(img)
    }

    fn perform_ocr(&self, image: Image) -> Result<String, Error>;
}
