use image::{RgbImage, Rgb};
use crate::modules::gradient::{WIDTH,HEIGHT};

pub fn xor_images(img1: RgbImage, img2: RgbImage) -> RgbImage {
    // todo: unravel both images arrays because those are faster
    // to work with. 

    let mut img1_iter = img1.pixels().into_iter();
    let mut img2_iter = img2.pixels().into_iter();

    RgbImage::from_fn(WIDTH, HEIGHT, |_, _| {
        let a = match img1_iter.next() {
            Some(a) => a,
            None => &Rgb([0,0,0]), 
        };
        let b = match img2_iter.next() {
            Some(b) => b,
            None => &Rgb([0,0,0]), 
        };
        Rgb([
            a.0[0] ^ b.0[0],
            a.0[1] ^ b.0[1],
            a.0[2] ^ b.0[2],
            ])
    })
}