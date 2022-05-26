package crawler

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestUploadImage(t *testing.T) {
	iu := NewImageUploader(client, os.Getenv("uid"), os.Getenv("hash"))
	file, err := ioutil.ReadFile(os.Getenv("fileName"))
	assert.NoError(t, err, "error load image")
	aid, err := iu.UploadImage(file, "jpg")
	assert.NoError(t, err, "error upload image")
	err = iu.UseImageInPID("0")
	assert.NoError(t, err, "error using image")
	err = iu.RemoveImage(aid)
	assert.NoError(t, err, "error remove image")
}

func TestGetImageFromText(t *testing.T) {
	text := "111[img]https://lain.bgm.tv/pic/cover/l/de/4a/329906_hmtVD_1.jpg[/img]11\n" +
		"1111[img]https://lain.bgm.tv/pic/cover/l/de/4a/329906_hmtVD_2.jpg[/img]1111\n" +
		"1[img]https://lain.bgm.tv/pic/cover/l/de/4a/329906_hmtVD_3.jpg[/ano]111\n" +
		"1[img]https://lain.bgm.tv/pic/cover/l/de/4a/329906_hmtVD_3.jpg[/img]111\n"
	result1 := GetAllImgFromPTGen(text)
	log.Print(result1)
	result2, err := ReplaceImgWithTagID(text, []string{"1000", "1001", "1002"})
	assert.NoError(t, err, "error replace")
	log.Print(result2)
}
