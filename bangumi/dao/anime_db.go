package dao

import (
	"gopkg.in/yaml.v3"
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
	if f, err := os.OpenFile("./data/names.yaml", os.O_CREATE|os.O_RDONLY, 0666); err != nil {
		return nil, err
	} else {
		_ = yaml.NewDecoder(f).Decode(&a.stateList)
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
	a.stateList = append(a.stateList, NamePair{Name: name, Replace: replace})
	if f, err := os.OpenFile("./data/names.yaml", os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		return err
	} else {
		if err = yaml.NewEncoder(f).Encode(a.stateList); err != nil {
			return err
		}
	}
	return nil
}
