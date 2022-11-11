#![feature(allow_internal_unstable)]
#![feature(iter_next_chunk)]
#![feature(async_closure)]

pub mod modules;
pub mod debug;
pub mod image;

use clap::{command, Command, Arg, ArgAction, arg};

fn main() {
    let mut matches = command!()
        .propagate_version(true)
        .subcommand_required(true)
        .arg_required_else_help(true)
        .subcommand(
            Command::new("images")
                .about("xor two provided images together")
                .arg(arg!(--width <VALUE>).required(false).action(ArgAction::Set))
                .arg(arg!(--height <VALUE>).required(false).action(ArgAction::Set))
                .arg(Arg::new("paths").required(true).action(ArgAction::Append))
        )
        .subcommand(
            Command::new("generate")
                .about("generate two gradients, randomly or from a prompt, and xor them together")
                .arg(arg!(--width <VALUE>).required(false).action(ArgAction::Set))
                .arg(arg!(--height <VALUE>).required(false).action(ArgAction::Set))
                .arg(arg!(--types <VALUE>).required(false).action(ArgAction::Append))
        );

    let matches_list = matches.clone().get_matches();

    match matches_list.subcommand() {
        Some(("images", sub_matches)) => {
            let paths = sub_matches
            .get_many::<String>("path")
            .unwrap_or_default()
            .map(|v| v.as_str())
            .collect::<Vec<_>>();
        }
        Some(("types", sub_matches)) => {
            let paths = sub_matches
            .get_many::<String>("paths")
            .unwrap_or_default()
            .map(|v| v.as_str())
            .collect::<Vec<_>>();

            if paths.len() <= 0 {
                println!("Usage: glitchfuck generate [type1] [type2] [type3] ...")
            } else {
                gen_image(paths)
            }
        }
        _ => {_ = matches.print_help();}
    };
}

fn gen_image(paths: Vec<&str>) {

}

fn gen_gradient(types: Vec<&str>) {

}