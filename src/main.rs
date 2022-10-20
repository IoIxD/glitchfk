#![feature(allow_internal_unstable)]
#![feature(iter_next_chunk)]
#![feature(async_closure)]

pub mod modules;
pub mod debug;
pub mod image;
pub mod services;

//use modules::gradient;
use services::{web};

#[tokio::main]
async fn main() {
    tokio::select! {
        done = web::web_thread() => println!("{}",done.unwrap()),
    };
}