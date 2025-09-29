pub struct ArgOptions {
    pub main_operation: String,
    pub sub_operation: String,
    pub parameters: Vec<String>,
}

pub fn parse_args() -> ArgOptions {
    let args: Vec<String> = std::env::args().collect();
    let mut options = ArgOptions {
        main_operation: String::new(),
        sub_operation: String::new(),
        parameters: Vec::new(),
    };

    if args.len() > 2 {
        options.main_operation = args[1].clone();
        if options.main_operation == "play" || options.main_operation == "help" {
            options.parameters = args[2..].to_vec();
        } else {
            options.sub_operation = args[2].clone();
            if args.len() > 3 {
                options.parameters = args[3..].to_vec();
            }
        }
    } else {
        eprintln!("Error: Not enough arguments provided!");
        std::process::exit(1);
    }
    
    options
}