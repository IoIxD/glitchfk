#![feature(allow_internal_unstable)]
#![feature(iter_next_chunk)]
#![feature(async_closure)]

pub mod modules;
pub mod debug;
pub mod image;
pub mod services;
pub mod array_vec_generics;
pub mod tinier_gradient;

//use modules::gradient;
use services::{web};

#[tokio::main]
async fn main() {
    tokio::select! {
        done = web::web_thread() => println!("{}",done.unwrap()),
    };
}