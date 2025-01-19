package proxy

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
)

func TestStreamRequest(t *testing.T) {
	// Create a pair of in-memory connections
	clientConn, clientEnd := net.Pipe()
	server, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer clientConn.Close()
	defer clientEnd.Close()
	defer server.Close()
	// Write on the mock server
	go func() {
		serverConn, err := server.Accept()
		if err != nil {
			t.Errorf("Failed to accept connection: %v", err)
			return
		}
		defer serverConn.Close()
		if _, err := serverConn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 11\r\n\r\nHello world")); err != nil {
			t.Errorf("Failed to write response: %v", err)
			return
		}
	}()

	ch := make(chan *Proxy)
	go func() {
		proxy, err := NewProxy(clientEnd)
		if err != nil {
			t.Errorf("Failed to create proxy: %v", err)
			ch <- nil
			return
		}
		ch <- proxy
	}()

	host := server.Addr().String()
	port := host[strings.IndexByte(host, ':')+1:]

	request := "GET http://localhost:" + port + "/ HTTP/1.1\r\nHost: localhost\r\n\r\n"
	if _, err := clientConn.Write([]byte(request)); err != nil {
		t.Fatalf("Failed to write request: %v", err)
	}
	// clientConn.Close()

	proxy := <-ch
	if proxy == nil {
		t.Errorf("Failed to create proxy")
	}

	respCh := make(chan *http.Response)
	go func() {
		resp, err := http.ReadResponse(bufio.NewReader(clientConn), proxy.ClientReq)
		if err != nil {
			t.Errorf("Response read error: %v", err)
			respCh <- nil
			return
		}

		respCh <- resp
	}()

	// Testing buffered writer
	// bufio.NewWriter(clientEnd).Write([]byte("Hello"))

	// Perform the stream request
	err = proxy.StreamRequest()
	if err != nil {
		t.Errorf("StreamRequest failed: %v", err)
	}

	resp := <-respCh
	if resp == nil {
		t.Errorf("Failed to read response")
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: got %v, want %v", resp.StatusCode, http.StatusOK)
	}

	// Check the response body
	body := make([]byte, 11)
	_, err = resp.Body.Read(body)
	if err != nil && err != io.EOF {
		t.Errorf("Failed to read response body: %v", err)
	}
	if string(body) != "Hello world" {
		t.Errorf("Unexpected response body: got %v, want %v", string(body), "Hello world")
	}
}
