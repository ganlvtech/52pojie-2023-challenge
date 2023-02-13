extern crate core;

use std::f64::consts::PI;
use std::fs::File;
use std::io::BufReader;
use fft2d::slice::{ifft_2d, ifftshift};
use image::{ImageBuffer, ImageFormat, Luma};
use num_complex::Complex;

fn main() -> anyhow::Result<()> {
    let img = image::load(BufReader::new(File::open("flag12{3ac97e24}.png")?), ImageFormat::Png)?.into_luma8();
    let noise_img = image::load(BufReader::new(File::open("noise.png")?), ImageFormat::Png)?.into_luma8();
    let buffer = img.iter().zip(noise_img.iter())
        .map(|(x, y)| Complex::from_polar(*x as f64, *y as f64 / 256.0 * 2.0 * PI)).collect::<Vec<Complex<f64>>>();
    let mut buffer = ifftshift(img.width() as usize, img.height() as usize, &buffer);
    ifft_2d(img.width() as usize, img.height() as usize, &mut buffer);
    let mut re_img = ImageBuffer::new(img.width(), img.height());
    let mut im_img = ImageBuffer::new(img.width(), img.height());
    let re_buffer = buffer.iter().map(|x| x.re).collect::<Vec<f64>>();
    let im_buffer = buffer.iter().map(|x| x.im).collect::<Vec<f64>>();
    print_f64_min_max("re_buffer", &re_buffer);
    print_f64_min_max("im_buffer", &im_buffer);
    let width = img.width();
    for i in 0..img.height() {
        for j in 0..img.width() {
            let complex = buffer[(i * width + j) as usize];
            re_img.put_pixel(i, j, Luma([(complex.re / 512.0 / 2.0 + 128.0).clamp(0.0, 255.9) as u8]));
            im_img.put_pixel(i, j, Luma([(complex.im / 512.0 / 2.0 + 128.0).clamp(0.0, 255.9) as u8]));
        }
    }
    re_img.save("re_img.png")?;
    im_img.save("im_img.png")?;
    Ok(())
}

pub fn print_f64_min_max(name: &str, modulus_buffer: &[f64]) {
    let min = modulus_buffer.iter().map(|x| *x).reduce(f64::min).unwrap();
    let max = modulus_buffer.iter().map(|x| *x).reduce(f64::max).unwrap();
    println!("{} min: {} max: {}", name, min, max);
}
