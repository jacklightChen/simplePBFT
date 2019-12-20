package common

import (
	. "github.com/simplePBFT/config"
	. "github.com/simplePBFT/protocol"
	"net"
)

type Peer struct {
	ID string
	url URL
	conn       net.Conn // source connection
}

func (peer *Peer) Send(msg []byte) bool {
	if peer.conn != nil {
		msg := Packet(msg)
		n,err := peer.conn.Write(msg)
		if err != nil {
			return false
		}
		if n == len(msg) {
			return true
		} else {
			return true
		}
	}
	return false

}

