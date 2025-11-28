use std::path::PathBuf;

//TODO:
// -- Implement profiles

fn main() {
    let redirector_dll_path = PathBuf::from("redirector.dll");
    let redirector_dll_exists = std::fs::exists(redirector_dll_path.clone());
    match redirector_dll_exists {
        Ok(_exists) => {
            if _exists == false {
                println!("Redirector DLL not found!");
                panic!();
            }
        },

        Err(e) => {
            println!("Error! {}", e);
            panic!();
        }
    }

    let mut play_path: Option<&str> = None;
    let mut tweaks: Vec<&str> = Vec::new();
    let mut is_server: bool = false;
    let mut redirect_path: Option<&str> = None;

    let args: Vec<String> = std::env::args().collect();
    for (i, arg) in args.clone().into_iter().enumerate() {
        if i == 0 {
            continue;
        }
    
        if arg == "-p" || arg == "--play" {
            play_path = Some(args.get(i+1).unwrap().as_str());
        }

        if arg == "-t" || arg == "--tweak" {
            tweaks.push(args.get(i+1).unwrap().as_str());
        }

        if arg == "-s" || arg == "--server" {
            is_server = true;
        }

        if arg == "--redirect" {
            redirect_path = Some(args.get(i+1).unwrap().as_str());
        }
    }

    if play_path == None {
        println!("Error: game path not specified!");
        panic!();
    }
    println!("play_path: {}", play_path.unwrap());
    println!("is_server: {}", is_server);
    if redirect_path != None {
        println!("redirect_path: {}", redirect_path.unwrap());
    }

    for (i, tweak) in tweaks.into_iter().enumerate() {
        println!("Tweak {}: {}", i, tweak);
    }
}
