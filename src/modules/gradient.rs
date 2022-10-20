use std::{fmt::{self}, time::SystemTime};
use image::{ImageBuffer, RgbImage, Rgb};
use rand::{
    distributions::{Distribution, Standard},
    Rng,
};
use crate::debug;

type Gradient = Vec<Rgb<u8>>;

// constant values
pub const BLACK_COLOR: Rgb<u8> = Rgb::<u8>([0,0,0]);

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

// linear gradient (code from samhza)
fn new(from: Rgb<u8>, to: Rgb<u8>, n: usize) -> Vec<Rgb<u8>> {
    let mut v: Vec<Rgb<u8>> = Vec::with_capacity(n);
    for i in 0..n {
        let i = i as f32;
        let byeah = |step| {
            from.0[step] + (((to.0[step] as f32 - from.0[step] as f32) * i) / n as f32) as u8
        };
        let r = byeah(0);
        let g = byeah(1);
        let b = byeah(2);
        v.push(Rgb::<u8>([r,g,b]));
    }
    v
}

// generate a gradient object
fn new_random(width: u32, height: u32) -> Gradient {
    debug!("generating gradient ");
    let rand = || {
        let r = rand::thread_rng().gen_range(0..255);
        let g = rand::thread_rng().gen_range(0..255);
        let b = rand::thread_rng().gen_range(0..255);
        Rgb::<u8>([r,g,b])
    };

    new(rand(),rand(),((width*height)+1) as usize)
}

// generate an image from the gradient.
pub fn new_image(gradient_type: GradientType, width: u32, height: u32) -> RgbImage {
    let now = SystemTime::now();

    let grad: Gradient = new_random(width,height);
                                
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
    let mut color = BLACK_COLOR;
    for (x, y, pix) in img.enumerate_pixels_mut() {
            let position = match gradient_type {
                GradientType::Horizontal => (x*height)+offset_h,
                GradientType::Vertical => (y*height)+offset_w,
                GradientType::Radial => (y*x)+offset,
                GradientType::Diagonal => {
                    if (y+offset_h) < (x+offset_h) {
                        ((x+offset_h)-(y+offset_h))*height
                    } else {
                        x
                    }
                }
                GradientType::DiagonalBidirectional => {
                    ((x as f32 - y as f32) * height as f32).abs() as u32
                }
            };
            if position > grad.len() as u32 {
                if !warned {
                    println!("calculation for {} went out of bounds.",gradient_type);
                    warned = true;
                }
            } else {
                color = grad[position as usize];
            }

            let (r, g, b) = (color.0[0], color.0[1], color.0[2]);
            debug!("\t\t{:5}:\t\t({:3}, {:3}, {:3})\n",
            position, r, g, b);
    
            *pix = image::Rgb([r, g, b]);
        } 
    debug!("\ndone\n");
    println!("{}ms",now.elapsed().unwrap().as_millis());
    img
}

// generate a random gradient from any of the times.
pub fn random_gradient() -> RgbImage {
    println!("");
    new_image(rand::random(), 800, 600)
    
}