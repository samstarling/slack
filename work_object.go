package slack

type WorkObject struct {
}

type WorkObjectOption func(*WorkObject)

func NewBlockWorkObject(blockID string, options ...WorkObjectOption) *WorkObject {
	block := &WorkObject{}

	for _, option := range options {
		option(block)
	}

	return block
}
