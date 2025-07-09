use std::collections::BTreeMap;
use std::fs::File;
use std::path::{Path, PathBuf};
use std::{error::Error, io::Read};

use csv::ReaderBuilder;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone, Eq, PartialEq, Default)]
pub struct Data {
    pub asset_id: String,
    pub value: String,
    #[serde(default)]
    pub message: String,
    #[serde(default)]
    pub region: String,
    #[serde(default)]
    pub sub_region: String,
}

#[derive(Debug)]
pub struct GameData {
    pub characters: BTreeMap<String, Data>,
    pub locations: BTreeMap<String, Data>,
    pub menus: BTreeMap<String, Data>,
    pub sub_menus: BTreeMap<String, Data>,
    pub bosses: BTreeMap<String, Data>,
}

impl GameData {
    fn load_data<P: AsRef<Path>>(file_path: P) -> Result<BTreeMap<String, Data>, Box<dyn Error>> {
        let mut file = File::open(file_path)?;
        let mut contents = String::new();
        file.read_to_string(&mut contents)?;

        let mut reader = ReaderBuilder::new().from_reader(contents.as_bytes());
        let mut data = BTreeMap::new();

        for result in reader.deserialize() {
            let record: Data = result?;
            data.insert(record.value.clone(), record);
        }

        Ok(data)
    }

    pub fn load(resource_dir: &PathBuf) -> Result<GameData, Box<dyn Error>> {
        let characters_path = resource_dir.join("characters.csv");
        let locations_path = resource_dir.join("locations.csv");
        let menus_path = resource_dir.join("menus.csv");
        let sub_menus_path = resource_dir.join("sub_menus.csv");
        let bosses_path = resource_dir.join("bosses.csv");

        let characters = Self::load_data(characters_path)?;
        let locations = Self::load_data(locations_path)?;
        let menus = Self::load_data(menus_path)?;
        let sub_menus = Self::load_data(sub_menus_path)?;
        let bosses = Self::load_data(bosses_path)?;

        Ok(GameData {
            characters,
            locations,
            menus,
            sub_menus,
            bosses,
        })
    }
}
