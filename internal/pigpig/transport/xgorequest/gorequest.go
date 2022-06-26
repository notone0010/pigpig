// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// xgorequest
package xgorequest

import (
	"compress/flate"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
	"github.com/marmotedu/errors"
	"github.com/notone/pigpig/internal/pigpig/transport"
	"github.com/notone/pigpig/internal/pigpig/dudu"
	"github.com/notone/pigpig/pkg/log"
	"github.com/parnurzeal/gorequest"
)

type RequestTransport struct {
	// you can set some options
}

var _ transport.GoRequestTransport = (*RequestTransport)(nil)

var requestEngine transport.Factory

func (e RequestTransport) GoRequest() transport.GoRequestTransport {
	return &RequestTransport{}
}

func (e RequestTransport) Close() error {
	return nil
}

func (e RequestTransport) FetchRemoteResponse(c *dudu.RequestDetail) (*http.Response, []byte, error) {
	var (
		err error
	)
	sa := gorequest.New()
	if err = PrepareRequestDetail(sa, c); err != nil {
		log.Error(err.Error())
		return nil, nil, err
	}
	var (
		resp   gorequest.Response
		body   []byte
		reqErr []error
	)
	resp, body, reqErr = RequestRemote(sa)
	if reqErr != nil && len(reqErr) != 0 {
		err = errors.NewAggregate(reqErr)
		return nil, nil, err
	}
	if !c.DisableCompression {
		v := resp.Header.Get("Content-Encoding")
		if v != "" {
			resp.Header.Set("X-Pigpig-Origin-Content-Length", strconv.Itoa(len(body)))
			body, err = UncompressContentEncoding(resp.Header.Get("Content-Encoding"), resp.Body)
			resp.Header.Set("X-Pigpig-Origin-Content-Encoding", v)
			if err == nil {
				c.UnCompression = true
				resp.Header.Del("Content-Encoding")
			} else {
				c.UnCompression = false
			}
		}
	}

	if v := resp.Header.Get("Connection"); v != "" {
		resp.Header.Set("X-Pigpig-Origin-Connection", resp.Header.Get("Connection"))
		resp.Header.Del("Connection")
	}
	return resp, body, err
}

// UncompressContentEncoding will hand content from remote response
func UncompressContentEncoding(compressionType string, content io.Reader) (uncompressContent []byte, err error) {
	switch compressionType {
	case "gzip":
		var reader *gzip.Reader
		reader, err = gzip.NewReader(content)
		if err != nil {
			return uncompressContent, err
		}
		defer reader.Close()
		uncompressContent, err = ioutil.ReadAll(reader)

		return uncompressContent, err
	case "deflate":
		var (
			reader io.ReadCloser
		)
		reader = flate.NewReader(content)
		defer reader.Close()

		uncompressContent, err = ioutil.ReadAll(reader)
		return uncompressContent, err
	case "br":
		var reader *brotli.Reader

		reader = brotli.NewReader(content)

		uncompressContent, err = ioutil.ReadAll(reader)
		return uncompressContent, err

	default:
		uncompressContent, err = ioutil.ReadAll(content)
		return uncompressContent, err
	}
}

// PrepareRequestDetail will prepare some options and request-details when a request fetch to remote
func PrepareRequestDetail(s *gorequest.SuperAgent, c *dudu.RequestDetail) error {
	// 更新请求body content-length
	// 携带代理
	s.Header = c.Request.Header
	if c.Proxy != "" {
		addr := net.ParseIP(strings.Split(c.Proxy, ":")[0])
		if addr == nil {
			return errors.New("invalid ip address")
		} else {
			s.Proxy(c.Proxy)
		}
	}
	reqRemoteUrl := c.Request.Header.Get(dudu.InternalHeaderFullPath)
	reqRemoteURL, _ := url.Parse(reqRemoteUrl)
	if !reqRemoteURL.IsAbs() {
		scheme := c.Protocol
		if strings.Contains(scheme, "https") || strings.Contains(scheme, "HTTPS") || c.Request.TLS != nil {
			reqRemoteURL.Scheme = "https"
		}
		reqRemoteURL.Host = c.Host

	}

	s.Url = reqRemoteURL.String()
	s.Method = c.Method
	s.Cookies = c.Cookies()
	s.Transport.DisableCompression = true
	s.Transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: false,
	}
	s.RedirectPolicy(func(req gorequest.Request, via []gorequest.Request) error {
		if via != nil && len(via) != 0 {
			log.Debugf("%s wants move to %s", via[0].URL.String(), req.URL.String())

		}
		return http.ErrUseLastResponse
	})
	s.FormData = c.RequestData
	return nil
}

// // PrepareRequestBody will handle a request body
// func PrepareRequestBody(s *gorequest.SuperAgent) {
// 	// 更新新的content-length
// 	// 记录旧的content-length
// }

// RequestRemote send request to remote server
func RequestRemote(s *gorequest.SuperAgent) (response gorequest.Response, body []byte, errs []error) {
	var _body string
	response, _body, errs = s.End()
	body = []byte(_body)
	return
}

func GetGorequestTransport() transport.Factory {
	requestEngine = &RequestTransport{}
	return requestEngine
}
