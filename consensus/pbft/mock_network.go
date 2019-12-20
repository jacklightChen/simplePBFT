package pbft

import (
	"encoding/json"
	"fmt"
	"github.com/simplePBFT/common"
	pb "github.com/simplePBFT/protos"
	"strconv"
)

type endpoint interface {
	stop()
	deliver([]byte, *pb.PeerID)
	getHandle() *pb.PeerID
	getID() uint64
	isBusy() bool
}

type taggedMsg struct {
	src int
	dst int
	msg []byte
}

type testnet struct {
	debug     bool
	N         int
	closed    chan struct{}
	endpoints []endpoint
	msgs      chan taggedMsg
	filterFn  func(int, int, []byte) []byte
}

type testEndpoint struct {
	id  uint64
	committee *common.Committee
}

func makeTestEndpoint(id uint64, committee *common.Committee) *testEndpoint {
	ep := &testEndpoint{}
	ep.id = id
	ep.committee = committee
	return ep
}

func (ep *testEndpoint) getID() uint64 {
	return ep.id
}

func (ep *testEndpoint) getHandle() *pb.PeerID {
	return &pb.PeerID{Name: fmt.Sprintf("vp%d", ep.id)}
}

func (ep *testEndpoint) GetNetworkInfo() (self *pb.PeerEndpoint, network []*pb.PeerEndpoint, err error) {
	oSelf, oNetwork, _ := ep.GetNetworkHandles()
	self = &pb.PeerEndpoint{
		ID:   oSelf,
		Type: pb.PeerEndpoint_VALIDATOR,
	}

	network = make([]*pb.PeerEndpoint, len(oNetwork))
	for i, id := range oNetwork {
		network[i] = &pb.PeerEndpoint{
			ID:   id,
			Type: pb.PeerEndpoint_VALIDATOR,
		}
	}
	return
}

func (ep *testEndpoint) GetNetworkHandles() (self *pb.PeerID, network []*pb.PeerID, err error) {
	if nil == ep.committee {
			err = fmt.Errorf("Network not initialized")
			return
		}
		self = ep.getHandle()
		network = make([]*pb.PeerID, len(ep.committee.NodeInfos))
		for i, oep := range ep.committee.NodeInfos {
				// In case this is invoked before all endpoints are initialized, this emulates a real network as well
				id, _ := strconv.Atoi(oep.NodeID)
				network[i] = &pb.PeerID{Name: fmt.Sprintf("vp%d",uint64(id))}
		}
		return
	return
}

// Broadcast delivers to all endpoints.  In contrast to the stack
// Broadcast, this will also deliver back to the replica.  We keep
// this behavior, because it exposes subtle bugs in the
// implementation.
func (ep *testEndpoint) Broadcast(msg *pb.Message, peerType pb.PeerEndpoint_Type) error {
	jsonbytes ,_ := json.Marshal(msg)
	ep.committee.Broadcast(jsonbytes)
	return nil
}

func (ep *testEndpoint) Unicast(msg *pb.Message, receiverHandle *pb.PeerID) error {
	//receiverID, err := getValidatorID(receiverHandle)
	//if err != nil {
	//	return fmt.Errorf("Couldn't unicast message to %s: %v", receiverHandle.Name, err)
	//}
	//internalQueueMessage(ep.net.msgs, taggedMsg{int(ep.id), int(receiverID), msg.Payload})
	return nil
}