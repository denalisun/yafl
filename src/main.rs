mod utils;

use std::{os::windows::process::CommandExt, path::{Path, PathBuf}, process::Command};

use windows::Win32::System::Threading::{CREATE_NEW_PROCESS_GROUP, DETACHED_PROCESS};

//TODO:
// -- Implement profiles

fn main() {
    let redirector_dll_path = PathBuf::from("redirector.dll");
    let redirector_dll_exists = std::fs::exists(redirector_dll_path.clone());
    match redirector_dll_exists {
        Ok(_exists) => {
            if _exists == false {
                panic!("Redirector DLL not found!");
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
        panic!("Game path not specified!");
    }

    let server_dll = PathBuf::from("server.dll");
    if is_server {
        let server_dll_exists = std::fs::exists(server_dll.clone());
        match server_dll_exists {
            Ok(_exists) => {
                if _exists == false {
                    panic!("Server DLL not found!");
                }
            },

            Err(e) => {
                println!("Error! {}", e);
                panic!();
            }
        }
    }

    if redirect_path.is_some() {
        let redirect_path_buf = PathBuf::from((env!("LOCALAPPDATA").to_string() + "\\.yaflredirect").as_str());
        let res = std::fs::write(redirect_path_buf, redirect_path.unwrap());
        match res {
            Ok(_) => {

            },
            Err(_) => {

            },
        }
    }

    let fortnite_binaries = Path::new(play_path.unwrap()).join("FortniteGame\\Binaries\\Win64");
    let fortnite_launcher_path = fortnite_binaries.clone().as_path().join("FortniteLauncher.exe");
    let fortnite_eac_path = fortnite_binaries.clone().as_path().join("FortniteClient-Win64-Shipping_EAC.exe");
    let fortnite_client_path = fortnite_binaries.clone().as_path().join("FortniteClient-Win64-Shipping.exe");

    let creation_flags: u32 = (DETACHED_PROCESS | CREATE_NEW_PROCESS_GROUP).0;

    let launch_args: Vec<&str> = "-epicapp=Fortnite -epicenv=Prod -epiclocale=en-us -epicportal -skippatchcheck -NOSSLPINNING -nobe -fromfl=eac -fltoken=3db3ba5dcbd2e16703f3978d -caldera=eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiYmU5ZGE1YzJmYmVhNDQwN2IyZjQwZWJhYWQ4NTlhZDQiLCJnZW5lcmF0ZWQiOjE2Mzg3MTcyNzgsImNhbGRlcmFHdWlkIjoiMzgxMGI4NjMtMmE2NS00NDU3LTliNTgtNGRhYjNiNDgyYTg2IiwiYWNQcm92aWRlciI6IkVhc3lBbnRpQ2hlYXQiLCJub3RlcyI6IiIsImZhbGxiYWNrIjpmYWxzZX0.VAWQB67RTxhiWOxx7DBjnzDnXyyEnX7OljJm-j2d88G_WgwQ9wrE6lwMEHZHjBd1ISJdUO1UVUqkfLdU5nofBQ".split_whitespace().collect();

    let mut launcher_process = Command::new(fortnite_launcher_path)
        .current_dir(&fortnite_binaries)
        .creation_flags(creation_flags)
        .args(&launch_args)
        .spawn()
        .expect("Failed to launch FortniteLauncher!");
    utils::nt_suspend_process(launcher_process.id());

    let mut eac_process = match Command::new(fortnite_eac_path.to_str().unwrap())
        .current_dir(&fortnite_binaries)
        .creation_flags(creation_flags)
        .args(&launch_args)
        .spawn() {
            Ok(result) => {
                utils::nt_suspend_process(result.id());
                result
            },
            Err(_) => {
                let _ = launcher_process.kill();
                panic!("Failed to start FortniteClient-Win64-Shipping_EAC!");
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
                let _ = launcher_process.kill();
                let _ = eac_process.kill();
                panic!("Failed to launch FortniteClient-Win64-Shipping!");
            }
        };

    utils::inject_dll(fortnite_process.id(), PathBuf::from(
        std::env::current_dir().unwrap()
    ).join(redirector_dll_path.to_str().unwrap()).to_str().unwrap());

    let _ = fortnite_process.wait();

    let _ = launcher_process.kill();
    let _ = eac_process.kill();
}