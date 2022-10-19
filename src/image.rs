use std::{io::BufWriter, borrow::BorrowMut};

use image::{RgbImage, Rgb};

pub fn xor_images(img1: RgbImage, img2: RgbImage) -> RgbImage {
    // todo: unravel both images arrays because those are faster
    // to work with. 

    let mut img1_iter = img1.pixels().into_iter();
    let mut img2_iter = img2.pixels().into_iter();

    RgbImage::from_fn(img1.width(), img1.height(), |_, _| {
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

pub fn png_from_u8(pixels_raw: Vec<u8>) -> Vec<u8> {
    let mut pixels = vec![0xFFu8; 0];
    
    // we wrap the following in a block so that we can
    // mutably borrow the writer and also give it back.
    {
        let p_ref: &mut Vec<u8> = pixels.borrow_mut();
        let mut w = BufWriter::new(p_ref);
        let mut encoder = png::Encoder::new(w.borrow_mut(),800,600);
        encoder.set_color(png::ColorType::Rgb);
        encoder.set_depth(png::BitDepth::Eight);
        encoder.set_trns(vec!(0xFFu8, 0xFFu8, 0xFFu8, 0xFFu8));

        let mut writer = encoder.write_header().unwrap();
        writer.write_image_data(&pixels_raw).unwrap();
    }
    
    pixels
}