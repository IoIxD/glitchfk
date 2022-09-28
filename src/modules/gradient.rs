use image::{ImageBuffer, RgbImage};
use rand::Rng;
use tiny_gradient::{gradient::Gradient, RGB};

use crate::debug;

const WIDTH: u32 = 640;
const HEIGHT: u32 = 480;
const SIZE: usize = (WIDTH*HEIGHT) as usize;
const BLACK_COLOR: RGB = RGB::new(0,0,0);

// generate a gradient object
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
pub fn new_image() -> RgbImage {
    // we unravel the gradient into an array as soon as we get it,
    // because working with them is faster. 
    let grad: [RGB<u8>; SIZE] = new()
                                .into_iter()
                                .collect::<Vec<RGB>>()
                                .try_into()
                                .unwrap();
                                
    // create a blank image
    let mut img = ImageBuffer::new(WIDTH, HEIGHT);
    
    // for each y and each x
    // (we can't just iterate over the gradient since again, that's too slow)
    for y in 0..HEIGHT{
        for x in 0..WIDTH {
            let color = grad[(y*x) as usize];

            debug!("\r\t\t{:5}, {:5}:\t\t({:3}, {:3}, {:3})",
            x,y,
            color.r,
            color.g,
            color.b);
    
            img.put_pixel(x, y, image::Rgb([color.r,color.g,color.b]))
        } 
    }
    
    debug!("\ndone\n");
    img
}
