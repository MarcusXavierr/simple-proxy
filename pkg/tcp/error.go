package tcp

import (
	"fmt"
	"io"
)

func ClientError(conn io.Writer, cause, errnum, shortmsg, longmsg string) {
	// Header
	_, _ = conn.Write([]byte(fmt.Sprintf("HTTP/1.0 %s %s\r\n", errnum, shortmsg)))
	_, _ = conn.Write([]byte("Content-Type: text/html \r\n\r\n"))

	// Body
	_, _ = conn.Write([]byte(fmt.Sprintf("<html><title>%s</title>", shortmsg)))
	_, _ = conn.Write([]byte("<body>\r\n"))
	_, _ = conn.Write([]byte(fmt.Sprintf("<h1>%s: %s</h1>\r\n", errnum, shortmsg)))
	_, _ = conn.Write([]byte(fmt.Sprintf("<p>%s: %s</p>\r\n", longmsg, cause)))
	_, _ = conn.Write([]byte("<hr><em>The webserver<em></body></html>\r\n"))
}

func BadRequestError(conn io.Writer, cause string) {
	if cause == "" {
		cause = "Client Error"
	}

	ClientError(conn, cause, "400", "Bad Request", "The request could not be processed ")
}

func ServerError(conn io.Writer, cause string) {
	if cause == "" {
		cause = "Internal Error"
	}

	ClientError(conn, cause, "500", "Server Error", "The server could not process this request")
}
