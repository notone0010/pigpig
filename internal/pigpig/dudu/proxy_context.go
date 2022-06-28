// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dudu

import (
	"io"
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// abortIndex represents a typical value used in abort functions.
const abortIndex int8 = math.MaxInt8 >> 1

// Context custom context.
type Context struct {
	writermem responseWriter

	Request *http.Request

	Writer ResponseWriter

	RequestDetail *RequestDetail

	ResponseDetail *ResponseDetail

	// Handlers HandlersChain
	Handlers HandlersChain

	engine *ProxyHttpMux

	index int8

	// This mutex protect Keys map
	mu sync.RWMutex

	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]interface{}

	Errors []error

	fullPath string
}

// ResponseDetail ...
type ResponseDetail struct {
	// StatusCode the code that response
	StatusCode int `json:"status_code"`

	// ElapsedTime elapsed time while the request fetch remote
	ElapsedTime time.Duration `json:"elapsed_time"` // 客户端请求所消耗的时间

	// Header response header
	Header http.Header

	// RawBody
	RawBody io.Reader

	// Body final body
	Body []byte

	// Response
	*http.Response
}

// RequestDetail 代理请求客户端对象.
type RequestDetail struct {
	// Instance the client request instance, for example, tianyancha.com
	Instance string `json:"instance"`

	RequestOptions *RequestOptions

	Protocol string `json:"protocol"`

	// Proxy the proxy is assigned to the client
	Proxy string `json:"proxy"` // 客户端分配的代理IP，default: 0.0.0.0:65536

	// RequestURL the request`s url
	RequestURL string `json:"request_url"`

	RequestData url.Values `json:"request_data"`

	UseType string `json:"use_type"`

	// Remoter The visitors ip addr
	RemoteAddr string

	// DisableCompression, if true, prevents the Transport from
	// requesting compression with an "Accept-Encoding: gzip"
	// request header when the Request contains no existing
	// Accept-Encoding value. If the Transport requests gzip on
	// its own and gets a gzipped response, it's transparently
	// decoded in the Response.Body. However, if the user
	// explicitly requested gzip it is not automatically
	// uncompressed.
	DisableCompression bool

	// UnCompression expressed remote response content is not compression
	UnCompression bool

	// // IsKeep label ture if the request is served first
	// IsKeep bool `json:"is_keep"`
	//
	// // IsChange records the request by be changed as the IsKeep is ture
	// IsChange bool `json:"is_change"`

	// Plugins []string

	*http.Request

	// Extra
	Extra interface{} `json:"extra"`

	// CreateAt records the request created time
	CreateAt time.Time `json:"create_at"`

	// OutOffAt records the request out off time
	OutOffAt time.Time `json:"out_off_at"`
}

// RequestOptions request options.
type RequestOptions struct {
	Hostname string

	Port string

	Path string

	Method string

	Header http.Header
}

// NewContext returns new Context.
func NewContext(engine *ProxyHttpMux) *Context {
	return &Context{engine: engine}
}

// FullPath returns a request full path
// url: /login  calls FullPath() -> http://example.com/login/
func (c *Context) FullPath() string {
	return c.fullPath
}

// SerializeFullPath returns a request full path
// url: /login  calls SerializeFullPath() -> http://example.com/login/
func (c *Context) SerializeFullPath() string {
	return c.Request.URL.String()
}

// InitContext initial context.
func (c *Context) InitContext() {
	c.Writer = &c.writermem
	c.fullPath = c.SerializeFullPath()
	// c.handlers = nil
	c.NewPrepareRequest()
}

// NewPrepareRequest returns NewPrepareRequest.
func (c *Context) NewPrepareRequest() {
	options := &RequestOptions{
		Method:   c.Request.Method,
		Hostname: c.Request.Host,
		Path:     c.Request.RequestURI,
		Header:   c.Request.Header,
	}
	detail := &RequestDetail{
		Instance:           c.Request.Host,
		RequestURL:         c.Request.RequestURI,
		UseType:            "gstunnel",
		Protocol:           c.Request.Proto,
		RequestOptions:     options,
		Request:            c.Request,
		RequestData:        c.Request.PostForm,
		DisableCompression: false,
		RemoteAddr:         c.Request.RemoteAddr,

		CreateAt: time.Now(),
	}
	// if v, exist := c.Request.Header["USE-RULE"]; exist {
	// 	if v[0] == "true" {
	// 		detail.Plugins = append(detail.Plugins, "ProxyPlugin")
	// 	}
	// }
	c.RequestDetail = detail
}

// GetContextObj reset context object.
func (c *Context) GetContextObj(w http.ResponseWriter, r *http.Request, engine *ProxyHttpMux) {
	c.writermem.reset(w)
	c.Writer = &c.writermem
	c.Request = r
	c.engine = engine
	c.RequestDetail = nil
	c.ResponseDetail = nil
	c.index = -1
	c.Handlers = nil
	c.fullPath = ""
	c.Keys = nil
	c.Errors = c.Errors[:0]
	c.InitContext()
}

/************************************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

// Deadline always returns that there is no deadline (ok==false),
// maybe you want to use Request.Context().Deadline() instead.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done always returns nil (chan which will wait forever),
// if you want to abort your work when the connection was closed
// you should use Request.Context().Done() instead.
func (c *Context) Done() <-chan struct{} {
	return nil
}

// Err always returns nil, maybe you want to use Request.Context().Err() instead.
func (c *Context) Err() error {
	return nil
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *Context) Value(key interface{}) interface{} {
	if key == 0 {
		return c.Request
	}
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Set is used to discover a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	c.Keys[key] = value
	c.mu.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false).
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetUint returns the value associated with the key as an unsigned integer.
func (c *Context) GetUint(key string) (ui uint) {
	if val, ok := c.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (c *Context) GetUint64(key string) (ui64 uint64) {
	if val, ok := c.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}

	return
}

// GetDuration returns the value associated with the key as a duration.
func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}

	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}

	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}

	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

// GetHeader returns value from request headers.
func (c *Context) GetHeader(key string) string {
	return c.requestHeader(key)
}

func (c *Context) requestHeader(key string) string {
	return c.Request.Header.Get(key)
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = abortIndex
}

// IsAborted context whether abort.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.Handlers)) {
		c.Handlers[c.index](c)
		c.index++
	}
}
