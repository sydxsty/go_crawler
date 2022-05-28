package dao

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSetCHSName(t *testing.T) {
	name1 := "LoveLive! 虹咲学园学园偶像同好会"
	name2 := "LoveLive! 虹咲学园学园偶像同好会 第2期"
	replace1 := "LL-1"
	replace2 := "LL-2"
	title := `【极影字幕社】LoveLive! 虹咲学园学园偶像同好会 第2期 第08集 加料剪辑版 GB_CN HEVC_opus 1080p`
	_ = os.WriteFile(`./data/names.yaml`, nil, 0666)
	a, err := NewAnimeDB()
	assert.NoError(t, err, "can not create object")
	assert.True(t, a != nil)
	// test replacement
	err = a.AddNewCHSName(name1, "")
	assert.NoError(t, err, "can not add name")
	result := a.GetAliasCHSName(title)
	assert.True(t, result == name1)
	err = a.AddNewCHSName(name1, replace1)
	assert.NoError(t, err, "can not add name")
	result = a.GetAliasCHSName(title)
	assert.True(t, result == replace1)
	// test override
	err = a.AddNewCHSName(name2, "")
	assert.NoError(t, err, "can not add name")
	result = a.GetAliasCHSName(title)
	assert.True(t, result == name2)
	err = a.AddNewCHSName(name2, replace2)
	assert.NoError(t, err, "can not add name")
	result = a.GetAliasCHSName(title)
	assert.True(t, result == replace2)
}
