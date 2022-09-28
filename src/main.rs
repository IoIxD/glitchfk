#![feature(iter_next_chunk)]


pub mod modules;

use modules::gradient;

fn main() {
    let image = gradient::new_image();
    image.save("test.png").unwrap();
}