use std::fmt;

use image::{ImageBuffer, RgbImage};
use rand::{
    distributions::{Distribution, Standard},
    Rng,
};
use tiny_gradient::{gradient::Gradient, RGB};
use crate::debug;
use unroll::unroll_for_loops;

// constant values
pub const BLACK_COLOR: RGB = RGB::new(0,0,0);

// gradient types
pub enum GradientType {
    Horizontal,
    Vertical,
    Diagonal,
    Radial,
    DiagonalBidirectional,
}

impl Distribution<GradientType> for Standard {
    fn sample<R: Rng + ?Sized>(&self, rng: &mut R) -> GradientType {
        match rng.gen_range(0_i8..=4_i8) {
            0 => GradientType::Horizontal,
            1 => GradientType::Vertical,
            2 => GradientType::Diagonal,
            3 => GradientType::Radial,
            4 => GradientType::DiagonalBidirectional,
            _ => GradientType::Horizontal,
        }
    }
}

impl fmt::Display for GradientType {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            GradientType::Horizontal => return write!(f, "horizontal"),
            GradientType::Vertical => return write!(f, "vertical"),
            GradientType::Diagonal => return write!(f, "diagonal"),
            GradientType::Radial => return write!(f, "radial"),
            GradientType::DiagonalBidirectional => return write!(f, "diagonal-bidirectional"),
        }
    }
}

// generate a gradient object
#[unroll_for_loops]
fn new(width: u32, height: u32) -> Gradient {
    debug!("generating gradient ");
    let rand = || {
        let r: u8 = rand::thread_rng().gen_range(0..255);
        let g = rand::thread_rng().gen_range(0..255);
        let b = rand::thread_rng().gen_range(0..255);

        RGB::new(r,g,b)
    };
    debug!("done\n");

    Gradient::new(rand(),rand(),((width*height)+1) as usize)
}

// generate an image from the gradient.
pub fn new_image(gradient_type: GradientType, width: u32, height: u32) -> RgbImage {
    // we unravel the gradient into a vector as soon as we get it,
    // because working with those is faster. 
    let grad: Vec<RGB> = new(width,height)
                                .into_iter()
                                .collect::<Vec<RGB>>();
                                
    debug!("{}",grad.len());
    // create a blank image
    let mut img = ImageBuffer::new(width, height);

    // random offset

    let offset_w = rand::thread_rng().gen_range(0..width);
    let offset_h = rand::thread_rng().gen_range(0..height);
    let offset = match rand::thread_rng().gen_range(0..1) as u32 {
        0 => offset_w,
        1_u32..=u32::MAX => offset_h, // we should only ever get 1 or above 
                                      // but the compiler says otherwise
    };

    let mut warned = false;

    // for each y and each x
    // (we can't just iterate over the gradient since again, that's too slow)
    for (x, y, pix) in img.enumerate_pixels_mut() {
            let position = match gradient_type {
                GradientType::Horizontal => (x*height)+offset_h,
                GradientType::Vertical => (y*height)+offset_w,
                GradientType::Radial => (y*x)+offset,
                GradientType::Diagonal => {
                    if ((y+offset_h) as f32) < ((x+offset_h) as f32) {
                        ((x+offset_h)-(y+offset_h))*height
                    } else {
                        x
                    }
                }
                GradientType::DiagonalBidirectional => {
                    ((x as f32 - y as f32) * height as f32).abs() as u32
                }
            };
            let mut color = grad[0];
            if position > grad.len() as u32 {
                if !warned {
                    println!("calculation for {} went out of bounds.",gradient_type);
                    warned = true;
                }
            } else {
                color = grad[position as usize];
            }

            debug!("\t\t{:5}:\t\t({:3}, {:3}, {:3})\n",
            position,
            color.r,
            color.g,
            color.b);
    
            *pix = image::Rgb([color.r, color.g, color.b]);
        } 
    debug!("\ndone\n");
    img
}

#[inline]
// generate a random gradient from any of the times.
pub fn random_gradient() -> RgbImage {
    new_image(rand::random(), 800, 600)
}