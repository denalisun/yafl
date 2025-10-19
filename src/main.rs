use std::path::Path;

use crate::utils::path_exists;

mod utils;

fn print_usage() {
    println!("USAGE: idk implement this");
}

fn main() {
    let mut play_path: String = "".to_string();
    let mut all_tweaks: Vec<String> = Vec::new();
    let mut is_server: bool = false;

    let cobalt_path = "./assets/Cobalt.dll";
    if !path_exists(cobalt_path) {
        panic!("Could not find Cobalt DLL!");
    }

    let args: Vec<String> = std::env::args().collect();
    for i in 0..args.len() {
        let arg: String = args[i].clone();
        if arg == "-p" || arg == "--play" {
            if i+1 <= args.len() {
                play_path = args[i+1].clone().to_string();
            }
        } else if arg == "-t" || arg == "--tweak" {
            if i+1 <= args.len() {
                all_tweaks.push(args[i+1].clone().to_string());
            }
        } else if arg == "-s" || arg == "--server" {
            is_server = true;
        }
    }

    if play_path == "".to_string() {
        print_usage();
        return;
    }

    let server_path = "./assets/Reboot.dll";
    if !path_exists(server_path) && is_server {
        panic!("Could not find Reboot DLL!");
    }

    // let fn_binaries = Path::new(play_path.as_str()).join("\\FortniteGame\\Binaries\\Win64\\");
    // I gotta implement starting processes tmrw
}