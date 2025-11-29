package slack

type WorkObject struct {
}

type WorkObjectOption func(*WorkObject)

func NewBlockWorkObject(blockID string, options ...WorkObjectOption) *WorkObject {
	o := &WorkObject{}

	for _, option := range options {
		option(o)
	}

	return o
}
