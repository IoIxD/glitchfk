#![feature(allow_internal_unstable)]
#![feature(iter_next_chunk)]
#![feature(async_closure)]

pub mod modules;
pub mod debug;
pub mod image;
pub mod services;
pub mod config;

use config::{config, Config};
//use modules::gradient;
use services::{twitter::{self, TwitterConfig},discord};

#[tokio::main]
async fn main() {
    let cfg: Config = match config() {
        Ok(x) => x,
        Err(err) => panic!("{}", err),
    };

    /*let grad1 = gradient::random_gradient();
    let grad2 = gradient::random_gradient();

    let final_grad = image::xor_images(grad1, grad2);*/

    let twitcfg = &TwitterConfig{
        consumer_key: cfg.TwitterConsumerKey,
        consumer_secret: cfg.TwitterConsumerSecret,
        access_token: cfg.TwitterAccessToken,
        access_secret: cfg.TwitterAccessSecret,
        interval: cfg.TwitterInterval,
    };

    let discordcfg = &discord::DiscordConfig{
        auth_token: cfg.DiscordAuthToken,
        id: cfg.DiscordID,
    };

    tokio::select! {
        _ = twitter::twitter_thread(twitcfg) => {
            println!("twitter thread finished");
        }
        _ = discord::discord_thread(discordcfg) => {
            println!("discord thread finished");
        }
    }
}