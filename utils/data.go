package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

type YAFLInstance struct {
	Name           string
	ModsPath       string
	BuildPath      string
	AdditionalArgs string
}

func initData() ([]YAFLInstance, error) {
	data := []YAFLInstance{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Could not marshal JSON: %s\n", err)
		return nil, err
	}
	yaflPath := filepath.Join(os.Getenv("LOCALAPPDATA"), ".yafl")
	if _, err = os.Stat(yaflPath); errors.Is(err, fs.ErrNotExist) {
		os.Mkdir(yaflPath, 0666)
	}
	if err = os.WriteFile(path.Join(yaflPath, "data.json"), jsonData, 0666); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil
}

func GetData() ([]YAFLInstance, error) {
	var data []YAFLInstance
	yaflPath := filepath.Join(os.Getenv("LOCALAPPDATA"), ".yafl")
	dataPath := filepath.Join(yaflPath, "data.json")
	if _, err := os.Stat(dataPath); errors.Is(err, fs.ErrNotExist) {
		data, err = initData()
		if err != nil {
			return nil, err
		}
	} else if err == nil {
		bData, err := os.ReadFile(dataPath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bData, &data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func CreateInstance(data *[]YAFLInstance, name string, buildPath string) error {
	newInstance := YAFLInstance{
		Name:           name,
		ModsPath:       filepath.Join(buildPath, "Mods"),
		BuildPath:      buildPath,
		AdditionalArgs: "",
	}
	for _, v := range *data {
		if v.Name == name {
			return fmt.Errorf("instance with the name \"%s\" already exists", name)
		}
	}
	*data = append(*data, newInstance)
	return nil
}

func RemoveInstance(data *[]YAFLInstance, name string) {
	for i, v := range *data {
		if v.Name == name {
			(*data)[i] = (*data)[len(*data)-1]
			*data = (*data)[:len(*data)-1]
			break
		}
	}
}

func FetchInstance(data *[]YAFLInstance, name string) *YAFLInstance {
	for _, v := range *data {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func SaveData(data *[]YAFLInstance) error {
	jsonData, err := json.Marshal(*data)
	if err != nil {
		return err
	}
	yaflPath := filepath.Join(os.Getenv("LOCALAPPDATA"), ".yafl")
	if _, err = os.Stat(yaflPath); errors.Is(err, fs.ErrNotExist) {
		os.Mkdir(yaflPath, 0666)
	}
	if err = os.WriteFile(filepath.Join(yaflPath, "data.json"), jsonData, 0666); err != nil {
		return err
	}
	return nil
}
