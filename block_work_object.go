package slack

type BlockWorkObject struct {
	BlockID string `json:"block_id,omitempty"`
}

func (b BlockWorkObject) BlockType() MessageBlockType {
	return MBTWorkObject
}

func (b BlockWorkObject) ID() string {
	return b.BlockID
}
