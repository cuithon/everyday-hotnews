// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
)

var serverAddr string
var once sync.Once

func echoServer(ws *Conn) { io.Copy(ws, ws) }

type Count struct {
	S string
	N int
}

func countServer(ws *Conn) {
	for {
		var count Count
		err := JSON.Receive(ws, &count)
		if err != nil {
			return
		}
		count.N++
		count.S = strings.Repeat(count.S, count.N)
		err = JSON.Send(ws, count)
		if err != nil {
			return
		}
	}
}

func startServer() {
	http.Handle("/echo", Handler(echoServer))
	http.Handle("/count", Handler(countServer))
	server := httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
	log.Print("Test WebSocket server listening on ", serverAddr)
}

func newConfig(t *testing.T, path string) *Config {
	config, _ := NewConfig(fmt.Sprintf("ws://%s%s", serverAddr, path), "http://localhost")
	return config
}

func TestEcho(t *testing.T) {
	once.Do(startServer)

	// websocket.Dial()
	client, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatal("dialing", err)
	}
	conn, err := NewClient(newConfig(t, "/echo"), client)
	if err != nil {
		t.Errorf("WebSocket handshake error: %v", err)
		return
	}

	msg := []byte("hello, world\n")
	if _, err := conn.Write(msg); err != nil {
		t.Errorf("Write: %v", err)
	}
	var actual_msg = make([]byte, 512)
	n, err := conn.Read(actual_msg)
	if err != nil {
		t.Errorf("Read: %v", err)
	}
	actual_msg = actual_msg[0:n]
	if !bytes.Equal(msg, actual_msg) {
		t.Errorf("Echo: expected %q got %q", msg, actual_msg)
	}
	conn.Close()
}

func TestAddr(t *testing.T) {
	once.Do(startServer)

	// websocket.Dial()
	client, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatal("dialing", err)
	}
	conn, err := NewClient(newConfig(t, "/echo"), client)
	if err != nil {
		t.Errorf("WebSocket handshake error: %v", err)
		return
	}

	ra := conn.RemoteAddr().String()
	if !strings.HasPrefix(ra, "ws://") || !strings.HasSuffix(ra, "/echo") {
		t.Errorf("Bad remote addr: %v", ra)
	}
	la := conn.LocalAddr().String()
	if !strings.HasPrefix(la, "http://") {
		t.Errorf("Bad local addr: %v", la)
	}
	conn.Close()
}

func TestCount(t *testing.T) {
	once.Do(startServer)

	// websocket.Dial()
	client, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatal("dialing", err)
	}
	conn, err := NewClient(newConfig(t, "/count"), client)
	if err != nil {
		t.Errorf("WebSocket handshake error: %v", err)
		return
	}

	var count Count
	count.S = "hello"
	if err := JSON.Send(conn, count); err != nil {
		t.Errorf("Write: %v", err)
	}
	if err := JSON.Receive(conn, &count); err != nil {
		t.Errorf("Read: %v", err)
	}
	if count.N != 1 {
		t.Errorf("count: expected %d got %d", 1, count.N)
	}
	if count.S != "hello" {
		t.Errorf("count: expected %q got %q", "hello", count.S)
	}
	if err := JSON.Send(conn, count); err != nil {
		t.Errorf("Write: %v", err)
	}
	if err := JSON.Receive(conn, &count); err != nil {
		t.Errorf("Read: %v", err)
	}
	if count.N != 2 {
		t.Errorf("count: expected %d got %d", 2, count.N)
	}
	if count.S != "hellohello" {
		t.Errorf("count: expected %q got %q", "hellohello", count.S)
	}
	conn.Close()
}

func TestWithQuery(t *testing.T) {
	once.Do(startServer)

	client, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatal("dialing", err)
	}

	config := newConfig(t, "/echo")
	config.Location, err = url.ParseRequest(fmt.Sprintf("ws://%s/echo?q=v", serverAddr))
	if err != nil {
		t.Fatal("location url", err)
	}

	ws, err := NewClient(config, client)
	if err != nil {
		t.Errorf("WebSocket handshake: %v", err)
		return
	}
	ws.Close()
}

func TestWithProtocol(t *testing.T) {
	once.Do(startServer)

	client, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatal("dialing", err)
	}

	config := newConfig(t, "/echo")
	config.Protocol = append(config.Protocol, "test")

	ws, err := NewClient(config, client)
	if err != nil {
		t.Errorf("WebSocket handshake: %v", err)
		return
	}
	ws.Close()
}

func TestHTTP(t *testing.T) {
	once.Do(startServer)

	// If the client did not send a handshake that matches the protocol
	// specification, the server MUST return an HTTP respose with an
	// appropriate error code (such as 400 Bad Request)
	resp, err := http.Get(fmt.Sprintf("http://%s/echo", serverAddr))
	if err != nil {
		t.Errorf("Get: error %#v", err)
		return
	}
	if resp == nil {
		t.Error("Get: resp is null")
		return
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Get: expected %q got %q", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestTrailingSpaces(t *testing.T) {
	// http://code.google.com/p/go/issues/detail?id=955
	// The last runs of this create keys with trailing spaces that should not be
	// generated by the client.
	once.Do(startServer)
	config := newConfig(t, "/echo")
	for i := 0; i < 30; i++ {
		// body
		ws, err := DialConfig(config)
		if err != nil {
			t.Errorf("Dial #%d failed: %v", i, err)
			break
		}
		ws.Close()
	}
}

func TestSmallBuffer(t *testing.T) {
	// http://code.google.com/p/go/issues/detail?id=1145
	// Read should be able to handle reading a fragment of a frame.
	once.Do(startServer)

	// websocket.Dial()
	client, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatal("dialing", err)
	}
	conn, err := NewClient(newConfig(t, "/echo"), client)
	if err != nil {
		t.Errorf("WebSocket handshake error: %v", err)
		return
	}

	msg := []byte("hello, world\n")
	if _, err := conn.Write(msg); err != nil {
		t.Errorf("Write: %v", err)
	}
	var small_msg = make([]byte, 8)
	n, err := conn.Read(small_msg)
	if err != nil {
		t.Errorf("Read: %v", err)
	}
	if !bytes.Equal(msg[:len(small_msg)], small_msg) {
		t.Errorf("Echo: expected %q got %q", msg[:len(small_msg)], small_msg)
	}
	var second_msg = make([]byte, len(msg))
	n, err = conn.Read(second_msg)
	if err != nil {
		t.Errorf("Read: %v", err)
	}
	second_msg = second_msg[0:n]
	if !bytes.Equal(msg[len(small_msg):], second_msg) {
		t.Errorf("Echo: expected %q got %q", msg[len(small_msg):], second_msg)
	}
	conn.Close()
}
