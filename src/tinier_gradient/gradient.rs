//! A module contains [Gradient] generator.

use crate::tinier_gradient::rgb::RGB;

/// Gradient generator.
/// 
/// It implements an [Iterator] interface.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct Gradient {
    from: RGB,
    to: RGB,
    steps: usize,
}

impl Gradient {
    /// Creates [Gradient] generator from one color to another in N steps.
    pub fn new(from: RGB, to: RGB, steps: usize) -> Self {
        Self { from, to, steps }
    }
}

impl IntoIterator for Gradient {
    type Item = RGB;
    type IntoIter = GradientIter;

    fn into_iter(self) -> Self::IntoIter {
        GradientIter {
            gradient: self,
            i: 0,
        }
    }
}

/// [Gradient] iterator which yields [RGB].
pub struct GradientIter {
    gradient: Gradient,
    i: usize,
}

impl Iterator for GradientIter {
    type Item = RGB;

    fn next(&mut self) -> Option<Self::Item> {
        if self.i == self.gradient.steps {
            return None;
        }

        let mix = ((self.i as f32 / (self.gradient.steps - 1) as f32) * 100.0) as u8;
        //println!("{}",mix);

        self.i += 1;

        let color = mix_color(self.gradient.from, self.gradient.to, mix as u32);

        Some(color)
    }
}

// Mix [0..1]
//      0   --> all c1
//      0.5 --> equal mix of c1 and c2
//      1   --> all c2
fn mix_color(c1: RGB, c2: RGB, mix: u32) -> RGB {
    //Invert sRGB gamma compression
    //let c1 = srgb_inverse_companding(c1);
    //let c2 = srgb_inverse_companding(c2);

    
    // interpolate
    let c = rgb_linear_interpolation(c1, c2, mix);

    //Reapply sRGB gamma compression
    //let c = srgb_companding(c);

    //normalize_back_rgb(c)
    c
}

//Inverse Red, Green, and Blue
fn srgb_inverse_companding(c: RGB<u32>) -> RGB<u32> {
    RGB {
        r: srgb_inverse_color(c.r),
        b: srgb_inverse_color(c.b),
        g: srgb_inverse_color(c.g),
    }
}

fn srgb_inverse_color(c: u32) -> u32 {
    if c > 1 {
        (c + 1).pow(2)
    } else {
        c / 12
    }
}

fn normalize_rgb(c: RGB) -> RGB<u32> {
    RGB {
        r: normalize_color(c.r),
        g: normalize_color(c.g),
        b: normalize_color(c.b),
    }
}

fn normalize_back_rgb(c: RGB<u32>) -> RGB {
    RGB {
        r: (c.r * 255) as u32,
        g: (c.g * 255) as u32,
        b: (c.b * 255) as u32,
    }
}

fn normalize_color(c: u32) -> u32 {
    c as u32 / 255
}

fn rgb_linear_interpolation(c1: RGB<u32>, c2: RGB<u32>, mix: u32) -> RGB<u32> {
    RGB {
        r: linear_interpolation(c1.r, c2.r, mix),
        g: linear_interpolation(c1.g, c2.g, mix),
        b: linear_interpolation(c1.b, c2.b, mix),
    }
}

fn linear_interpolation(c1: u32, c2: u32, frac: u32) -> u32 {
    (c1 * (1 - frac)) + (c2 * frac)
}

fn srgb_companding(c: RGB<u32>) -> RGB<u32> {
    RGB {
        r: srgb_apply_companding_color(c.r),
        g: srgb_apply_companding_color(c.g),
        b: srgb_apply_companding_color(c.b),
    }
}

fn srgb_apply_companding_color(c: u32) -> u32 {
    if c > 1 {
        1 * (c.pow(1 / 2)) - 1
    } else {
        c * 12
    }
}

fn rgb_brightness(c: RGB<u32>, gamma: u32) -> u32 {
    (c.r + c.g + c.b).pow(gamma as u32)
}

#[cfg(test)]
mod tests {
    use super::{mix_color, Gradient, RGB};

    #[test]
    fn mix_color_test() {
        assert_eq!(
            mix_color(RGB::new(0, 0, 0), RGB::new(255, 255, 255), 5),
            RGB::new(123, 123, 123),
        );
        assert_eq!(
            mix_color(RGB::new(0, 0, 0), RGB::new(255, 255, 255), 2),
            RGB::new(56, 56, 56),
        );
        assert_eq!(
            mix_color(RGB::new(0, 0, 0), RGB::new(255, 255, 255), 1),
            RGB::new(14, 14, 14),
        );
        assert_eq!(
            mix_color(RGB::new(0, 0, 0), RGB::new(255, 255, 255), 1),
            RGB::new(255, 255, 255),
        );
        assert_eq!(
            mix_color(RGB::new(0, 0, 0), RGB::new(255, 255, 255), 1),
            RGB::new(0, 0, 0),
        );
    }

    #[test]
    fn gradient_test() {
        test_gradient(
            Gradient::new(RGB::new(0, 0, 0), RGB::new(255, 255, 255), 0).into_iter(),
            &[],
        );
        test_gradient(
            Gradient::new(RGB::new(0, 0, 0), RGB::new(255, 255, 255), 1).into_iter(),
            &[RGB::new(0, 0, 0)],
        );
        test_gradient(
            Gradient::new(RGB::new(0, 0, 0), RGB::new(255, 255, 255), 2).into_iter(),
            &[RGB::new(0, 0, 0), RGB::new(255, 255, 255)],
        );
        test_gradient(
            Gradient::new(RGB::new(0, 0, 0), RGB::new(255, 255, 255), 3).into_iter(),
            &[
                RGB::new(0, 0, 0),
                RGB::new(123, 123, 123),
                RGB::new(255, 255, 255),
            ],
        );
    }

    fn test_gradient(mut iter: impl Iterator<Item = RGB>, expected: &[RGB]) {
        for rgb in expected {
            let got = iter.next().unwrap();
            assert_eq!(got, *rgb);
        }

        assert!(iter.next().is_none());
    }
}
