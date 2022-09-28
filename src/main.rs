#![feature(iter_next_chunk)]


pub mod modules;

use modules::gradient;

fn main() {
    gradient::new_image().save("test.png").unwrap();
}