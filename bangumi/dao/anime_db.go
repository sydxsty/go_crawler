package dao

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type NamePair struct {
	Name    string
	Replace string
}

type AnimeDB struct {
	stateList []NamePair
}

func NewAnimeDB() (*AnimeDB, error) {
	a := &AnimeDB{}
	data, err := os.ReadFile("./data/names.yaml")
	if err != nil {
		return a, nil
	}
	err = yaml.Unmarshal(data, &a.stateList)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// GetAliasCHSName if the original chs name can not get valid result, try to generate another name
func (a *AnimeDB) GetAliasCHSName(title string) string {
	index := -1
	for i, v := range a.stateList {
		if strings.Contains(title, v.Name) {
			if index != -1 && len(a.stateList[index].Name) > len(v.Name) {
				continue
			}
			// if the new name is longer
			index = i
		}
	}
	if index != -1 {
		if a.stateList[index].Replace != "" {
			return a.stateList[index].Replace
		}
		return a.stateList[index].Name
	}
	return ""
}

func (a *AnimeDB) AddNewCHSName(name, replace string) error {
	for k, v := range a.stateList {
		if name == v.Name { // delete value
			a.stateList = append(a.stateList[:k], a.stateList[k+1:]...)
			break
		}
	}
	return a.appendCHSName(name, replace)
}

func (a *AnimeDB) InsertNewCHSName(name, replace string) error {
	for _, v := range a.stateList {
		if name == v.Name {
			return errors.Errorf("key already exist: %s", name)
		}
	}
	return a.appendCHSName(name, replace)
}

func (a *AnimeDB) appendCHSName(name, replace string) error {
	a.stateList = append(a.stateList, NamePair{Name: name, Replace: replace})
	data, err := yaml.Marshal(a.stateList)
	if err != nil {
		return err
	}
	err = os.WriteFile("./data/names.yaml", data, 0666)
	if err != nil {
		return err
	}
	return nil
}
