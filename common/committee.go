package common

import (
	"fmt"
	. "github.com/simplePBFT/config"
	"net"
	"sync"
	"time"
)

type NodeInfo struct {
	NodeID  string
	NodeURL URL
}

/*
a committee is a sharding, which manages message sending
*/
type Committee struct {
	mutex       sync.Mutex
	CommitteeID string
	NodeInfos   []NodeInfo
	peers       []*Peer
	peerMap     map[string]*Peer
}

func (cmt *Committee) Broadcast(msgBytes []byte) chan bool {
	fmt.Println("broadcast")
	successChan := make(chan bool, len(cmt.peers))
	for _, peer := range cmt.peers {
		go func(peer *Peer) {
			fmt.Println("send to peer:", peer.ID)
			success := peer.Send(msgBytes)
			fmt.Println("result:", success)
			successChan <- success
		}(peer)
	}
	return successChan
}

func (cmt *Committee) Send(id string, msgBytes []byte) bool {
	if peer, ok := cmt.peerMap[id]; ok {
		return peer.Send(msgBytes)
	}
	return false
}

func (cmt *Committee) Dails(localId string) {
	cmt.peers = make([]*Peer, 0)
	cmt.peerMap = make(map[string]*Peer)
	dail := func(node NodeInfo) {
		var count uint32 = 1
		d := time.Duration(2 * time.Second)
		for {
			conn, err := net.Dial("tcp", string(node.NodeURL))
			if err == nil && conn != nil {
				cmt.mutex.Lock()
				peer := &Peer{node.NodeID, node.NodeURL, conn}
				cmt.peers = append(cmt.peers, peer)
				cmt.peerMap[node.NodeID] = peer
				cmt.mutex.Unlock()
				fmt.Printf("connect to %v successfully\n", node.NodeURL)
				return
			}
			//fmt.Printf("dail url:%v falied: %v\n", node.NodeURL, err )
			count++
			if count%100 == 0 {
				d = d + time.Duration(time.Second)
			}
			time.Sleep(d)
		}
	}
	for _, node := range cmt.NodeInfos {
		if node.NodeID != localId {
			go dail(node)
		}
	}
}
