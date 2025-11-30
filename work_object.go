package slack

type ExternalRef struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type WorkObjectEntityType string

const WorkObjectEntityTypeItem WorkObjectEntityType = "slack#/entities/item"

type WorkObjectEntity struct {
	AppUnfurlURL  string                  `json:"app_unfurl_url"`
	URL           string                  `json:"url"`
	ExternalRef   ExternalRef             `json:"external_ref"`
	EntityType    WorkObjectEntityType    `json:"entity_type"`
	EntityPayload WorkObjectEntityPayload `json:"entity_payload"`
}

type WorkObjectEntityAttributes struct {
	Title struct {
		Text string `json:"text"`
	} `json:"title"`
}

type WorkObjectEntityPayload struct {
	Attributes WorkObjectEntityAttributes `json:"attributes,omitempty"`
	Fields     map[string]any             `json:"fields,omitempty"`
}

type WorkObjectOption func(*WorkObjectEntity)

func NewBlockWorkObject(blockID string, options ...WorkObjectOption) *WorkObjectEntity {
	o := &WorkObjectEntity{}

	for _, option := range options {
		option(o)
	}

	return o
}
