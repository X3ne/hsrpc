use serde::{Deserialize, Serialize};
use std::fs;
use strsim::levenshtein;
use xcap::image::{DynamicImage, GenericImageView};

use crate::error::Error;
use crate::game::data::Data;
use crate::ocr::windows::WindowsOcrManager;
use crate::ocr::OcrManager;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Resolution {
    pub width: u32,
    pub height: u32,
}

impl Resolution {
    pub fn new(width: u32, height: u32) -> Self {
        Resolution { width, height }
    }
}

#[derive(Debug, Serialize, Deserialize, Clone, Copy)]
pub struct Rect {
    pub x: u32,
    pub y: u32,
    pub width: u32,
    pub height: u32,
}

pub fn adjust_size(original_size: i32, scale: f64) -> u32 {
    (original_size as f64 * scale).round() as u32
}

// This function is useful to mitigate OCR errors by finding
// the closest correspondence to the given text
pub fn find_closest_correspondence(text: &str, candidates: &[Data]) -> Option<Data> {
    let threshold: usize = match text.len() {
        0..=3 => 1,
        4..=6 => 2,
        _ => 5,
    };

    let mut min_distance = text.len();
    let mut closest = Data::default();

    for candidate in candidates {
        let distance = levenshtein(text, &candidate.value);
        if distance < min_distance {
            min_distance = distance;
            closest = candidate.clone();
        }
    }

    if min_distance > threshold || closest.value.is_empty() {
        return None;
    }

    Some(closest)
}

pub fn find_current_character(ocr_manager: &WindowsOcrManager, coords: &[Rect]) -> i32 {
    let mut whitest_position = -1;
    let mut whitest_value = 0;

    for (i, coord) in coords.iter().enumerate() {
        match ocr_manager.make_window_screenshot(*coord) {
            Ok(image) => {
                let (r, g, b) = get_average_color(&image);
                let brightness = r + g + b;

                if brightness > whitest_value {
                    whitest_value = brightness;
                    whitest_position = i as i32;
                }
            }
            Err(err) => {
                log::error!("Error capturing image: {}", err);
                continue;
            }
        }
    }

    whitest_position
}

pub fn get_average_color(image: &DynamicImage) -> (u32, u32, u32) {
    let mut r = 0u64;
    let mut g = 0u64;
    let mut b = 0u64;

    for y in 0..image.height() {
        for x in 0..image.width() {
            let pixel = image.get_pixel(x, y);
            r += pixel[0] as u64;
            g += pixel[1] as u64;
            b += pixel[2] as u64;
        }
    }

    let total_pixels = (image.width() * image.height()) as u64;
    (
        (r / total_pixels) as u32,
        (g / total_pixels) as u32,
        (b / total_pixels) as u32,
    )
}
