package system

import (
	"fmt"
	"net"
	"os"
)

func ExtractListenerFile(ln net.Listener) (*os.File, error) {
	tcpListener, ok := ln.(*net.TCPListener)
	if !ok {
		return nil, fmt.Errorf("unsupported listener: %T", ln)
	}
	return tcpListener.File()
}
