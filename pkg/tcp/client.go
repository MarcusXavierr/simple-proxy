package tcp

import (
	"net"

	"github.com/pkg/errors"
)

func DialServer(host string) (*net.TCPConn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return nil, errors.Wrap(err, "Getting TCP Address")
	}

	return net.DialTCP("tcp", nil, raddr)
}
