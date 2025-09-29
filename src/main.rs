mod utils;

use utils::args::*;

use crate::utils::data::{get_data, YAFLData};

fn main() {
    let args: ArgOptions = parse_args();
    println!("Main Operation: {}", args.main_operation);
    println!("Sub Operation: {}", args.sub_operation);
    println!("Parameters: {:?}", args.parameters);

    let data: YAFLData = get_data();
    dbg!(data);
}