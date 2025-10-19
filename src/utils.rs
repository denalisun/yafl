use std::fs;

pub fn path_exists(path: &str) -> bool {
    let metadata = fs::metadata(path);
    match metadata {
        Ok(_) => {
            true
        },
        Err(_) => {
            false
        },
    }
}