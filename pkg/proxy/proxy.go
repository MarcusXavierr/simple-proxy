package proxy

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"proxygo/pkg/tcp"
	"regexp"

	"github.com/pkg/errors"
)

type Proxy struct {
	ClientReq    *http.Request
	clientConn   *net.TCPConn
	clientWriter *bufio.Writer
	clientReader *bufio.Reader
	serverConn   *net.TCPConn
}

func NewProxy(clientConn *net.TCPConn) (*Proxy, error) {
	clientReader := bufio.NewReader(clientConn)
	req, err := http.ReadRequest(clientReader)

	if err != nil {
		return nil, errors.Wrap(err, "Parsing request from client")
	}

	return &Proxy{
		clientConn:   clientConn,
		ClientReq:    req,
		clientWriter: bufio.NewWriter(clientConn),
		clientReader: clientReader,
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
	if _, err := p.clientWriter.ReadFrom(proxyReader); err != nil {
		return errors.Wrap(err, "Streaming response from server to client")
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
	hasPortSuffix, _ := regexp.MatchString("^[a-zA-Z0-9]+\\.[a-zA-Z0-9]+:[0-9]+$", host)
	if !hasPortSuffix {
		host += ":80"
	}

	proxyfd, err := tcp.DialServer(host)
	if err != nil {
		return err
	}

	p.serverConn = proxyfd
	return nil
}

func (p *Proxy) mountRequestHeader() string {
	req := p.ClientReq
	path := req.URL.Path
	if path == "" {
		path = "/"
	}

	return fmt.Sprintf("%s %s %s\r\nHost: %s\r\n\r\n", req.Method, path, req.Proto, req.Host)
}
