package tinykv

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

type httpGetter struct {
	baseURL string
}

var bufferPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	url := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	b := bufferPool.Get().(*bytes.Buffer)
	b.Reset()
	defer bufferPool.Put(b)
	_, err = io.Copy(b, res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return b.Bytes(), nil
}
