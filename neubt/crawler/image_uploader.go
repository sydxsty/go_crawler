package crawler

import (
	"bytes"
	"crawler/neubt"
	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
	"mime/multipart"
	"net/url"
	"regexp"
	"strings"
)

// ImageUploader upload the image to the specific thread
type ImageUploader interface {
	// UploadImage upload the image, return the image AID
	UploadImage(data []byte, fileType string) (string, error)
	RemoveImage(aid string) error
}

type ImageUploaderImpl struct {
	client         neubt.Client
	uploadURL      string
	deleteImageURL string
	uid            string
	hash           string
}

func NewImageUploader(client neubt.Client, uid, hash string) ImageUploader {
	i := &ImageUploaderImpl{
		client:         client.Clone(),
		uploadURL:      `/misc.php?mod=swfupload&operation=upload&simple=1&type=image`,
		deleteImageURL: `/forum.php?mod=ajax&action=deleteattach&inajax=yes&tid=0&pid=0&aids[]=`,
		uid:            uid,
		hash:           hash,
	}
	return i
}

func (i *ImageUploaderImpl) UploadImage(data []byte, fileType string) (string, error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	_ = w.WriteField(UTF82GB2312("uid"), UTF82GB2312(i.uid))
	_ = w.WriteField(UTF82GB2312("hash"), UTF82GB2312(i.hash))

	pa, _ := w.CreateFormFile(UTF82GB2312("Filedata"), UTF82GB2312("poster."+fileType))
	if _, err := pa.Write(data); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	i.client.SetRequestCallback(func(r *colly.Request) {
		r.Headers.Set("Content-Type", w.FormDataContentType())
	})
	resp, err := i.client.SyncPostRaw(i.uploadURL, body.Bytes())
	if err != nil {
		return "", err
	}
	// response body like DISCUZUPLOAD|0|5427062|1|0
	if v := regexp.MustCompile(`([0-9]{5,})`).FindAllString(string(resp.Body), -1); len(v) == 0 {
		return "", errors.New("no matching aid found")
	} else {
		return v[0], nil
	}
}

func (i *ImageUploaderImpl) RemoveImage(aid string) error {
	_, err := i.client.SyncVisit(i.deleteImageURL + url.QueryEscape(aid))
	return err
}

func GetAllImgFromPTGen(text string) []string {
	var result []string
	for _, r := range regexp.MustCompile(`(\[img])(.*)(\[/img])`).FindAllStringSubmatch(text, -1) {
		result = append(result, r[2])
	}
	return result
}

func ReplaceImgWithTagID(text string, aid []string) (string, error) {
	target := regexp.MustCompile(`(\[img])(.*)(\[/img])`).FindAllString(text, -1)
	if len(target) != len(aid) {
		return "", errors.New("the size of image and aid not equal")
	}
	for i, r := range target {
		text = strings.Replace(text, r, "[attachimg]"+aid[i]+"[/attachimg]", 1)
	}
	return text, nil
}
