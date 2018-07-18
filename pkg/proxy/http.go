package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	charsetUTF8                    = "charset=UTF-8"
	MIMEApplicationJSON            = "application/json"
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8
)

func newJSONResponse(status int, i interface{}) (*http.Response, error) {
	body, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Body:          ioutil.NopCloser(bytes.NewReader(body)),
		StatusCode:    status,
		ContentLength: int64(len(body)),
		Header: http.Header{
			"Content-Type": []string{MIMEApplicationJSONCharsetUTF8},
		},
	}, nil
}

func copyBody(r io.ReadCloser) (io.ReadCloser, []byte) {
	body := make([]byte, 0)
	if r != nil {
		body, _ = ioutil.ReadAll(r)
		r = ioutil.NopCloser(bytes.NewReader(body))
	}
	return r, body
}

func bind(r *http.Request, i interface{}) error {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer closeBody(r)

	return json.Unmarshal(bodyBytes, i)
}

func closeBody(r *http.Request) {
	if r.Body != nil {
		r.Body.Close()
	}
}
