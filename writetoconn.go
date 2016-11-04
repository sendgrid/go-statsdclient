package statsdclient

import "net"

type writeToConn struct {
	remoteAddr *net.UDPAddr

	udpConn *net.UDPConn
}

func newWriteToConn(raddr string) (*writeToConn, error) {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, err
	}

	remoteAddr, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		return nil, err
	}

	conn := &writeToConn{
		udpConn:    udpConn,
		remoteAddr: remoteAddr,
	}
	return conn, nil
}

func (w *writeToConn) Write(p []byte) (int, error) {
	n, err := w.udpConn.WriteToUDP(p, w.remoteAddr)
	return n, err
}

func (w *writeToConn) Close() error {
	return w.udpConn.Close()
}
