package channeltalk

import (
	"encoding/json"
	"fmt"
	"time"
)

type nativeFunctionRequest struct {
	Method string `json:"method"`
	Params any    `json:"params"`
}

type nativeFunctionResponse struct {
	Result json.RawMessage              `json:"result"`
	Error  *nativeFunctionErrorResponse `json:"error"`
}

type nativeFunctionErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (r *nativeFunctionErrorResponse) Error() string {
	return fmt.Sprintf("client error: %s: %s", r.Type, r.Message)
}

type nativeFuntcionParams interface {
	Method() string
}

type errorResponse struct {
}

type issueTokenRequest struct {
	Secret    string  `json:"secret"`
	ChannelId *string `json:"channelId"`
}

func (r *issueTokenRequest) Method() string {
	return "issueToken"
}

type refreshIssueTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (r *refreshIssueTokenRequest) Method() string {
	return "issueToken"
}

type IssueTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiersIn"`

	Expiry time.Time
}

type WriteGroupMessageRequest struct {
	ChannelID     string       `json:"channelId"`
	GroupID       string       `json:"groupId"`
	RootMessageID string       `json:"rootMessageId"`
	Broadcast     bool         `json:"broadcast"`
	DTO           GroupMessage `json:"dto"`
}

func (r *WriteGroupMessageRequest) Method() string {
	return "writeGroupMessage"
}

type GroupMessage struct {
	Blocks    []MessageBlock `json:"blocks"`
	RequestId string         `json:"requestId"`
	BotName   string         `json:"botName"`
	Buttons   []Button       `json:"buttons,omitempty"`
}

type WriteGroupMessageResponse struct {
	Message struct {
		ID string `json:"id"`
	} `json:"message"`
}
