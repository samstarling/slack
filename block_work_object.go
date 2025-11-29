package slack

type BlockWorkObject struct {
	BlockID string `json:"block_id,omitempty"`
}

type WorkObjectBlockOption func(*BlockWorkObject)

func NewBlockWorkObject(blockID string, options ...WorkObjectBlockOption) *BlockWorkObject {
	block := &BlockWorkObject{BlockID: blockID}

	for _, option := range options {
		option(block)
	}

	return block
}

func (b BlockWorkObject) BlockType() MessageBlockType {
	return MBTWorkObject
}

func (b BlockWorkObject) ID() string {
	return b.BlockID
}
