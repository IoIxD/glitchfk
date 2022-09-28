use image::{ImageBuffer, RgbaImage};
use rand::Rng;
use tiny_gradient::{gradient::Gradient, RGB};

const WIDTH: u32 = 1024;
const HEIGHT: u32 = 768;
const BLACK_COLOR: RGB = RGB::new(0,0,0);

fn new() -> Gradient {
    let mut colors: [RGB<u8>; 2] = [BLACK_COLOR,BLACK_COLOR];
    for i in 0..=1 {
        let r: u8 = rand::thread_rng().gen_range(0..255);
        let g = rand::thread_rng().gen_range(0..255);
        let b = rand::thread_rng().gen_range(0..255);

        let color = RGB::new(r,g,b);
        colors[i] = color;
    }
    println!("generated gradient");
    Gradient::new(colors[0],colors[1],(WIDTH*HEIGHT) as usize)
    
}

pub fn new_image() -> RgbaImage {
    let mut grad = new().into_iter();

    let img: RgbaImage = ImageBuffer::from_fn(WIDTH,HEIGHT, 
                                            |x, y| {

        let color_ = grad.next();
        if color_ == None {
            return image::Rgba([0, 0, 0, 255]);
        }
        let color = color_.unwrap();

        print!("\r\t\t{:5}, {:5}:\t\t({:3}, {:3}, {:3})",
        x, y,
        color.r,
        color.g,
        color.b);

        image::Rgba([color.r,color.g,color.b,255])
    });
    println!("\nsaved gradient");
    img
}