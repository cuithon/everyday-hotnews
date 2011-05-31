// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Primitive HTTP client. See RFC 2616.

package http

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

// A Client is an HTTP client. Its zero value (DefaultClient) is a usable client
// that uses DefaultTransport.
// Client is not yet very configurable.
type Client struct {
	Transport RoundTripper // if nil, DefaultTransport is used

	// If CheckRedirect is not nil, the client calls it before
	// following an HTTP redirect. The arguments req and via
	// are the upcoming request and the requests made already,
	// oldest first. If CheckRedirect returns an error, the client
	// returns that error instead of issue the Request req.
	//
	// If CheckRedirect is nil, the Client uses its default policy,
	// which is to stop after 10 consecutive requests.
	CheckRedirect func(req *Request, via []*Request) os.Error
}

// DefaultClient is the default Client and is used by Get, Head, and Post.
var DefaultClient = &Client{}

// RoundTripper is an interface representing the ability to execute a
// single HTTP transaction, obtaining the Response for a given Request.
type RoundTripper interface {
	// RoundTrip executes a single HTTP transaction, returning
	// the Response for the request req.  RoundTrip should not
	// attempt to interpret the response.  In particular,
	// RoundTrip must return err == nil if it obtained a response,
	// regardless of the response's HTTP status code.  A non-nil
	// err should be reserved for failure to obtain a response.
	// Similarly, RoundTrip should not attempt to handle
	// higher-level protocol details such as redirects,
	// authentication, or cookies.
	//
	// RoundTrip may modify the request. The request Headers field is
	// guaranteed to be initialized.
	RoundTrip(req *Request) (resp *Response, err os.Error)
}

// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

// Used in Send to implement io.ReadCloser by bundling together the
// bufio.Reader through which we read the response, and the underlying
// network connection.
type readClose struct {
	io.Reader
	io.Closer
}

// Do sends an HTTP request and returns an HTTP response, following
// policy (e.g. redirects, cookies, auth) as configured on the client.
//
// Callers should close resp.Body when done reading from it.
//
// Generally Get, Post, or PostForm will be used instead of Do.
func (c *Client) Do(req *Request) (resp *Response, err os.Error) {
	if req.Method == "GET" || req.Method == "HEAD" {
		return c.doFollowingRedirects(req)
	}
	return send(req, c.Transport)
}


// send issues an HTTP request.  Caller should close resp.Body when done reading from it.
func send(req *Request, t RoundTripper) (resp *Response, err os.Error) {
	if t == nil {
		t = DefaultTransport
		if t == nil {
			err = os.NewError("no http.Client.Transport or http.DefaultTransport")
			return
		}
	}

	// Most the callers of send (Get, Post, et al) don't need
	// Headers, leaving it uninitialized.  We guarantee to the
	// Transport that this has been initialized, though.
	if req.Header == nil {
		req.Header = make(Header)
	}

	info := req.URL.RawUserinfo
	if len(info) > 0 {
		if req.Header == nil {
			req.Header = make(Header)
		}
		req.Header.Set("Authorization", "Basic "+base64.URLEncoding.EncodeToString([]byte(info)))
	}
	return t.RoundTrip(req)
}

// True if the specified HTTP status code is one for which the Get utility should
// automatically redirect.
func shouldRedirect(statusCode int) bool {
	switch statusCode {
	case StatusMovedPermanently, StatusFound, StatusSeeOther, StatusTemporaryRedirect:
		return true
	}
	return false
}

// Get issues a GET to the specified URL.  If the response is one of the following
// redirect codes, Get follows the redirect, up to a maximum of 10 redirects:
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//
// Caller should close r.Body when done reading from it.
//
// Get is a convenience wrapper around DefaultClient.Get.
func Get(url string) (r *Response, err os.Error) {
	return DefaultClient.Get(url)
}

// Get issues a GET to the specified URL.  If the response is one of the
// following redirect codes, Get follows the redirect after calling the
// Client's CheckRedirect function.
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//
// Caller should close r.Body when done reading from it.
func (c *Client) Get(url string) (r *Response, err os.Error) {
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.doFollowingRedirects(req)
}

func (c *Client) doFollowingRedirects(ireq *Request) (r *Response, err os.Error) {
	// TODO: if/when we add cookie support, the redirected request shouldn't
	// necessarily supply the same cookies as the original.
	var base *URL
	redirectChecker := c.CheckRedirect
	if redirectChecker == nil {
		redirectChecker = defaultCheckRedirect
	}
	var via []*Request

	req := ireq
	url := "" // next relative or absolute URL to fetch (after first request)
	for redirect := 0; ; redirect++ {
		if redirect != 0 {
			req = new(Request)
			req.Method = ireq.Method
			req.Header = make(Header)
			req.URL, err = base.ParseURL(url)
			if err != nil {
				break
			}
			if len(via) > 0 {
				// Add the Referer header.
				lastReq := via[len(via)-1]
				if lastReq.URL.Scheme != "https" {
					req.Referer = lastReq.URL.String()
				}

				err = redirectChecker(req, via)
				if err != nil {
					break
				}
			}
		}

		url = req.URL.String()
		if r, err = send(req, c.Transport); err != nil {
			break
		}
		if shouldRedirect(r.StatusCode) {
			r.Body.Close()
			if url = r.Header.Get("Location"); url == "" {
				err = os.ErrorString(fmt.Sprintf("%d response missing Location header", r.StatusCode))
				break
			}
			base = req.URL
			via = append(via, req)
			continue
		}
		return
	}

	method := ireq.Method
	err = &URLError{method[0:1] + strings.ToLower(method[1:]), url, err}
	return
}

func defaultCheckRedirect(req *Request, via []*Request) os.Error {
	if len(via) >= 10 {
		return os.ErrorString("stopped after 10 redirects")
	}
	return nil
}

// Post issues a POST to the specified URL.
//
// Caller should close r.Body when done reading from it.
//
// Post is a wrapper around DefaultClient.Post
func Post(url string, bodyType string, body io.Reader) (r *Response, err os.Error) {
	return DefaultClient.Post(url, bodyType, body)
}

// Post issues a POST to the specified URL.
//
// Caller should close r.Body when done reading from it.
func (c *Client) Post(url string, bodyType string, body io.Reader) (r *Response, err os.Error) {
	req, err := NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	return send(req, c.Transport)
}

// PostForm issues a POST to the specified URL, 
// with data's keys and values urlencoded as the request body.
//
// Caller should close r.Body when done reading from it.
//
// PostForm is a wrapper around DefaultClient.PostForm
func PostForm(url string, data map[string]string) (r *Response, err os.Error) {
	return DefaultClient.PostForm(url, data)
}

// PostForm issues a POST to the specified URL, 
// with data's keys and values urlencoded as the request body.
//
// Caller should close r.Body when done reading from it.
func (c *Client) PostForm(url string, data map[string]string) (r *Response, err os.Error) {
	return c.Post(url, "application/x-www-form-urlencoded", urlencode(data))
}

// TODO: remove this function when PostForm takes a multimap.
func urlencode(data map[string]string) (b *bytes.Buffer) {
	m := make(map[string][]string, len(data))
	for k, v := range data {
		m[k] = []string{v}
	}
	return bytes.NewBuffer([]byte(EncodeQuery(m)))
}

// Head issues a HEAD to the specified URL.  If the response is one of the
// following redirect codes, Head follows the redirect after calling the
// Client's CheckRedirect function.
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//
// Head is a wrapper around DefaultClient.Head
func Head(url string) (r *Response, err os.Error) {
	return DefaultClient.Head(url)
}

// Head issues a HEAD to the specified URL.  If the response is one of the
// following redirect codes, Head follows the redirect after calling the
// Client's CheckRedirect function.
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
func (c *Client) Head(url string) (r *Response, err os.Error) {
	req, err := NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return c.doFollowingRedirects(req)
}
