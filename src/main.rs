#![feature(allow_internal_unstable)]
#![feature(iter_next_chunk)]

pub mod modules;
pub mod debug;

use modules::gradient;

fn main() {
    gradient::new_image().save("test.png").unwrap();
}