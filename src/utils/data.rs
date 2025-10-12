use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone)]
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
    
    pub fn create_instance(&mut self, name: String, build_path: String) -> Result<(), String> {
        for instance in &self.instances {
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
        self.instances.push(new_instance);
        self.selected_instance = Some(name);

        let json_str = serde_json::to_string(&self).unwrap();
        std::fs::write("yafl_data.json", json_str).expect("Unable to write to yafl_data.json!");

        Ok(())
    }

    pub fn remove_instance(&mut self, name: String) {
        let mut l: i32 = -1;
        for i in 0..(&self.instances).len() {
            if self.instances[i].name == name {
                l = i as i32;
            }
        }
        if l != -1 { self.instances.remove(l.try_into().unwrap()); };
    }

    pub fn get_instance(&mut self, name: String) -> Option<&YAFLInstance> {
        for i in &self.instances {
            if i.name == name {
                return Some(i);
            }
        }
        None
    }

    pub fn save_current_data(&self) {
        let json_str = serde_json::to_string(self).unwrap();
        std::fs::write("yafl_data.json", json_str).expect("Unable to write to yafl_data.json!");
    }
}

pub fn init_data() -> YAFLData {
    let data: YAFLData = YAFLData::new();
    let json_str = serde_json::to_string(&data).unwrap();
    std::fs::write("yafl_data.json", json_str).expect("Unable to write to yafl_data.json!");

    data
}

pub fn get_data() -> Result<YAFLData, String> {
    let mut data: YAFLData = init_data();
    let data_path_str: &str = "./yafl_data.json";
    let data_path = std::path::Path::new(data_path_str);
    if data_path.exists() {
        let json_str = std::fs::read_to_string(data_path_str).expect("Unable to read yafl_data.json!");
        let res: Result<YAFLData, serde_json::Error> = serde_json::from_str(&json_str);
        match res {
            Ok(v) => {
                data = v;
            },
            Err(e) => {
                eprintln!("Failed to read JSON!");
                return Err(e.to_string());
            },
        }
    }
    Ok(data)
}