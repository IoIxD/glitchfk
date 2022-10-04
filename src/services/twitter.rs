use std::thread;
use std::time::Duration;
use crossbeam::{channel, select};

pub fn twitter_thread() {
    let interval = channel::tick(Duration::from_secs(216000));
    loop {
        select! {
            recv(interval) -> _ => {
                twitter_post()
            }
        }
    }
}

fn twitter_post() {
    println!("a");
}