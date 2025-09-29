use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
pub struct YAFLInstance {
    pub name: String,
    pub mods_path: String,
    pub build_path: String,
    pub additional_args: Vec<String>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct YAFLData {
    pub instances: Vec<YAFLInstance>,
    pub selected_instance: Option<String>,
}

impl YAFLData {
    pub fn new() -> Self {
        YAFLData {
            instances: Vec::new(),
            selected_instance: None,
        }
    }
}

pub fn init_data() -> YAFLData {
    let data: YAFLData = YAFLData::new();
    let json_str = serde_json::to_string(&data).unwrap();
    std::fs::write("yafl_data.json", json_str).expect("Unable to write to yafl_data.json!");

    data
}

pub fn get_data() -> YAFLData {
    let mut data: YAFLData = init_data();
    let data_path_str: &str = "yafl_data.json";
    let data_path = std::path::Path::new(data_path_str);
    if data_path.exists() {
        let json_str = std::fs::read_to_string(data_path_str).expect("Unable to read yafl_data.json!");
        data = serde_json::from_str(&json_str).expect("JSON was not well-formatted");
    }
    data
}

pub fn create_instance(data: &mut YAFLData, name: String, build_path: String) -> Result<(), String> {
    for instance in &data.instances {
        if instance.name == name {
            return Err(format!("Instance with name '{}' already exists!", name));
        }
    }

    let mods_path = format!("{}/mods", build_path);
    let new_instance = YAFLInstance {
        name: name.clone(),
        mods_path,
        build_path,
        additional_args: Vec::new(),
    };
    data.instances.push(new_instance);
    data.selected_instance = Some(name);

    let json_str = serde_json::to_string(data).unwrap();
    std::fs::write("yafl_data.json", json_str).expect("Unable to write to yafl_data.json!");

    Ok(())
}