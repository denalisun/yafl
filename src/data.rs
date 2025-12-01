use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Debug)]
pub struct YAFLProfile {
    pub play_path: String,
    pub is_server: bool,
    pub redirect_path: Option<String>,
    pub tweaks: Vec<String>,
}

impl YAFLProfile {
    pub fn new(play_path: String, is_server: bool, redirect_path: Option<String>, tweaks: Vec<String>) -> YAFLProfile {
        let mut new_vec: Vec<String> = Vec::new();
        for tweak in tweaks {
            new_vec.push(tweak.to_string());
        }
        YAFLProfile { play_path: play_path.to_string(), is_server, redirect_path: redirect_path, tweaks: new_vec }
    }
}