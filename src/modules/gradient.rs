use image::{ImageBuffer, RgbaImage};
use rand::Rng;
use colorgrad::Color;

const WIDTH: u32 = 1024;
const HEIGHT: u32 = 768;

fn new() -> colorgrad::Gradient {
    let color_num = rand::thread_rng().gen_range(0..32768);
    let mut colors = Vec::new();
    for _ in 1..color_num {
        let r = rand::thread_rng().gen_range(0..255);
        let g = rand::thread_rng().gen_range(0..255);
        let b = rand::thread_rng().gen_range(0..255);

        let color = Color::from_rgba8(r,g,b,255);

        colors.push(color)
    }

    let g = 
        match colorgrad::CustomGradient::new()
        .colors(colors.as_mut_slice())
        .build() {
            Ok(g) => g,
            Err(err) => {
                panic!("{err}");
                
            }
        };
    g
}

pub fn new_image() -> RgbaImage {
    let grad = new();
    let img: RgbaImage = ImageBuffer::from_fn(WIDTH,HEIGHT, |x, y| {
        let color = grad.at((x) as f64 / (WIDTH) as f64);

        println!("\r{}, {}: ({}, {}, {})",
        x, y,
        color.r*255.0,
        color.g*255.0,
        color.b*255.0);

        image::Rgba([
            (color.r * 255.0) as u8, 
            (color.g * 255.0) as u8, 
            (color.b * 255.0) as u8,
            255
        ])
    });
    img
}