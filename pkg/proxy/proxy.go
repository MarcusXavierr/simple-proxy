package proxy

import (
	"bufio"
	"fmt"
	"net"
	"net/http"

	"github.com/MarcusXavierr/simple-proxy/pkg/tcp"

	"github.com/pkg/errors"
)

type Proxy struct {
	ClientReq    *http.Request
	clientConn   net.Conn
	ClientWriter *bufio.Writer
	ClientReader *bufio.Reader
	serverConn   net.Conn
}

// NewProxy creates a new Proxy instance from an incoming client connection.
// It reads an HTTP request from the client connection and initializes a Proxy struct
// with the necessary connection and request details.
//
// Parameters:
//   - clientConn: A network connection representing the client's connection
//
// Returns:
//   - A configured Proxy instance ready for request processing
//   - An error if the HTTP request cannot be parsed from the client connection
//
// Example:
//   proxy, err := NewProxy(clientConnection)
//   if err != nil {
//       // Handle error
//   }
func NewProxy(clientConn net.Conn) (*Proxy, error) {
	clientReader := bufio.NewReader(clientConn)
	req, err := http.ReadRequest(clientReader)

	if err != nil {
		return nil, errors.Wrap(err, "Parsing request from client")
	}

	return &Proxy{
		clientConn:   clientConn,
		ClientReq:    req,
		ClientWriter: bufio.NewWriter(clientConn),
		ClientReader: clientReader,
	}, nil
}

func (p *Proxy) StreamRequest() error {
	if err := p.connectToServer(); err != nil {
		return errors.Wrap(err, "Connecting to server")
	}
	defer p.serverConn.Close()

	// TODO: Add code to write the body of the request to server, so we can handle POST requests
	if _, err := p.serverConn.Write([]byte(p.mountRequestHeader())); err != nil {
		return errors.Wrap(err, "Writing request header")
	}

	proxyReader := bufio.NewReader(p.serverConn)

	// TODO: handle STREAMS and be able to read this on your unit test
	buffer := make([]byte, 1024)
	n, err := proxyReader.Read(buffer)
	if n > 0 {
		if _, writeErr := p.clientConn.Write(buffer[:n]); writeErr != nil {
			return errors.Wrap(writeErr, "Writing to client connection")
		}
	}
	if err != nil {
		return errors.Wrap(err, "Reading from server connection")
	}

	return nil
}

func (p *Proxy) ValidateImpementedMethods() error {
	if p.ClientReq.Method != "GET" {
		return errors.New("Only GET is implemented")
	}

	return nil
}

func (p *Proxy) connectToServer() error {
	host := p.ClientReq.Host
	if !hasPort(host) {
		host += ":80"
	}

	proxyfd, err := tcp.DialServer(host)
	if err != nil {
		return err
	}

	p.serverConn = proxyfd
	return nil
}

// Returns true if the host includes a valid port, false otherwise.
func hasPort(host string) bool {
	_, port, err := net.SplitHostPort(host)
	return err == nil && port != ""
}

func (p *Proxy) mountRequestHeader() string {
	req := p.ClientReq
	path := req.URL.Path
	if path == "" {
		path = "/"
	}

	return fmt.Sprintf("%s %s %s\r\nHost: %s\r\n\r\n", req.Method, path, req.Proto, req.Host)
}
