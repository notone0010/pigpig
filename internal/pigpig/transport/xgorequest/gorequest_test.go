// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// xgorequest
package xgorequest

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/notone0010/pigpig/internal/pigpig/dudu"
	"github.com/notone0010/pigpig/internal/pigpig/transport"
	"github.com/notone0010/pigpig/pkg/log"
	"github.com/panjf2000/ants"
	"github.com/valyala/fasthttp"
)

func TestGorequest(t *testing.T) {
	// Send()
	requestTransport := GetGorequestTransport()
	var wg sync.WaitGroup
	antPool, _ := ants.NewPool(100)
	defer antPool.Release()

	start := time.Now()
	counter := 0
	for i := 0; i < 500; i++ {
		wg.Add(1)
		counter++

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			requestGO(requestTransport)
		}(&wg)
		if counter%100 == 0 {
			wg.Wait()
		}
		// sendWorkFast(&wg)()
		// if err := antPool.Submit(func(wg *sync.WaitGroup) func() {
		// 	// defer wg.Done()
		// 	return sendWork(wg)
		// 	// return func() {
		// 	//
		// 	// }
		// }(&wg)); err != nil {
		// 	log.Fatal(err.Error())
		// }
		// go func() {
		// 	defer wg.Done()
		// 	get()
		// 	// requestGO(requestTransport)
		// }()
	}
	wg.Wait()
	end := time.Now().Sub(start)
	fmt.Println("elapsed: ", end.Seconds())
	// fmt.Println("body: ", body[:100])
}

func requestGO(requestTransport transport.Factory) {
	req, err := http.NewRequest("GET", "http://192.168.1.100:8000/", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Set(dudu.InternalHeaderFullPath, "http://192.168.1.100:8000/")
	options := &dudu.RequestOptions{
		Method:   "GET",
		Hostname: "192.168.1.100:8000",
		Path:     "http://192.168.1.100:8000/",
		Header:   req.Header,
	}

	requestDetail := &dudu.RequestDetail{
		Instance:           "192.168.1.100:8000",
		RequestURL:         "http://192.168.1.100:8000/",
		UseType:            "gstunnel",
		Protocol:           "HTTP1.1",
		RequestOptions:     options,
		Request:            req,
		RequestData:        url.Values{},
		DisableCompression: false,
		RemoteAddr:         "",

		CreateAt: time.Now(),
	}
	resp, _, err := requestTransport.GoRequest().FetchRemoteResponse(requestDetail)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	fmt.Println("response: ", resp.Header.Get("Server"))
}

func sendWork(wg *sync.WaitGroup) func() {
	defer wg.Done()
	return get
}

func get() {
	r, err := http.Get("http://192.168.1.100:8000/")
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer r.Body.Close()

	fmt.Println(r.Header.Get("Server"))
}

func sendWorkFast(wg *sync.WaitGroup) func() {
	defer wg.Done()
	return fastHttpGet
}

func fastHttpGet() {
	req := &fasthttp.Request{}
	req.SetRequestURI("http://192.168.1.100:8000/")
	req.Header.SetMethod("GET")
	resp := &fasthttp.Response{}
	client := &fasthttp.Client{}

	if err := client.Do(req, resp); err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println(string(resp.Header.Server()))
}
