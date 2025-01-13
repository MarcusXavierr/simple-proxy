package main

import (
	"flag"
	"net"
	"path/filepath"
	"proxygo/pkg/proxy"
	"proxygo/pkg/tcp"
	"time"

	"go.uber.org/zap"
)

func main() {
	var addr string
	flag.StringVar(&addr, "h", ":1234", "port")
	flag.Parse()

	zapLogger, _ := configLogger().Build()
	defer zapLogger.Sync()
	zap.ReplaceGlobals(zapLogger)

	server, err := tcp.NewServer(addr)
	if err != nil {
		zap.S().Errorf("Error starting tcp listener %v", err)
		return
	}

	server.HandleServerLoop(handler)
}

func handler(conn *net.TCPConn) {
	defer conn.Close()
	px, err := proxy.NewProxy(conn)
	if err != nil {
		zap.S().Error(err)
		tcp.BadRequestError(conn, "")
		return
	}

	if err := px.ValidateImpementedMethods(); err != nil {
		zap.S().Error(err)
		tcp.ClientError(
			conn,
			px.ClientReq.Method,
			"501", "Not implemented",
			"The "+px.ClientReq.Method+" is not implemented. Yet!",
		)
		return
	}

	if err := px.StreamRequest(); err != nil {
		zap.S().Error(err)
		tcp.ServerError(conn, "")
		return
	}
}

func configLogger() zap.Config {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"stdout",
		filepath.Join("logs", time.Now().Format("2006-01-02")+".log"),
	}
	config.Encoding = "console"
	return config
}
