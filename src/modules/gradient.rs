use image::{ImageBuffer, RgbImage};
use rand::{
    distributions::{Distribution, Standard},
    Rng,
};
use tiny_gradient::{gradient::Gradient, RGB};
use crate::debug;
use unroll::unroll_for_loops;

// constant values
pub const WIDTH: u32 = 1024;
pub const HEIGHT: u32 = 768;
pub const SIZE: usize = ((WIDTH*HEIGHT)+1) as usize;
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

// generate a gradient object
#[unroll_for_loops]
fn new() -> Gradient {
    debug!("generating gradient ");
    let mut colors: [RGB<u8>; 2] = [BLACK_COLOR,BLACK_COLOR];
    for i in 0..=1 {
        let r: u8 = rand::thread_rng().gen_range(0..255);
        let g = rand::thread_rng().gen_range(0..255);
        let b = rand::thread_rng().gen_range(0..255);

        let color = RGB::new(r,g,b);
        colors[i] = color;
    }
    debug!("done\n");
    Gradient::new(colors[0],colors[1],SIZE)
    
}

// generate an image from the gradient.
pub fn new_image(gradient_type: GradientType) -> RgbImage {
    // we unravel the gradient into an array as soon as we get it,
    // because working with those is faster. 
    let grad: [RGB<u8>; SIZE] = new()
                                .into_iter()
                                .collect::<Vec<RGB>>()
                                .try_into()
                                .unwrap();
                                
    debug!("{}",grad.len());
    // create a blank image
    let mut img = ImageBuffer::new(WIDTH, HEIGHT);
    
    // for each y and each x
    // (we can't just iterate over the gradient since again, that's too slow)
    let mut color: RGB;
    for y in 1..HEIGHT{
        for x in 1..WIDTH {
            let position: u32;
            match gradient_type {
                GradientType::Horizontal =>{
                    position = x*HEIGHT;
                }
                GradientType::Vertical => {
                    position = y*HEIGHT;
                }
                GradientType::Diagonal => {
                    if (y as f32) < (x as f32) 
                        {position = (x-y)*HEIGHT}
                    else {position = x};
                }
                GradientType::DiagonalBidirectional => {
                    position = ((x as f32 -y as f32)*HEIGHT as f32).abs() 
                        as u32;
                }
                GradientType::Radial =>{
                    position = y*x;
                }
            }
            color = grad[position as usize];

            debug!("\t\t{:5}:\t\t({:3}, {:3}, {:3})\n",
            position,
            color.r,
            color.g,
            color.b);
    
            img.put_pixel(x, y, image::Rgb([color.r,color.g,color.b]))
        } 
    }
    
    debug!("\ndone\n");
    img
}

#[inline]
// generate a random gradient from any of the times.
pub fn random_gradient() -> RgbImage {
    new_image(rand::random())
}