mod utils;
mod data;

use std::{fs, os::windows::process::CommandExt, path::{Path, PathBuf}, process::{Child, Command}};
use windows::Win32::System::Threading::{CREATE_NEW_PROCESS_GROUP, DETACHED_PROCESS};
use crate::data::YAFLProfile;

fn print_usage() {
    println!("Usage: yafl [options]\n\nOptions:");

    println!("\t--play <path> | -p <path>         Specify build to launch");
    println!("\t--tweak <path> | -t <path>        Specify a tweak to launch with, can be used multiple times");
    println!("\t--server | -s                     Enable server-mode, requires a server dll\n");
    println!("\t--redirect <url>                  Specify a URL to redirect Epic requests to");
    println!("\t--save                            Saves argument config to out_profile.json");
    println!("\t--from-profile <path>             Ignores all config, sourced from path");
}

fn safe_exit(s: Option<&str>) {
    if s.is_some() {
        println!("{}", s.unwrap());
    }
    std::process::exit(0);
}

fn main() {
    let redirector_dll_path = PathBuf::from("redirector.dll");
    let redirector_dll_exists = std::fs::exists(redirector_dll_path.clone());
    match redirector_dll_exists {
        Ok(_exists) => {
            if _exists == false {
                safe_exit(Some("ERROR: Redirector DLL not found!"));
                print_usage();
            }
        },

        Err(e) => {
            panic!("ERROR: {}", e);
        }
    }
    
    let args: Vec<String> = std::env::args().collect();

    let mut play_path: Option<String> = None;
    let mut tweaks: Vec<String> = Vec::new();
    let mut is_server: bool = false;
    let mut redirect_path: Option<String> = None;
    let mut should_save: bool = false;
    let mut profile_path: Option<String> = None;

    if args.contains(&"--help".to_string()) {
        print_usage();
        std::process::exit(0);
    }

    for (i, arg) in args.clone().into_iter().enumerate() {    
        if arg == "-p" || arg == "--play" {
            play_path = Some((*args.get(i+1).unwrap().clone()).to_string());
        }

        if arg == "-t" || arg == "--tweak" {
            tweaks.push((*args.get(i+1).unwrap().clone()).to_string());
        }

        if arg == "-s" || arg == "--server" {
            is_server = true;
        }

        if arg == "--redirect" {
            redirect_path = Some(args.get(i+1).unwrap().to_string());
        }

        if arg == "--save" {
            should_save = true;
        }

        if arg == "--from-profile" {
            profile_path = Some(args.get(i+1).unwrap().to_string());
        }
    }

    if profile_path.is_some() {
        if should_save {
            print_usage();
            safe_exit(Some("Cannot use --from-profile and --save simultaneously!"));
        }

        let contents = match fs::read_to_string(profile_path.unwrap()) {
            Ok(c) => {
                c
            },
            Err(_) => {
                panic!("ERROR: Failed to read from profile!");
            },
        };
        
        let deserialized: YAFLProfile = serde_json::from_str(&contents).unwrap();

        play_path = Some(deserialized.play_path);
        redirect_path = deserialized.redirect_path;
        is_server = deserialized.is_server;
        tweaks = deserialized.tweaks;
    }

    if play_path == None {
        print_usage();
        safe_exit(None);
    }

    let server_dll = PathBuf::from("server.dll");
    if is_server {
        let server_dll_exists = std::fs::exists(server_dll.clone());
        match server_dll_exists {
            Ok(_exists) => {
                if _exists == false {
                    print_usage();
                    safe_exit(Some("ERROR: Server DLL not found!"));
                }
            },

            Err(e) => {
                panic!("ERROR: {}", e.to_string());
            }
        }
    }

    if redirect_path.is_some() {
        let redirect_path_buf = PathBuf::from((std::env::var("LOCALAPPDATA").unwrap().to_string() + "\\.yaflredirect").as_str());
        let res = std::fs::write(redirect_path_buf, redirect_path.clone().unwrap());
        match res {
            Ok(_) => {},
            Err(_) => {
                panic!("Failed to write redirect path!");
            },
        }
    }

    if should_save {
        let profile: YAFLProfile = YAFLProfile::new(
            play_path.unwrap(),
            is_server,
            redirect_path.clone(),
            tweaks
        );
        
        let serialized = serde_json::to_string(&profile).unwrap();

        match std::fs::write("out_profile.json", serialized) {
            Ok(_) => {
                println!("Successfully wrote to out_profile.json! It is HIGHLY recommended you rename this file!");
            },
            Err(_) => {
                panic!("ERROR: Failed to write to out_profile.json!");
            },
        }

        std::process::exit(0);
    }

    let fortnite_binaries = Path::new(play_path.unwrap().as_str()).join("FortniteGame\\Binaries\\Win64");
    let fortnite_launcher_path = fortnite_binaries.clone().as_path().join("FortniteLauncher.exe");
    let fortnite_eac_path = fortnite_binaries.clone().as_path().join("FortniteClient-Win64-Shipping_EAC.exe");
    let fortnite_client_path = fortnite_binaries.clone().as_path().join("FortniteClient-Win64-Shipping.exe");

    let creation_flags: u32 = (DETACHED_PROCESS | CREATE_NEW_PROCESS_GROUP).0;

    let launch_args: Vec<&str> = "-epicapp=Fortnite -epicenv=Prod -epiclocale=en-us -epicportal -skippatchcheck -NOSSLPINNING -nobe -fromfl=eac -fltoken=3db3ba5dcbd2e16703f3978d -caldera=eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiYmU5ZGE1YzJmYmVhNDQwN2IyZjQwZWJhYWQ4NTlhZDQiLCJnZW5lcmF0ZWQiOjE2Mzg3MTcyNzgsImNhbGRlcmFHdWlkIjoiMzgxMGI4NjMtMmE2NS00NDU3LTliNTgtNGRhYjNiNDgyYTg2IiwiYWNQcm92aWRlciI6IkVhc3lBbnRpQ2hlYXQiLCJub3RlcyI6IiIsImZhbGxiYWNrIjpmYWxzZX0.VAWQB67RTxhiWOxx7DBjnzDnXyyEnX7OljJm-j2d88G_WgwQ9wrE6lwMEHZHjBd1ISJdUO1UVUqkfLdU5nofBQ".split_whitespace().collect();

    let launcher_process: Option<Child> = match Command::new(fortnite_launcher_path.to_str().unwrap())
        .current_dir(&fortnite_binaries)
        .creation_flags(creation_flags)
        .args(&launch_args)
        .spawn() {
            Ok(c) => {
                utils::nt_suspend_process(c.id());
                Some(c)
            },
            Err(_) => {
                None
            }
        };

    let eac_process: Option<Child> = match Command::new(fortnite_eac_path.to_str().unwrap())
        .current_dir(&fortnite_binaries)
        .creation_flags(creation_flags)
        .args(&launch_args)
        .spawn() {
            Ok(c) => {
                utils::nt_suspend_process(c.id());
                Some(c)
            },
            Err(_) => {
                None
            }
        };

    let mut fortnite_process = match Command::new(fortnite_client_path)
        .current_dir(&fortnite_binaries)
        .creation_flags(creation_flags)
        .args(&launch_args)
        .spawn() {
            Ok(res) => {
                res
            },
            Err(_) => {
                if launcher_process.is_some() {
                    let _ = launcher_process.unwrap().kill();
                }
                if eac_process.is_some() {
                    let _ = eac_process.unwrap().kill();
                }
                panic!("ERROR: Failed to launch FortniteClient-Win64-Shipping!");
            }
        };

    utils::inject_dll(fortnite_process.id(), PathBuf::from(
        std::env::current_dir().unwrap()
    ).join(redirector_dll_path.to_str().unwrap()).to_str().unwrap());

    let _ = fortnite_process.wait();

    if launcher_process.is_some() {
        let _ = launcher_process.unwrap().kill();
    }
    if eac_process.is_some() {
        let _ = eac_process.unwrap().kill();
    }
}
