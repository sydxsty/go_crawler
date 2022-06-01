package crawler

import (
	"bytes"
	"crawler/neubt"
	"crawler/util/html"
	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
	"mime/multipart"
	"strings"
)

type TorrentModifier struct {
	*TorrentPosterImpl
	updateFID string
	updateTID string
	updatePID string
}

func NewTorrentModifier(url string, client neubt.Client) (*TorrentModifier, error) {
	t := &TorrentModifier{}
	t.TorrentPosterImpl = &TorrentPosterImpl{
		client: client.Clone(),
	}
	resp, err := t.client.SyncVisit(url)
	if err != nil {
		return nil, err
	}
	node, err := html.NewNodeFromBytes(resp.Body)
	if err != nil {
		return nil, err
	}
	t.tid, _ = node.GetInnerString(`.//option[@selected="selected"]/@value`) // 222

	t.formHash, _ = node.GetInnerString(`.//input[@id="formhash"]/@value`)
	t.postTime, _ = node.GetInnerString(`.//input[@id="posttime"]/@value`)
	t.wysiwyg, _ = node.GetInnerString(`.//input[@name="wysiwyg"]/@value`)
	t.special, _ = node.GetInnerString(`.//input[@name="special"]/@value`)
	t.specialExtra, _ = node.GetInnerString(`.//input[@name="specialextra"]/@value`)
	t.subject, _ = node.GetInnerString(`.//input[@name="subject"]/@value`)

	t.updateFID, _ = node.GetInnerString(`.//input[@name="fid"]/@value`)
	t.updateTID, _ = node.GetInnerString(`.//input[@name="tid"]/@value`)
	t.updatePID, _ = node.GetInnerString(`.//input[@name="pid"]/@value`)
	// build the image loader
	uid, _ := node.GetInnerString(`.//input[@name="uid"]/@value`)
	hash, _ := node.GetInnerString(`.//input[@name="hash"]/@value`)
	t.imgUploader = NewImageUploader(client, uid, hash)

	tidNodeList, err := node.GetInnerNodeList(`.//select[@name="typeid"]/option`)
	if err != nil {
		return nil, err
	}
	t.tidList = make(map[string]string)
	for _, tidNode := range tidNodeList {
		t.tidList[tidNode.GetString()], _ = tidNode.GetInnerString(`./@value`)
	}
	return t, nil
}

func (t *TorrentModifier) UpdateTorrentMultiPart() error {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	_ = w.WriteField(UTF82GB2312("formhash"), UTF82GB2312(t.formHash))
	_ = w.WriteField(UTF82GB2312("posttime"), UTF82GB2312(t.postTime))
	_ = w.WriteField(UTF82GB2312("delattachop"), UTF82GB2312("0")) //1
	_ = w.WriteField(UTF82GB2312("wysiwyg"), UTF82GB2312(t.wysiwyg))
	_ = w.WriteField(UTF82GB2312("fid"), UTF82GB2312(t.updateFID)) //1
	_ = w.WriteField(UTF82GB2312("tid"), UTF82GB2312(t.updateTID)) //1
	_ = w.WriteField(UTF82GB2312("pid"), UTF82GB2312(t.updatePID)) //1
	_ = w.WriteField(UTF82GB2312("page"), UTF82GB2312("1"))        //1
	_ = w.WriteField(UTF82GB2312("special"), UTF82GB2312(t.special))
	_ = w.WriteField(UTF82GB2312("specialextra"), UTF82GB2312(t.specialExtra))
	_ = w.WriteField(UTF82GB2312("typeid"), UTF82GB2312(t.tid))
	_ = w.WriteField(UTF82GB2312("subject"), UTF82GB2312(t.subject))
	_ = w.WriteField(UTF82GB2312("message"), UTF82GB2312(t.getMessageBody()))
	_ = w.WriteField(UTF82GB2312("readperm"), UTF82GB2312(""))
	_ = w.WriteField(UTF82GB2312("tags"), UTF82GB2312(""))
	_ = w.WriteField(UTF82GB2312("allownoticeauthor"), UTF82GB2312("1"))
	_ = w.WriteField(UTF82GB2312("usesig"), UTF82GB2312("1"))
	_ = w.WriteField(UTF82GB2312("save"), UTF82GB2312(""))
	for _, aid := range append(t.aidList, t.thumbAidList...) {
		_ = w.WriteField(UTF82GB2312("attachupdate["+aid+"]"), UTF82GB2312(""))
		_ = w.WriteField(UTF82GB2312("attachnew["+aid+"][description]"), UTF82GB2312(""))
		_ = w.WriteField(UTF82GB2312("attachnew["+aid+"][readperm]"), UTF82GB2312(""))
		_ = w.WriteField(UTF82GB2312("attachnew["+aid+"][price]"), UTF82GB2312("0"))
	}
	t.client.SetRequestCallback(func(r *colly.Request) {
		r.Headers.Set("Content-Type", w.FormDataContentType())
	})
	resp, err := t.client.SyncPostRaw(`/forum.php?mod=post&action=edit&extra=&editsubmit=yes`, body.Bytes())
	if err != nil {
		return err
	}
	if strings.Contains(resp.Request.URL.RawQuery, "forum.php") {
		return errors.New("error publish torrent")
	}
	return nil
}
