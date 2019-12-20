package config

type URL string

type Config struct {
	CommitteeConfigs []CommitteeConfig `json:"committees"`
	LocalNodeID      string            `json:"node_id"`
	LocalURL         URL               `json:"local_url"`
	LocalCommittee   string            `json:"local_committee"`
	OpeURL           URL               `json:"ope_url"`
}

type CommitteeConfig struct {
	ID    string       `json:"id"`
	Nodes []NodeConfig `json:"nodes"`
}

type NodeConfig struct {
	URL    URL    `json:"url"`
	NodeID string `json:"node_id"`
}
