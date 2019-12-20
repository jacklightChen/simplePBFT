package pbft

import (
	"fmt"
	pb "github.com/simplePBFT/protos"
)

type noopSecurity struct{}

func (ns *noopSecurity) Sign(msg []byte) ([]byte, error) {
	return nil, nil
}

func (ns *noopSecurity) Verify(peerID *pb.PeerID, signature []byte, message []byte) error {
	return nil
}

type mockPersist struct {
	store map[string][]byte
}

func (p *mockPersist) initialize() {
	if p.store == nil {
		p.store = make(map[string][]byte)
	}
}

func (p *mockPersist) ReadState(key string) ([]byte, error) {
	p.initialize()
	if val, ok := p.store[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("cannot find key %s", key)
}

func (p *mockPersist) ReadStateSet(prefix string) (map[string][]byte, error) {
	if p.store == nil {
		return nil, fmt.Errorf("no state yet")
	}
	ret := make(map[string][]byte)
	for k, v := range p.store {
		if len(k) >= len(prefix) && k[0:len(prefix)] == prefix {
			ret[k] = v
		}
	}
	return ret, nil
}

func (p *mockPersist) StoreState(key string, value []byte) error {
	p.initialize()
	p.store[key] = value
	return nil
}

func (p *mockPersist) DelState(key string) {
	p.initialize()
	delete(p.store, key)
}

