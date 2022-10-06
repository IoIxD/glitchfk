use std::primitive::str;
use std::time::SystemTime;
use std::fmt;

use crate::modules::gradient;
use crate::image::xor_images;
use tokio::fs::File;
use image::RgbImage;
use serenity::async_trait;

use serenity::builder::{CreateApplicationCommand};
use serenity::model::application::command::{Command,CommandOptionType};
use serenity::model::application::interaction::{Interaction, InteractionResponseType};
use serenity::model::gateway::Ready;
use serenity::model::prelude::{AttachmentType};
use serenity::model::prelude::interaction::application_command::CommandDataOption;
use serenity::prelude::*;

struct Handler;

#[async_trait]
impl EventHandler for Handler {
    async fn ready(&self, ctx: Context, ready: Ready) {
        println!("{} is connected!", ready.user.name);

        Command::create_global_application_command(&ctx.http, |command| {
            discord_register(command)
        })
        .await
        .expect("couldn't register global command.");
    }

    async fn interaction_create(&self, ctx: Context, interaction: Interaction) {
        if let Interaction::ApplicationCommand(command) = interaction {
            let (content, file) = match command.data.name.as_str() {
                "glitchfuck" => discord_command(&command.data.options).await,
                _ => panic!("we should never reach this point.")
            };

            if let Err(why) = command
                .create_interaction_response(&ctx.http, |response| {
                    response
                        .kind(InteractionResponseType::ChannelMessageWithSource)
                        .interaction_response_data(|m| {
                            m.content(content);
                            m.add_file(AttachmentType::from((&file, "glitchfuck.png")));
                            m
                        })
                })
                .await
            {
                println!("Cannot respond to slash command: {}", why);
            }
        }
    }
}

pub struct DiscordConfig {
    pub auth_token: String,
    pub id: String,
}

pub async fn discord_thread(cfg: &DiscordConfig) -> String {
    let mut client = Client::builder(&cfg.auth_token, GatewayIntents::empty())
        .event_handler(Handler)
        .await
        .expect("Error creating client");
    
    match client.start().await {
        Ok(..) => "Discord thread exited gracefully".to_string(),
        Err(err) => format!("Discord thread exited with error: {:?}", err),
    }
}

async fn discord_command(options: &[CommandDataOption]) -> (String, tokio::fs::File) {
    let response = &mut String::new();

    /* options */
    let types_str: &str = match options.get(0) {
        Some(a) => match &a.value {
            Some(b) => b.as_str().unwrap(),
            None => "horizontal,vertical",
        },
        None => "horizontal,vertical",
    };
    let forcelowcontrast = match options.get(2) {
        Some(a) => match &a.value {
            Some(b) => b.as_bool().unwrap(),
            None => false,
        },
        None => false,
    };
    if forcelowcontrast {
        fmt::write(response, format_args!("WARNING: forcelowcontrast ignored, it's not implemented as noise images aren't implemented yet.\n")).expect("idk");
    }
    let width = match options.get(3) {
        Some(a) => match &a.value {
            Some(b) => b.as_u64().unwrap() as u32,
            None => 800,
        },
        None => 800,
    };
    let height = match options.get(4) {
        Some(a) => match &a.value {
            Some(b) => b.as_u64().unwrap() as u32,
            None => 600,
        },
        None => 600,
    };

    /* create a vector of user provided types. */
    let mut types: Vec<gradient::GradientType> = Vec::new();
    
    types_str.split(",").for_each(|t| {
        match t {
            "horizontal" => types.push(gradient::GradientType::Horizontal),
            "vertical" => types.push(gradient::GradientType::Vertical),
            "diagonal" => types.push(gradient::GradientType::Diagonal),
            "radial" => types.push(gradient::GradientType::Radial),
            "diagonal_bidirectional" => types.push(gradient::GradientType::DiagonalBidirectional),
            _ => {
                fmt::write(response, format_args!("WARNING: Type '{}' ignored, using horizontal in its place.\n",t)).expect("idk");
                types.push(gradient::GradientType::Horizontal);
            },
        }
    });

    let mut grads: Vec<RgbImage> = Vec::new();

    for grad_type in types {
        grads.push(gradient::new_image(grad_type,width,height))
    }

    let mut final_grad: RgbImage = match grads.get(0) {
        Some(a) => a.clone(),
        None => {
            panic!("no/invalid image");
        }
    };
    for i in 1..grads.len() {
        let next_image: RgbImage = match grads.get(i) {
            Some(a) => a.clone(),
            None => {
                panic!("no/invalid image");
            }
        };
        final_grad = xor_images(final_grad, next_image);
    }    

    let filename: &String = &format!("req_{:?}.png",SystemTime::now());
    final_grad.save(filename).expect("couldn't save file idk");
    let file: File = match File::open(filename).await {
        Ok(a) => a,
        Err(err) => panic!("{}",err),
    };

    ("".to_string(), file)
}

struct Option<'a> {
    name: &'a str,
    description: &'a str,
    kind: CommandOptionType,
    required: bool
}

fn discord_register(command: &mut CreateApplicationCommand) -> &mut CreateApplicationCommand {
    let options: [Option; 4] = [
        Option{
            name: "types",
            description: "The types of images to generate, seperated by commas. Random image by default. No docs yet.",
            kind: CommandOptionType::String,
            required: false,
        },
        Option{
            name: "forcelowcontrast",
            description: "Don't return image if average contrast is high. Used for Twitter bot.",
            kind: CommandOptionType::Boolean,
            required: false,
        },
        Option{
            name: "width",
            description: "Width of the image. Default is 1024.",
            kind: CommandOptionType::Integer,
            required: false,
        },
        Option{
            name: "height",
            description: "Height of the image. Default is 768.",
            kind: CommandOptionType::Integer,
            required: false,
        }
    ];
    let cmd = command.name("glitchfuck").description("Runs glitchfuck");
    
    for opt in options {
        cmd.create_option( |option| {
            option
                .name(opt.name)
                .description(opt.description)
                .kind(opt.kind)
                .required(opt.required)
        });
    }
    cmd
}

