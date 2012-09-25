// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// HTTP client implementation. See RFC 2616.
// 
// This is the low-level Transport implementation of RoundTripper.
// The high-level interface is in client.go.

package http

import (
	"bufio"
	"compress/gzip"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// DefaultTransport is the default implementation of Transport and is
// used by DefaultClient.  It establishes a new network connection for
// each call to Do and uses HTTP proxies as directed by the
// $HTTP_PROXY and $NO_PROXY (or $http_proxy and $no_proxy)
// environment variables.
var DefaultTransport RoundTripper = &Transport{Proxy: ProxyFromEnvironment}

// DefaultMaxIdleConnsPerHost is the default value of Transport's
// MaxIdleConnsPerHost.
const DefaultMaxIdleConnsPerHost = 2

// Transport is an implementation of RoundTripper that supports http,
// https, and http proxies (for either http or https with CONNECT).
// Transport can also cache connections for future re-use.
type Transport struct {
	idleLk   sync.Mutex
	idleConn map[string][]*persistConn
	altLk    sync.RWMutex
	altProto map[string]RoundTripper // nil or map of URI scheme => RoundTripper

	// TODO: tunable on global max cached connections
	// TODO: tunable on timeout on cached connections
	// TODO: optional pipelining

	// Proxy specifies a function to return a proxy for a given
	// Request. If the function returns a non-nil error, the
	// request is aborted with the provided error.
	// If Proxy is nil or returns a nil *URL, no proxy is used.
	Proxy func(*Request) (*url.URL, error)

	// Dial specifies the dial function for creating TCP
	// connections.
	// If Dial is nil, net.Dial is used.
	Dial func(net, addr string) (c net.Conn, err error)

	// TLSClientConfig specifies the TLS configuration to use with
	// tls.Client. If nil, the default configuration is used.
	TLSClientConfig *tls.Config

	DisableKeepAlives  bool
	DisableCompression bool

	// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
	// (keep-alive) to keep to keep per-host.  If zero,
	// DefaultMaxIdleConnsPerHost is used.
	MaxIdleConnsPerHost int
}

// ProxyFromEnvironment returns the URL of the proxy to use for a
// given request, as indicated by the environment variables
// $HTTP_PROXY and $NO_PROXY (or $http_proxy and $no_proxy).
// An error is returned if the proxy environment is invalid.
// A nil URL and nil error are returned if no proxy is defined in the
// environment, or a proxy should not be used for the given request.
func ProxyFromEnvironment(req *Request) (*url.URL, error) {
	proxy := getenvEitherCase("HTTP_PROXY")
	if proxy == "" {
		return nil, nil
	}
	if !useProxy(canonicalAddr(req.URL)) {
		return nil, nil
	}
	proxyURL, err := url.Parse(proxy)
	if err != nil || proxyURL.Scheme == "" {
		if u, err := url.Parse("http://" + proxy); err == nil {
			proxyURL = u
			err = nil
		}
	}
	if err != nil {
		return nil, fmt.Errorf("invalid proxy address %q: %v", proxy, err)
	}
	return proxyURL, nil
}

// ProxyURL returns a proxy function (for use in a Transport)
// that always returns the same URL.
func ProxyURL(fixedURL *url.URL) func(*Request) (*url.URL, error) {
	return func(*Request) (*url.URL, error) {
		return fixedURL, nil
	}
}

// transportRequest is a wrapper around a *Request that adds
// optional extra headers to write.
type transportRequest struct {
	*Request        // original request, not to be mutated
	extra    Header // extra headers to write, or nil
}

func (tr *transportRequest) extraHeaders() Header {
	if tr.extra == nil {
		tr.extra = make(Header)
	}
	return tr.extra
}

// RoundTrip implements the RoundTripper interface.
func (t *Transport) RoundTrip(req *Request) (resp *Response, err error) {
	if req.URL == nil {
		return nil, errors.New("http: nil Request.URL")
	}
	if req.Header == nil {
		return nil, errors.New("http: nil Request.Header")
	}
	if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
		t.altLk.RLock()
		var rt RoundTripper
		if t.altProto != nil {
			rt = t.altProto[req.URL.Scheme]
		}
		t.altLk.RUnlock()
		if rt == nil {
			return nil, &badStringError{"unsupported protocol scheme", req.URL.Scheme}
		}
		return rt.RoundTrip(req)
	}
	treq := &transportRequest{Request: req}
	cm, err := t.connectMethodForRequest(treq)
	if err != nil {
		return nil, err
	}

	// Get the cached or newly-created connection to either the
	// host (for http or https), the http proxy, or the http proxy
	// pre-CONNECTed to https server.  In any case, we'll be ready
	// to send it requests.
	pconn, err := t.getConn(cm)
	if err != nil {
		return nil, err
	}

	return pconn.roundTrip(treq)
}

// RegisterProtocol registers a new protocol with scheme.
// The Transport will pass requests using the given scheme to rt.
// It is rt's responsibility to simulate HTTP request semantics.
//
// RegisterProtocol can be used by other packages to provide
// implementations of protocol schemes like "ftp" or "file".
func (t *Transport) RegisterProtocol(scheme string, rt RoundTripper) {
	if scheme == "http" || scheme == "https" {
		panic("protocol " + scheme + " already registered")
	}
	t.altLk.Lock()
	defer t.altLk.Unlock()
	if t.altProto == nil {
		t.altProto = make(map[string]RoundTripper)
	}
	if _, exists := t.altProto[scheme]; exists {
		panic("protocol " + scheme + " already registered")
	}
	t.altProto[scheme] = rt
}

// CloseIdleConnections closes any connections which were previously
// connected from previous requests but are now sitting idle in
// a "keep-alive" state. It does not interrupt any connections currently
// in use.
func (t *Transport) CloseIdleConnections() {
	t.idleLk.Lock()
	m := t.idleConn
	t.idleConn = nil
	t.idleLk.Unlock()
	if m == nil {
		return
	}
	for _, conns := range m {
		for _, pconn := range conns {
			pconn.close()
		}
	}
}

//
// Private implementation past this point.
//

func getenvEitherCase(k string) string {
	if v := os.Getenv(strings.ToUpper(k)); v != "" {
		return v
	}
	return os.Getenv(strings.ToLower(k))
}

func (t *Transport) connectMethodForRequest(treq *transportRequest) (*connectMethod, error) {
	cm := &connectMethod{
		targetScheme: treq.URL.Scheme,
		targetAddr:   canonicalAddr(treq.URL),
	}
	if t.Proxy != nil {
		var err error
		cm.proxyURL, err = t.Proxy(treq.Request)
		if err != nil {
			return nil, err
		}
	}
	return cm, nil
}

// proxyAuth returns the Proxy-Authorization header to set
// on requests, if applicable.
func (cm *connectMethod) proxyAuth() string {
	if cm.proxyURL == nil {
		return ""
	}
	if u := cm.proxyURL.User; u != nil {
		return "Basic " + base64.URLEncoding.EncodeToString([]byte(u.String()))
	}
	return ""
}

// putIdleConn adds pconn to the list of idle persistent connections awaiting
// a new request.
// If pconn is no longer needed or not in a good state, putIdleConn
// returns false.
func (t *Transport) putIdleConn(pconn *persistConn) bool {
	if t.DisableKeepAlives || t.MaxIdleConnsPerHost < 0 {
		pconn.close()
		return false
	}
	if pconn.isBroken() {
		return false
	}
	key := pconn.cacheKey
	max := t.MaxIdleConnsPerHost
	if max == 0 {
		max = DefaultMaxIdleConnsPerHost
	}
	t.idleLk.Lock()
	if t.idleConn == nil {
		t.idleConn = make(map[string][]*persistConn)
	}
	if len(t.idleConn[key]) >= max {
		t.idleLk.Unlock()
		pconn.close()
		return false
	}
	for _, exist := range t.idleConn[key] {
		if exist == pconn {
			log.Fatalf("dup idle pconn %p in freelist", pconn)
		}
	}
	t.idleConn[key] = append(t.idleConn[key], pconn)
	t.idleLk.Unlock()
	return true
}

func (t *Transport) getIdleConn(cm *connectMethod) (pconn *persistConn) {
	key := cm.String()
	t.idleLk.Lock()
	defer t.idleLk.Unlock()
	if t.idleConn == nil {
		return nil
	}
	for {
		pconns, ok := t.idleConn[key]
		if !ok {
			return nil
		}
		if len(pconns) == 1 {
			pconn = pconns[0]
			delete(t.idleConn, key)
		} else {
			// 2 or more cached connections; pop last
			// TODO: queue?
			pconn = pconns[len(pconns)-1]
			t.idleConn[key] = pconns[0 : len(pconns)-1]
		}
		if !pconn.isBroken() {
			return
		}
	}
	panic("unreachable")
}

func (t *Transport) dial(network, addr string) (c net.Conn, err error) {
	if t.Dial != nil {
		return t.Dial(network, addr)
	}
	return net.Dial(network, addr)
}

// getConn dials and creates a new persistConn to the target as
// specified in the connectMethod.  This includes doing a proxy CONNECT
// and/or setting up TLS.  If this doesn't return an error, the persistConn
// is ready to write requests to.
func (t *Transport) getConn(cm *connectMethod) (*persistConn, error) {
	if pc := t.getIdleConn(cm); pc != nil {
		return pc, nil
	}

	conn, err := t.dial("tcp", cm.addr())
	if err != nil {
		if cm.proxyURL != nil {
			err = fmt.Errorf("http: error connecting to proxy %s: %v", cm.proxyURL, err)
		}
		return nil, err
	}

	pa := cm.proxyAuth()

	pconn := &persistConn{
		t:        t,
		cacheKey: cm.String(),
		conn:     conn,
		reqch:    make(chan requestAndChan, 50),
		writech:  make(chan writeRequest, 50),
		closech:  make(chan struct{}),
	}

	switch {
	case cm.proxyURL == nil:
		// Do nothing.
	case cm.targetScheme == "http":
		pconn.isProxy = true
		if pa != "" {
			pconn.mutateHeaderFunc = func(h Header) {
				h.Set("Proxy-Authorization", pa)
			}
		}
	case cm.targetScheme == "https":
		connectReq := &Request{
			Method: "CONNECT",
			URL:    &url.URL{Opaque: cm.targetAddr},
			Host:   cm.targetAddr,
			Header: make(Header),
		}
		if pa != "" {
			connectReq.Header.Set("Proxy-Authorization", pa)
		}
		connectReq.Write(conn)

		// Read response.
		// Okay to use and discard buffered reader here, because
		// TLS server will not speak until spoken to.
		br := bufio.NewReader(conn)
		resp, err := ReadResponse(br, connectReq)
		if err != nil {
			conn.Close()
			return nil, err
		}
		if resp.StatusCode != 200 {
			f := strings.SplitN(resp.Status, " ", 2)
			conn.Close()
			return nil, errors.New(f[1])
		}
	}

	if cm.targetScheme == "https" {
		// Initiate TLS and check remote host name against certificate.
		cfg := t.TLSClientConfig
		if cfg == nil || cfg.ServerName == "" {
			host := cm.tlsHost()
			if cfg == nil {
				cfg = &tls.Config{ServerName: host}
			} else {
				clone := *cfg // shallow clone
				clone.ServerName = host
				cfg = &clone
			}
		}
		conn = tls.Client(conn, cfg)
		if err = conn.(*tls.Conn).Handshake(); err != nil {
			return nil, err
		}
		if t.TLSClientConfig == nil || !t.TLSClientConfig.InsecureSkipVerify {
			if err = conn.(*tls.Conn).VerifyHostname(cm.tlsHost()); err != nil {
				return nil, err
			}
		}
		pconn.conn = conn
	}

	pconn.br = bufio.NewReader(pconn.conn)
	pconn.bw = bufio.NewWriter(pconn.conn)
	go pconn.readLoop()
	go pconn.writeLoop()
	return pconn, nil
}

// useProxy returns true if requests to addr should use a proxy,
// according to the NO_PROXY or no_proxy environment variable.
// addr is always a canonicalAddr with a host and port.
func useProxy(addr string) bool {
	if len(addr) == 0 {
		return true
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	if host == "localhost" {
		return false
	}
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() {
			return false
		}
	}

	no_proxy := getenvEitherCase("NO_PROXY")
	if no_proxy == "*" {
		return false
	}

	addr = strings.ToLower(strings.TrimSpace(addr))
	if hasPort(addr) {
		addr = addr[:strings.LastIndex(addr, ":")]
	}

	for _, p := range strings.Split(no_proxy, ",") {
		p = strings.ToLower(strings.TrimSpace(p))
		if len(p) == 0 {
			continue
		}
		if hasPort(p) {
			p = p[:strings.LastIndex(p, ":")]
		}
		if addr == p || (p[0] == '.' && (strings.HasSuffix(addr, p) || addr == p[1:])) {
			return false
		}
	}
	return true
}

// connectMethod is the map key (in its String form) for keeping persistent
// TCP connections alive for subsequent HTTP requests.
//
// A connect method may be of the following types:
//
// Cache key form                Description
// -----------------             -------------------------
// ||http|foo.com                http directly to server, no proxy
// ||https|foo.com               https directly to server, no proxy
// http://proxy.com|https|foo.com  http to proxy, then CONNECT to foo.com
// http://proxy.com|http           http to proxy, http to anywhere after that
//
// Note: no support to https to the proxy yet.
//
type connectMethod struct {
	proxyURL     *url.URL // nil for no proxy, else full proxy URL
	targetScheme string   // "http" or "https"
	targetAddr   string   // Not used if proxy + http targetScheme (4th example in table)
}

func (ck *connectMethod) String() string {
	proxyStr := ""
	targetAddr := ck.targetAddr
	if ck.proxyURL != nil {
		proxyStr = ck.proxyURL.String()
		if ck.targetScheme == "http" {
			targetAddr = ""
		}
	}
	return strings.Join([]string{proxyStr, ck.targetScheme, targetAddr}, "|")
}

// addr returns the first hop "host:port" to which we need to TCP connect.
func (cm *connectMethod) addr() string {
	if cm.proxyURL != nil {
		return canonicalAddr(cm.proxyURL)
	}
	return cm.targetAddr
}

// tlsHost returns the host name to match against the peer's
// TLS certificate.
func (cm *connectMethod) tlsHost() string {
	h := cm.targetAddr
	if hasPort(h) {
		h = h[:strings.LastIndex(h, ":")]
	}
	return h
}

// persistConn wraps a connection, usually a persistent one
// (but may be used for non-keep-alive requests as well)
type persistConn struct {
	t        *Transport
	cacheKey string // its connectMethod.String()
	conn     net.Conn
	closed   bool                // whether conn has been closed
	br       *bufio.Reader       // from conn
	bw       *bufio.Writer       // to conn
	reqch    chan requestAndChan // written by roundTrip; read by readLoop
	writech  chan writeRequest   // written by roundTrip; read by writeLoop
	closech  chan struct{}       // broadcast close when readLoop (TCP connection) closes
	isProxy  bool

	// mutateHeaderFunc is an optional func to modify extra
	// headers on each outbound request before it's written. (the
	// original Request given to RoundTrip is not modified)
	mutateHeaderFunc func(Header)

	lk                   sync.Mutex // guards numExpectedResponses and broken
	numExpectedResponses int
	broken               bool // an error has happened on this connection; marked broken so it's not reused.
}

func (pc *persistConn) isBroken() bool {
	pc.lk.Lock()
	b := pc.broken
	pc.lk.Unlock()
	return b
}

var remoteSideClosedFunc func(error) bool // or nil to use default

func remoteSideClosed(err error) bool {
	if err == io.EOF {
		return true
	}
	if remoteSideClosedFunc != nil {
		return remoteSideClosedFunc(err)
	}
	return false
}

func (pc *persistConn) readLoop() {
	defer close(pc.closech)
	alive := true
	var lastbody io.ReadCloser // last response body, if any, read on this connection

	for alive {
		pb, err := pc.br.Peek(1)

		pc.lk.Lock()
		if pc.numExpectedResponses == 0 {
			pc.closeLocked()
			pc.lk.Unlock()
			if len(pb) > 0 {
				log.Printf("Unsolicited response received on idle HTTP channel starting with %q; err=%v",
					string(pb), err)
			}
			return
		}
		pc.lk.Unlock()

		rc := <-pc.reqch

		// Advance past the previous response's body, if the
		// caller hasn't done so.
		if lastbody != nil {
			lastbody.Close() // assumed idempotent
			lastbody = nil
		}

		var resp *Response
		if err == nil {
			resp, err = ReadResponse(pc.br, rc.req)
		}

		if err != nil {
			pc.close()
		} else {
			hasBody := rc.req.Method != "HEAD" && resp.ContentLength != 0
			if rc.addedGzip && hasBody && resp.Header.Get("Content-Encoding") == "gzip" {
				resp.Header.Del("Content-Encoding")
				resp.Header.Del("Content-Length")
				resp.ContentLength = -1
				gzReader, zerr := gzip.NewReader(resp.Body)
				if zerr != nil {
					pc.close()
					err = zerr
				} else {
					resp.Body = &readFirstCloseBoth{&discardOnCloseReadCloser{gzReader}, resp.Body}
				}
			}
			resp.Body = &bodyEOFSignal{body: resp.Body}
		}

		if err != nil || resp.Close || rc.req.Close {
			alive = false
		}

		hasBody := resp != nil && resp.ContentLength != 0
		var waitForBodyRead chan bool
		if hasBody {
			lastbody = resp.Body
			waitForBodyRead = make(chan bool, 1)
			resp.Body.(*bodyEOFSignal).fn = func() {
				if alive && !pc.t.putIdleConn(pc) {
					alive = false
				}
				if !alive || pc.isBroken() {
					pc.close()
				}
				waitForBodyRead <- true
			}
		}

		if alive && !hasBody {
			// When there's no response body, we immediately
			// reuse the TCP connection (putIdleConn), but
			// we need to prevent ClientConn.Read from
			// closing the Response.Body on the next
			// loop, otherwise it might close the body
			// before the client code has had a chance to
			// read it (even though it'll just be 0, EOF).
			lastbody = nil

			if !pc.t.putIdleConn(pc) {
				alive = false
			}
		}

		rc.ch <- responseAndError{resp, err}

		// Wait for the just-returned response body to be fully consumed
		// before we race and peek on the underlying bufio reader.
		if waitForBodyRead != nil {
			<-waitForBodyRead
		}

		if !alive {
			pc.close()
		}
	}
}

func (pc *persistConn) writeLoop() {
	for {
		select {
		case wr := <-pc.writech:
			if pc.isBroken() {
				wr.ch <- errors.New("http: can't write HTTP request on broken connection")
				continue
			}
			err := wr.req.Request.write(pc.bw, pc.isProxy, wr.req.extra)
			if err == nil {
				err = pc.bw.Flush()
			}
			if err != nil {
				pc.markBroken()
			}
			wr.ch <- err
		case <-pc.closech:
			return
		}
	}
}

type responseAndError struct {
	res *Response
	err error
}

type requestAndChan struct {
	req *Request
	ch  chan responseAndError

	// did the Transport (as opposed to the client code) add an
	// Accept-Encoding gzip header? only if it we set it do
	// we transparently decode the gzip.
	addedGzip bool
}

// A writeRequest is sent by the readLoop's goroutine to the
// writeLoop's goroutine to write a request while the read loop
// concurrently waits on both the write response and the server's
// reply.
type writeRequest struct {
	req *transportRequest
	ch  chan<- error
}

func (pc *persistConn) roundTrip(req *transportRequest) (resp *Response, err error) {
	if pc.mutateHeaderFunc != nil {
		pc.mutateHeaderFunc(req.extraHeaders())
	}

	// Ask for a compressed version if the caller didn't set their
	// own value for Accept-Encoding. We only attempted to
	// uncompress the gzip stream if we were the layer that
	// requested it.
	requestedGzip := false
	if !pc.t.DisableCompression && req.Header.Get("Accept-Encoding") == "" {
		// Request gzip only, not deflate. Deflate is ambiguous and 
		// not as universally supported anyway.
		// See: http://www.gzip.org/zlib/zlib_faq.html#faq38
		requestedGzip = true
		req.extraHeaders().Set("Accept-Encoding", "gzip")
	}

	pc.lk.Lock()
	pc.numExpectedResponses++
	pc.lk.Unlock()

	// Write the request concurrently with waiting for a response,
	// in case the server decides to reply before reading our full
	// request body.
	writeErrCh := make(chan error, 1)
	pc.writech <- writeRequest{req, writeErrCh}

	resc := make(chan responseAndError, 1)
	pc.reqch <- requestAndChan{req.Request, resc, requestedGzip}

	var re responseAndError
	var pconnDeadCh = pc.closech
	var failTicker <-chan time.Time
WaitResponse:
	for {
		select {
		case err := <-writeErrCh:
			if err != nil {
				re = responseAndError{nil, err}
				break WaitResponse
			}
		case <-pconnDeadCh:
			// The persist connection is dead. This shouldn't
			// usually happen (only with Connection: close responses
			// with no response bodies), but if it does happen it
			// means either a) the remote server hung up on us
			// prematurely, or b) the readLoop sent us a response &
			// closed its closech at roughly the same time, and we
			// selected this case first, in which case a response
			// might still be coming soon.
			//
			// We can't avoid the select race in b) by using a unbuffered
			// resc channel instead, because then goroutines can
			// leak if we exit due to other errors.
			pconnDeadCh = nil                               // avoid spinning
			failTicker = time.After(100 * time.Millisecond) // arbitrary time to wait for resc
		case <-failTicker:
			re = responseAndError{nil, errors.New("net/http: transport closed before response was received")}
			break WaitResponse
		case re = <-resc:
			break WaitResponse
		}
	}

	pc.lk.Lock()
	pc.numExpectedResponses--
	pc.lk.Unlock()

	return re.res, re.err
}

// markBroken marks a connection as broken (so it's not reused).
// It differs from close in that it doesn't close the underlying
// connection for use when it's still being read.
func (pc *persistConn) markBroken() {
	pc.lk.Lock()
	defer pc.lk.Unlock()
	pc.broken = true
}

func (pc *persistConn) close() {
	pc.lk.Lock()
	defer pc.lk.Unlock()
	pc.closeLocked()
}

func (pc *persistConn) closeLocked() {
	pc.broken = true
	if !pc.closed {
		pc.conn.Close()
		pc.closed = true
	}
	pc.mutateHeaderFunc = nil
}

var portMap = map[string]string{
	"http":  "80",
	"https": "443",
}

// canonicalAddr returns url.Host but always with a ":port" suffix
func canonicalAddr(url *url.URL) string {
	addr := url.Host
	if !hasPort(addr) {
		return addr + ":" + portMap[url.Scheme]
	}
	return addr
}

func responseIsKeepAlive(res *Response) bool {
	// TODO: implement.  for now just always shutting down the connection.
	return false
}

// bodyEOFSignal wraps a ReadCloser but runs fn (if non-nil) at most
// once, right before the final Read() or Close() call returns, but after
// EOF has been seen.
type bodyEOFSignal struct {
	body     io.ReadCloser
	fn       func()
	isClosed bool
	once     sync.Once
}

func (es *bodyEOFSignal) Read(p []byte) (n int, err error) {
	n, err = es.body.Read(p)
	if es.isClosed && n > 0 {
		panic("http: unexpected bodyEOFSignal Read after Close; see issue 1725")
	}
	if err == io.EOF {
		es.condfn()
	}
	return
}

func (es *bodyEOFSignal) Close() (err error) {
	if es.isClosed {
		return nil
	}
	es.isClosed = true
	err = es.body.Close()
	if err == nil {
		es.condfn()
	}
	return
}

func (es *bodyEOFSignal) condfn() {
	if es.fn != nil {
		es.once.Do(es.fn)
	}
}

type readFirstCloseBoth struct {
	io.ReadCloser
	io.Closer
}

func (r *readFirstCloseBoth) Close() error {
	if err := r.ReadCloser.Close(); err != nil {
		r.Closer.Close()
		return err
	}
	if err := r.Closer.Close(); err != nil {
		return err
	}
	return nil
}

// discardOnCloseReadCloser consumes all its input on Close.
type discardOnCloseReadCloser struct {
	io.ReadCloser
}

func (d *discardOnCloseReadCloser) Close() error {
	io.Copy(ioutil.Discard, d.ReadCloser) // ignore errors; likely invalid or already closed
	return d.ReadCloser.Close()
}
