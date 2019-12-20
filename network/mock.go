package network

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/simplePBFT/protos"
	gp "google/protobuf"
	"time"
)

func (node Node) mockRequest() {
	time.Sleep(5e9)
	node.Consenter.RecvMsg(createTxMsg(1), &pb.PeerID{Name: fmt.Sprintf("vp%d", 0)})
}

func createTxMsg(tag int64) (msg *pb.Message) {
	tx := createTx(tag)
	txPacked := marshalTx(tx)
	msg = &pb.Message{
		Type:    pb.Message_CHAIN_TRANSACTION,
		Payload: txPacked,
	}
	return
}

func createTx(tag int64) (tx *pb.Transaction) {
	txTime := &gp.Timestamp{Seconds: tag, Nanos: 0}
	tx = &pb.Transaction{Type: pb.Transaction_CHAINCODE_DEPLOY,
		Timestamp: txTime,
		Payload:   []byte(fmt.Sprint(tag)),
	}
	return
}

func marshalTx(tx *pb.Transaction) (txPacked []byte) {
	txPacked, _ = proto.Marshal(tx)
	return
}
