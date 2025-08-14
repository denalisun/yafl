package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

type YAFLInstance struct {
	Name           string
	Mods           []string
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
	yaflPath := path.Join(os.Getenv("LOCALAPPDATA"), ".yafl")
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
	yaflPath := path.Join(os.Getenv("LOCALAPPDATA"), ".yafl")
	dataPath := path.Join(yaflPath, "data.json")
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

func AddInstance(data *[]YAFLInstance, name string, buildPath string) {
	newInstance := YAFLInstance{
		Name:           name,
		Mods:           []string{},
		BuildPath:      buildPath,
		AdditionalArgs: "",
	}
	*data = append(*data, newInstance)
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

func SaveData(data *[]YAFLInstance) error {
	jsonData, err := json.Marshal(*data)
	if err != nil {
		return err
	}
	yaflPath := path.Join(os.Getenv("LOCALAPPDATA"), ".yafl")
	if _, err = os.Stat(yaflPath); errors.Is(err, fs.ErrNotExist) {
		os.Mkdir(yaflPath, 0666)
	}
	if err = os.WriteFile(path.Join(yaflPath, "data.json"), jsonData, 0666); err != nil {
		return err
	}
	return nil
}
