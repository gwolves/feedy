package http

import "github.com/goccy/go-json"

type functionRequest struct {
	Method  string          `json:"method"`
	Params  functionParams  `json:"params"`
	Context functionContext `json:"context"`
}

type functionParams struct {
	Chat struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"chat"`

	Trigger struct {
		Type       string   `json:"type"`
		Attributes []string `json:"attributes"`
	} `json:"trigger"`

	Input    json.RawMessage `json:"input"`
	Language string          `json:"language"`
}

type functionContext struct {
	Channel struct {
		ID string `json:"id"`
	} `json:"channel"`
	Caller struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"caller"`
}

type functionResponse struct {
	Result any `json:"result"`
}

var succeedResponse = functionResponse{
	Result: map[string]any{
		"success": true,
	},
}

type subscribeInputs struct {
	Url     string `json:"url"`
	BotName string `json:"botname"`
}

type unsubscribeInputs struct {
	ID int64 `json:"id"`
}

type subscriptionsResponse struct {
	Subscriptions []subscription `json:"url"`
}

type subscription struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type autoCompleteUnsubscribeInputs struct {
	ID int64 `json:"id"`
}

type autoCompleteUnsubscribeResponse struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

const (
	methodString = "arst"
)
