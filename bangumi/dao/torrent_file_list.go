package dao

import (
	ac "crawler/bangumi/anime_control"
	"github.com/pkg/errors"
	"strings"
)

const folderPrefix = "|-  "

type TorrentFileList struct {
	name    string
	content TorrentContent
}

type TorrentContent struct {
	isFile bool
	// if is file, print the size
	Size string
	// if is file, content is nil
	content map[string]*TorrentContent
}

func NewTorrentFileList(contents []interface{}) (*TorrentFileList, error) {
	tf := &TorrentFileList{}
	for _, v := range contents {
		vList, ok := v.([]interface{})
		if !ok {
			return nil, errors.Errorf("can not unmarshal %s", v)
		}
		if len(vList) != 2 {
			return nil, errors.Errorf("len(vlist) != 2, %s", vList)
		}
		path := ac.SplitByDelimiter(vList[0].(string), `\/`)
		if len(path) == 0 { // no valid path found
			continue
		}
		err := tf.Upsert(path, vList[1].(string))
		if err != nil {
			return nil, err
		}
	}
	return tf, nil
}

func (tf *TorrentFileList) Upsert(path []string, size string) error {
	name := path[0]
	if tf.name == "" {
		tf.name = name
	}
	if tf.name != name {
		return errors.Errorf("path inconsistent, %s vs. %s", name, tf.name)
	}
	return tf.content.Upsert(path, size)
}

func (tc *TorrentContent) Upsert(path []string, size string) error {
	// current folder(file) name is path[0], path must >1
	if len(path) == 0 {
		return errors.Errorf("path string == 0, this is a bug, %s", path)
	}
	if len(path) == 1 {
		tc.isFile = true
		tc.Size = size
		return nil
	}
	// is folder
	tc.isFile = false
	subPath := path[1:]
	if tc.content == nil {
		tc.content = make(map[string]*TorrentContent)
	}
	subContent, ok := tc.content[subPath[0]]
	if !ok {
		subContent = &TorrentContent{}
		tc.content[subPath[0]] = subContent
	}
	// recurse into sub dir
	return subContent.Upsert(subPath, size)
}

func (tf *TorrentFileList) PrintToStringList() ([]string, error) {
	var result []string
	if tf.name == "" {
		return nil, errors.New("empty tfl")
	}
	result = append(result, tf.name+"\t"+tf.content.Size)
	result = append(result, tf.content.PrintToStringList(folderPrefix)...)
	return result, nil
}

func (tc *TorrentContent) PrintToStringList(prefix string) []string {
	if tc.isFile {
		return nil
	}
	// is folder
	var result []string
	rep := strings.Repeat(" ", len(prefix))
	for _, i := range []bool{false, true} {
		for k, v := range tc.content {
			if v.isFile == i {
				continue
			}
			// append sub folder name
			result = append(result, prefix+k+"\t"+v.Size)
			result = append(result, v.PrintToStringList(rep+folderPrefix)...)
		}
	}
	return result
}

func (tf *TorrentFileList) GetTorrentName() string {
	return tf.name
}

func (tf *TorrentFileList) PrintToString(limit int) (string, error) {
	sl, err := tf.PrintToStringList()
	if err != nil {
		return "", err
	}
	// out of range
	if limit > len(sl) || limit <= 0 {
		limit = len(sl)
	}
	result := ""
	for i := 0; i < limit; i++ {
		result += sl[i] + "\n"
	}
	if limit != len(sl) {
		result += "..."
	}
	return result, nil
}
