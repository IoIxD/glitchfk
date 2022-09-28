use image::{ImageBuffer, RgbaImage};
use rand::Rng;
use tiny_gradient::{gradient::Gradient, RGB};

const WIDTH: u32 = 640;
const HEIGHT: u32 = 480;
const SIZE: usize = (WIDTH*HEIGHT) as usize;
const BLACK_COLOR: RGB = RGB::new(0,0,0);

fn new() -> Gradient {
    print!("generating gradient ");
    let mut colors: [RGB<u8>; 2] = [BLACK_COLOR,BLACK_COLOR];
    for i in 0..=1 {
        let r: u8 = rand::thread_rng().gen_range(0..255);
        let g = rand::thread_rng().gen_range(0..255);
        let b = rand::thread_rng().gen_range(0..255);

        let color = RGB::new(r,g,b);
        colors[i] = color;
    }
    print!("done\n");
    Gradient::new(colors[0],colors[1],SIZE)
    
}

pub fn new_image() -> RgbaImage {
    let grad_ = new();
    print!("unwrapping gradient into array ");
    let grad: [RGB<u8>; SIZE] = grad_.into_iter()
    .collect::<Vec<RGB>>()
    .try_into()
    .unwrap();
    print!("done\n");


    print!("bringing gradient to image\n");
    let mut img = ImageBuffer::new(WIDTH, HEIGHT);
    for y in 0..HEIGHT{
        for x in 0..WIDTH {
            let color = grad[(y*x) as usize];

            print!("\r\t\t{:5}, {:5}:\t\t({:3}, {:3}, {:3})",
            x,y,
            color.r,
            color.g,
            color.b);
    
            img.put_pixel(x, y, image::Rgba([color.r,color.g,color.b,255]))
        } 
    }
    print!("\ndone\n");
    img
}
