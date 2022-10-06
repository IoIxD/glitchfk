use std::{string::String, fs::{self}, error::Error};
use serde::Deserialize;

#[allow(non_snake_case)]
#[derive(Deserialize)]
pub struct Config {
    pub TwitterConsumerKey: String,
    pub TwitterConsumerSecret: String,
    pub TwitterAccessToken: String,
    pub TwitterAccessSecret: String,
    pub TwitterInterval: String,
    
    pub DiscordAuthToken: String,
    pub DiscordID: String,
    
    pub InProduction: bool
}

pub fn config() -> Result<Config, Box<dyn Error>> {
    let config_file: String = fs::read_to_string("config.toml")?.parse()?;
    let obj: Config = toml::from_str(config_file.as_str())?;
    Ok(obj)
}
