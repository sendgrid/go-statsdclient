package statsdclient

import (
	"net"
)

type writeToConn struct {
	packetListener net.PacketConn
	remoteAddr     *net.UDPAddr
}

func newWriteToConn(raddr string) (*writeToConn, error) {
	// need to create the local socket for sending messages from
	packetListener, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	// resolve the udp address provided
	remoteAddr, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		return nil, err
	}

	return &writeToConn{packetListener: packetListener, remoteAddr: remoteAddr}, nil
}

func (w *writeToConn) Write(p []byte) (int, error) {
	return w.packetListener.WriteTo(p, w.remoteAddr)
}

func (w *writeToConn) Close() error {
	return w.packetListener.Close()
}
