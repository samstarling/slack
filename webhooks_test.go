package slack

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostWebhook_OK(t *testing.T) {
	s := startServer()
	defer s.Close()

	var receivedPayload WebhookMessage

	s.RegisterHandler("/webhook", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&receivedPayload)
		if err != nil {
			t.Errorf("Request contained invalid JSON, %s", err)
		}

		response := []byte(`{}`)
		rw.Write(response)
	})

	url := "http://" + serverAddr + "/webhook"

	payload := &WebhookMessage{
		Text: "Test Text",
		Attachments: []Attachment{
			{
				Text: "Foo",
			},
		},
	}

	err := PostWebhook(url, payload)

	if err != nil {
		t.Errorf("Expected not to receive error: %s", err)
	}

	if !reflect.DeepEqual(payload, &receivedPayload) {
		t.Errorf("Payload did not match\nwant: %#v\n got: %#v", payload, receivedPayload)
	}
}

func TestPostWebhook_NotOK(t *testing.T) {
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/webhook2", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("500 - Something bad happened!"))
	})

	url := "http://" + serverAddr + "/webhook2"

	err := PostWebhook(url, &WebhookMessage{})

	if err == nil {
		t.Errorf("Expected to receive error")
	}
}

func TestPostWebhook_MessageLimitExceeded(t *testing.T) {
	s := startServer()
	defer s.Close()

	s.RegisterHandler("/message_limit_exceeded", func(rw http.ResponseWriter, r *http.Request) {
		// When a workspace's message limit is exceeded we get a 429 without a Retry-After header
		rw.WriteHeader(http.StatusTooManyRequests)
		rw.Write([]byte("message_limit_exceeded"))
	})

	url := "http://" + serverAddr + "/message_limit_exceeded"

	err := PostWebhook(url, &WebhookMessage{})

	if err == nil {
		t.Errorf("Expected to receive error")
	}
	assert.IsType(t, StatusCodeError{}, err)
}

func TestWebhookMessage_WithBlocks(t *testing.T) {
	textBlockObject := NewTextBlockObject("plain_text", "text", false, false)
	sectionBlock := NewSectionBlock(textBlockObject, nil, nil)

	singleBlock := &Blocks{BlockSet: []Block{sectionBlock}}
	twoBlocks := &Blocks{BlockSet: []Block{sectionBlock, sectionBlock}}

	msgSingleBlock := WebhookMessage{Blocks: singleBlock}
	assert.Equal(t, 1, len(msgSingleBlock.Blocks.BlockSet))

	msgJsonSingleBlock, _ := json.Marshal(msgSingleBlock)
	assert.Equal(t, `{"blocks":[{"type":"section","text":{"type":"plain_text","text":"text","emoji":false}}],"replace_original":false,"delete_original":false}`, string(msgJsonSingleBlock))

	msgTwoBlocks := WebhookMessage{Blocks: twoBlocks}
	assert.Equal(t, 2, len(msgTwoBlocks.Blocks.BlockSet))

	msgJsonTwoBlocks, _ := json.Marshal(msgTwoBlocks)
	assert.Equal(t, `{"blocks":[{"type":"section","text":{"type":"plain_text","text":"text","emoji":false}},{"type":"section","text":{"type":"plain_text","text":"text","emoji":false}}],"replace_original":false,"delete_original":false}`, string(msgJsonTwoBlocks))

	msgNoBlocks := WebhookMessage{Text: "foo"}
	msgJsonNoBlocks, _ := json.Marshal(msgNoBlocks)
	assert.Equal(t, `{"text":"foo","replace_original":false,"delete_original":false}`, string(msgJsonNoBlocks))
}
