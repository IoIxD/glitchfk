use std::time::Duration;
use tokio::time;

pub struct TwitterConfig {
    pub consumer_key: String,
    pub consumer_secret: String,
    pub access_token: String,
    pub access_secret: String,
    pub interval: String,
}

pub async fn twitter_thread(cfg: &TwitterConfig) {
    let mut interval = time::interval(Duration::from_secs(30));

    loop {
        tokio::select! {
            _ = interval.tick() => twitter_post(cfg)
        }
    }

}

fn twitter_post(_: &TwitterConfig) {
    println!("a");
}