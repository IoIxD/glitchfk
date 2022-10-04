#![feature(allow_internal_unstable)]
#![feature(iter_next_chunk)]

pub mod modules;
pub mod debug;
pub mod image;
pub mod services;

use modules::gradient;
use services::{twitter,discord};

use std::thread;

fn main() {
    let grad1 = gradient::random_gradient();
    let grad2 = gradient::random_gradient();

    let final_grad = image::xor_images(grad1, grad2);

    //thread::spawn(twitter::twitter_thread);
    twitter::twitter_thread();
}