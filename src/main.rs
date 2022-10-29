#![feature(allow_internal_unstable)]
#![feature(iter_next_chunk)]
#![feature(async_closure)]

pub mod modules;
pub mod debug;
pub mod image;
pub mod services;

// Compiling for server or anything else.
#[cfg(not(target_arch = "wasm32"))]
use services::server::server;

#[cfg(not(target_arch = "wasm32"))]
#[tokio::main]
async fn main() {
    
    tokio::select! {
        done = server::web_thread() => println!("{}",done.unwrap()),
    };
}

// Compiling for WASM.
#[cfg(all(target_arch = "wasm32", target_os = "unknown"))]
use services::wasm::wasm;

#[cfg(all(target_arch = "wasm32", target_os = "unknown"))]
fn main() {
    wasm::main();
}