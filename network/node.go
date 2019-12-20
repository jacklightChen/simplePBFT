package network

import (
	"fmt"
	"github.com/simplePBFT/config"
	"github.com/simplePBFT/consensus"
	"github.com/simplePBFT/consensus/pbft"
	. "github.com/simplePBFT/protocol"
	pb "github.com/simplePBFT/protos"
	"strconv"
	"time"
)

type Node struct {
	Config      *config.Config
	NodeID      string
	View        *View
	MsgBuffer   *MsgBuffer
	MsgEntrance chan interface{}
	MsgDelivery chan interface{}
	Alarm       chan bool
	Srv         *Server
	Consenter   consensus.Consenter
}

type MsgBuffer struct {
	RecvMsgs []*pb.Message
}

type View struct {
	ID      int64
	Primary string
}

const ResolvingTimeDuration = time.Millisecond * 1000 // 1 second.

func NewNode(server *Server) *Node {
	const viewID = 10000000000 // temporary.

	config := server.config

	node := &Node{
		// Hard-coded for test.
		NodeID: config.LocalNodeID,
		View: &View{
			ID:      viewID,
			Primary: server.localCmt.NodeInfos[0].NodeID, // first node is the primary
		},

		// Channels
		MsgEntrance: make(chan interface{}),
		MsgDelivery: make(chan interface{}),
		Alarm:       make(chan bool),
		Config:      config,
		Srv:         server,
	}
	id, _ := strconv.Atoi(node.NodeID[len(node.NodeID)-1:len(node.NodeID)])
	node.Consenter = pbft.NewHandler(uint64(id), node.Srv.localCmt)
	// if node is primary,mock request
	// should wait for all nodes in committee start
	//if node.NodeID == server.localCmt.NodeInfos[0].NodeID {
	//	go func(node *Node){
	//		time.Sleep(5e9)
	//		var msg consensus.RequestMsg
	//		msg.ClientID="123"
	//		msg.Operation="ssfaafs"
	//		node.MsgEntrance <- &msg
	//	}(node)
	//}

	// Start message dispatcher
	go node.dispatchMsg()

	// if node is primary,mock request
	// should wait for all nodes in committee start
	if node.NodeID == server.localCmt.NodeInfos[0].NodeID {
		go node.mockRequest()
	}
	// Start alarm trigger
	//go node.alarmToDispatcher()

	// Start message resolver
	//go node.resolveMsg()

	return node
}

func (node *Node) broadcast(msg pb.Message) chan bool {
	//msg.SetCommitteeID(node.Config.LocalCommittee)
	jsonMsg, err := Serial(msg)
	if err == nil {
		return node.Srv.localCmt.Broadcast(jsonMsg)
	}
	return nil
}

func (node *Node) Broadcast(jsonBytes []byte) chan bool {
	return node.Srv.localCmt.Broadcast(jsonBytes)
}

func (node *Node) dispatchMsg() {
	for {
		select {
		case msg := <-node.MsgEntrance:
			tmp :=msg.(pb.Message)
			fmt.Sprintf("receive message from replicaId:%d\n",tmp.ReplicaId)
			//id, _ := strconv.Atoi(node.NodeID[len(node.NodeID)-1:len(node.NodeID)])
			node.Consenter.RecvMsg(&tmp, &pb.PeerID{Name: fmt.Sprintf("vp%d", tmp.ReplicaId)})


			//err := node.routeMsg(msg)
			//if err != nil {
			//	fmt.Println(err)
			//	//TODO: send err to ErrorChannel
			//}
			//case <-node.Alarm:
			//	err := node.routeMsgWhenAlarmed()
			//	if err != nil {
			//		fmt.Println(err)
			//		// TODO: send err to ErrorChannel
			//	}
		}
	}
}

//func (node *Node) routeMsg(msg interface{}) []error {
//	switch msg.(type) {
//	case consensus.RequestMsg:
//		if node.CurrentState == nil {
//			// Copy buffered messages first.
//			msgs := make([]*consensus.RequestMsg, len(node.MsgBuffer.ReqMsgs))
//			copy(msgs, node.MsgBuffer.ReqMsgs)
//			temp := msg.(consensus.RequestMsg)
//			// Append a newly arrived message.
//			msgs = append(msgs, &temp)
//
//			// Empty the buffer.
//			node.MsgBuffer.ReqMsgs = make([]*consensus.RequestMsg, 0)
//
//			// Send messages.
//			node.MsgDelivery <- msgs
//		} else {
//			node.MsgBuffer.ReqMsgs = append(node.MsgBuffer.ReqMsgs, msg.(*consensus.RequestMsg))
//		}
//	}
//
//	return nil
//}

//func (node *Node) routeMsgWhenAlarmed() []error {
//	if node.CurrentState == nil {
//		// Check ReqMsgs, send them.
//		if len(node.MsgBuffer.ReqMsgs) != 0 {
//			msgs := make([]*consensus.RequestMsg, len(node.MsgBuffer.ReqMsgs))
//			copy(msgs, node.MsgBuffer.ReqMsgs)
//
//			node.MsgDelivery <- msgs
//		}
//
//		// Check PrePrepareMsgs, send them.
//		if len(node.MsgBuffer.PrePrepareMsgs) != 0 {
//			msgs := make([]*consensus.PrePrepareMsg, len(node.MsgBuffer.PrePrepareMsgs))
//			copy(msgs, node.MsgBuffer.PrePrepareMsgs)
//
//			node.MsgDelivery <- msgs
//		}
//	} else {
//		switch node.CurrentState.CurrentStage {
//		case consensus.PrePrepared:
//			// Check PrepareMsgs, send them.
//			if len(node.MsgBuffer.PrepareMsgs) != 0 {
//				msgs := make([]*consensus.VoteMsg, len(node.MsgBuffer.PrepareMsgs))
//				copy(msgs, node.MsgBuffer.PrepareMsgs)
//
//				node.MsgDelivery <- msgs
//			}
//		case consensus.Prepared:
//			// Check CommitMsgs, send them.
//			if len(node.MsgBuffer.CommitMsgs) != 0 {
//				msgs := make([]*consensus.VoteMsg, len(node.MsgBuffer.CommitMsgs))
//				copy(msgs, node.MsgBuffer.CommitMsgs)
//
//				node.MsgDelivery <- msgs
//			}
//		}
//	}
//
//	return nil
//}

//func (node *Node) resolveMsg() {
//	for {
//		// Get buffered messages from the dispatcher.
//		msgs := <-node.MsgDelivery
//		switch msgs.(type) {
//		case []*consensus.RequestMsg:
//			errs := node.resolveRequestMsg(msgs.([]*consensus.RequestMsg))
//			if len(errs) != 0 {
//				for _, err := range errs {
//					fmt.Println(err)
//				}
//				// TODO: send err to ErrorChannel
//			}
//		case []*consensus.PrePrepareMsg:
//			errs := node.resolvePrePrepareMsg(msgs.([]*consensus.PrePrepareMsg))
//			if len(errs) != 0 {
//				for _, err := range errs {
//					fmt.Println(err)
//				}
//				// TODO: send err to ErrorChannel
//			}
//		case []*consensus.VoteMsg:
//			voteMsgs := msgs.([]*consensus.VoteMsg)
//			if len(voteMsgs) == 0 {
//				break
//			}
//
//			if voteMsgs[0].VoteType == consensus.PrepareMsg {
//				errs := node.resolvePrepareMsg(voteMsgs)
//				if len(errs) != 0 {
//					for _, err := range errs {
//						fmt.Println(err)
//					}
//					// TODO: send err to ErrorChannel
//				}
//			} else if voteMsgs[0].VoteType == consensus.CommitMsg {
//				errs := node.resolveCommitMsg(voteMsgs)
//				if len(errs) != 0 {
//					for _, err := range errs {
//						fmt.Println(err)
//					}
//					// TODO: send err to ErrorChannel
//				}
//			}
//		}
//	}
//}
//
//func (node *Node) alarmToDispatcher() {
//	for {
//		time.Sleep(ResolvingTimeDuration)
//		node.Alarm <- true
//	}
//}
