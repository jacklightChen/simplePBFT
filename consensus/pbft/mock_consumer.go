package pbft

import (
	"github.com/simplePBFT/consensus"
	"github.com/simplePBFT/consensus/util/events"
	pb "github.com/simplePBFT/protos"
	"math/rand"
	"time"
)

type consumerEndpoint struct {
	*testEndpoint
	consumer     pbftConsumer
	execTxResult func([]*pb.Transaction) ([]byte, error)
}

func (ce *consumerEndpoint) stop() {
	ce.consumer.Close()
}

func (ce *consumerEndpoint) isBusy() bool {
	pbft := ce.consumer.getPBFTCore()
	if pbft.timerActive || pbft.skipInProgress || pbft.currentExec != nil {
		//ce.net.debugMsg("Reporting busy because of timer (%v) or skipInProgress (%v) or currentExec (%v)\n", pbft.timerActive, pbft.skipInProgress, pbft.currentExec)
		return true
	}

	select {
	case <-ce.consumer.idleChannel():
	default:
		//ce.net.debugMsg("Reporting busy because consumer not idle\n")
		return true
	}

	select {
	case ce.consumer.getManager().Queue() <- nil:
		//ce.net.debugMsg("Reporting busy because pbft not idle\n")
	default:
		return true
	}
	//}

	return false
}

func (ce *consumerEndpoint) deliver(msg []byte, senderHandle *pb.PeerID) {
	ce.consumer.RecvMsg(&pb.Message{Type: pb.Message_CONSENSUS, Payload: msg}, senderHandle)
}

type completeStack struct {
	*consumerEndpoint
	*noopSecurity
	*MockLedger
	mockPersist
	skipTarget chan struct{}
}

const MaxStateTransferTime int = 200

func (cs *completeStack) ValidateState()   {}
func (cs *completeStack) InvalidateState() {}
func (cs *completeStack) Start()           {}
func (cs *completeStack) Halt()            {}

func (cs *completeStack) UpdateState(tag interface{}, target *pb.BlockchainInfo, peers []*pb.PeerID) {
	select {
	// This guarantees the first SkipTo call is the one that's queued, whereas a mutex can be raced for
	case cs.skipTarget <- struct{}{}:
		go func() {
			// State transfer takes time, not simulating this hides bugs
			time.Sleep(time.Duration((MaxStateTransferTime/2)+rand.Intn(MaxStateTransferTime/2)) * time.Millisecond)
			cs.simulateStateTransfer(target, peers)
			cs.consumer.StateUpdated(tag, cs.GetBlockchainInfo())
			<-cs.skipTarget // Basically like releasing a mutex
		}()
	default:
		//cs.net.debugMsg("Ignoring skipTo because one is already in progress\n")
	}
}

type pbftConsumer interface {
	innerStack
	consensus.Consenter
	getPBFTCore() *pbftCore
	getManager() events.Manager // TODO, remove, this is a temporary measure
	Close()
	idleChannel() <-chan struct{}
}