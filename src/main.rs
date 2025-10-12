mod utils;

use std::process::exit;

use utils::args::*;

use crate::utils::data::{get_data, YAFLData};

fn main() {
    let args: ArgOptions = parse_args();
    let mut data: YAFLData;
    match get_data() {
        Ok(dta) => {
            data = dta;
        },
        Err(_) => {
            exit(-1);
        },
    }

    match args.main_operation.as_str() {
        "profiles" => {
            match args.sub_operation.as_str() {
                "add" => {
                    if args.parameters.len() < 2 {
                        eprintln!("Wrong parameter count! {} provided, 2 or more required!", args.parameters.len());
                        exit(-1);
                    }

                    let res = data.create_instance(args.parameters[0].clone(), args.parameters[1].clone());
                    match res {
                        Ok(_) => {},
                        Err(e) => { eprintln!("{}", e) },
                    }
                },
                "remove" => {
                    if args.parameters.len() != 1 {
                        eprintln!("Wrong parameter count! {} provided, 1 required!", args.parameters.len());
                        exit(-1);
                    }

                    data.remove_instance(args.parameters[0].clone());
                },
                "list" => {
                    let mut all_instances_format: Vec<String> = Vec::new();
                    for v in data.instances {
                        all_instances_format.push(format!("\t- {} ({})", v.name, v.build_path));
                    }
                    println!("All instances: {}:\n{}", all_instances_format.len(), all_instances_format.join("\n"));
                },
                "select" => {
                    if args.parameters.len() != 1 {
                        eprintln!("Wrong parameter count! {} provided, 1 required!", args.parameters.len());
                        exit(-1);
                    }

                    let selected_instance = data.get_instance(args.parameters[0].clone());
                    match selected_instance {
                        Some(inst) => {
                            println!("Selected {}!", inst.name);
                            data.selected_instance = Some(inst.name.clone());
                        },
                        None => {},
                    }
                },
                _ => {},
            }
        },
        "play" => {
            if args.parameters.len() != 1 {
                eprintln!("Wrong parameter count! {} provided, 1 required!", args.parameters.len());
                exit(-1);
            }
        },
        _ => {},
    }
}