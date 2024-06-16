package channeltalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Message struct {
	Blocks []MessageBlock `json:"blocks"`
}

type BlockType string

const (
	BlockTypeText    BlockType = "text"
	BlockTypeCode    BlockType = "code"
	BlockTypeBullets BlockType = "bullets"
)

func NewTextBlock(s string) MessageBlock {
	return MessageBlock{
		Type: BlockTypeText,
		Text: Text{
			Value: s,
		},
	}
}

func NewCodeBlock(value string, language *string) MessageBlock {
	return MessageBlock{
		Type: BlockTypeCode,
		Code: Code{
			Language: language,
			Value:    value,
		},
	}
}

func NewBulletsBlock(textBlocks []MessageBlock) MessageBlock {
	return MessageBlock{
		Type: BlockTypeBullets,
		Bullets: Bullets{
			Blocks: textBlocks,
		},
	}
}

// Block := Text | Code | Bullets
type MessageBlock struct {
	Type BlockType

	Text    Text    // Populated if Type is Text
	Code    Code    // Populated if Type is Code
	Bullets Bullets // Populated if Type is Bullets
}

func (b MessageBlock) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"type": b.Type,
	}
	switch b.Type {
	case BlockTypeText:
		m["value"] = b.Text.Value

	case BlockTypeCode:
		if b.Code.Language != nil {
			m["language"] = *b.Code.Language
		}
		m["value"] = b.Code.Value

	case BlockTypeBullets:
		m["blocks"] = b.Bullets.Blocks

	default:
		return nil, errors.New("unknown block type")
	}

	return marshalJSONWithoutEscapeHTML(m)
}

// Text := { type: "text", value: ANTLRString }
type Text struct {
	Value string
}

// Code := { type: "code", language: String | null, value: String }
type Code struct {
	Language *string
	Value    string
}

// Bullets := { type: "bullets", blocks: [Text] }
type Bullets struct {
	Blocks []MessageBlock
}

type Button struct {
	Title        string `json:"title"`
	ColorVariant int    `json:"color_variant"`
	Action       struct {
		WebAction struct {
			Attributes struct {
				URL string `json:"url"`
			} `json:"attributes"`
		} `json:"web_action"`
	} `json:"action"`
}

// Note: cannot use json.Marshaler as it always escapes HTML characters (<, >, &)
// e.g. "<b>" encodes to "\u003cb\u003e"
// https://pkg.go.dev/encoding/json#Marshal
func marshalJSONWithoutEscapeHTML(m any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(m)
	if err != nil {
		return nil, err
	}
	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}

// ANTLRString := (Pattern | String)+
// Pattern     := Emoji | Mention | Variable | Bold | Italic | InlineLink

// Emoji     := :{EmojiName}:
// EmojiName := RegExp(/^[-+_0-9a-zA-Z]+$/)
func Emoji(name string) string {
	return fmt.Sprintf(":%s:", name)
}

// MentionType := "manager" | "team"
type MentionType string

const (
	MentionTypeManager MentionType = "manager"
	MentionTypeTeam    MentionType = "team"
)

// Mention     := <link type="{MentionType}" value="{MentionId}">{MentionName}</link>
// MentionId   := EscapedString
// MentionName := {ANTLRString}
func Mention(mt MentionType, id string, mentionName string) string {
	return fmt.Sprintf("<link type=\"%s\" value=\"%s\">%s</link>", mt, EscapedString(id), mentionName)
}

// Bold := <b>{ANTLRString}</b>
func Bold(s string) string {
	return fmt.Sprintf("<b>%s</b>", s)
}

// Italic := <i>{ANTLRString}</i>
func Italic(s string) string {
	return fmt.Sprintf("<i>%s</i>", s)
}

// InlineLink     := <link type="url" value="{InlineLinkHref}">{ANTLRString}</link>
// InlineLinkHref := EscapedString
func InlineLink(href string, s string) string {
	return fmt.Sprintf("<link type=\"url\" value=\"%s\">%s</link>", EscapedString(href), s)
}

// Variable    := ${{VariableKey}} | ${{VariableKey}|{VariableAlt}}
// VariableKey := RegExp(/^\w+(?:\.[^<>.\s|$]+)*$/) | ""
// VariableAlt := RegExp(/^[^\s}]*[^\v}]+[^\s}]*$/) | ""
func Variable(key, alt string) string {
	if alt == "" {
		return fmt.Sprintf("${%s}", key)
	}
	return fmt.Sprintf("${%s|%s}", key, alt)
}

var escapeMap = map[rune]string{
	'"': "&quot;",
	'&': "&amp;",
	'<': "&lt;",
	'>': "&gt;",
}

// EscapedString := String - ("\"" | "&" | "<" | ">")
func EscapedString(s string) string {
	var b strings.Builder
	for _, c := range s {
		if escaped, ok := escapeMap[c]; ok {
			b.WriteString(escaped)
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}
