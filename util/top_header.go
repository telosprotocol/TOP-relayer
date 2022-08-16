package util

type BlockList struct {
	Hash  string `json:"blockHash"`
	Index string `json:"blockIndex"`
}
type TopHeader struct {
	Number      string      `json:"number"`
	Hash        string      `json:"hash"`
	Header      string      `json:"header"`
	BlockType   string      `json:"blockType"`
	ChainBits   string      `json:"chainBits"`
	RelatedList []BlockList `json:"blockList"`
}
