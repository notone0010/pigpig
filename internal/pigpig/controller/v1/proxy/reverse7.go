// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/notone0010/pigpig/pkg/log"
)

// NewReverseProxy7 returns reverse proxy object for OSI 7 level: application level.
func NewReverseProxy7(targetHost string) (*httputil.ReverseProxy, error) {
	_url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(_url)

	originalDirector := reverseProxy.Director
	reverseProxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}

	// reverseProxy.ModifyResponse = modifyResponse()
	reverseProxy.ErrorHandler = errorHandler()
	return reverseProxy, nil
}

func modifyRequest(req *http.Request) {
	req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		log.Debugf("Got error while modifying response: %v \n", err)
	}
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		defer resp.Body.Close()

		return nil
	}
}

// GoReverseProxy handles the http request using proxy.
func GoReverseProxy(targetHost string) func(http.ResponseWriter, *http.Request) {
	reverseProxy, err := NewReverseProxy7(targetHost)
	if err != nil {
		log.Errorf(err.Error())
		return nil
	}
	return func(w http.ResponseWriter, r *http.Request) {
		reverseProxy.ServeHTTP(w, r)
	}
}
