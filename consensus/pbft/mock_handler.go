package pbft

import (
	"github.com/simplePBFT/common"
	"github.com/simplePBFT/consensus"
	"github.com/spf13/viper"
)

func obcBatchHelper(id uint64, config *viper.Viper, stack consensus.Stack) pbftConsumer {
	// It's not entirely obvious why the compiler likes the parent function, but not newObcBatch directly
	return newObcBatch(id, config, stack)
}
func NewHandler(id uint64,committee *common.Committee) pbftConsumer {
	tep := makeTestEndpoint(id, committee)

	ce := &consumerEndpoint{
		testEndpoint: tep,
	}

	ml := NewMockLedger(nil)
	ml.ce = ce

	cs := &completeStack{
		consumerEndpoint: ce,
		noopSecurity:     &noopSecurity{},
		MockLedger:       ml,
		skipTarget:       make(chan struct{}, 1),
	}
	ce.consumer = obcBatchHelper(id, loadConfig(), cs)
	ce.consumer.getPBFTCore().N = 4
	ce.consumer.getPBFTCore().f = 1
	return ce.consumer
}
