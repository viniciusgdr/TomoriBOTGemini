package sender

import "encoding/json"

type ButtonType string

const (
	ButtonURL   ButtonType = "url"
	ButtonCopy  ButtonType = "copy"
	ButtonReply ButtonType = "reply"
)

type InteractiveButtons struct {
	DisplayText string
	ID          string
	Type        ButtonType
}
type NativeFlowButtonURL struct {
	DisplayText string `json:"display_text"`
	URL         string `json:"url"`
	MerchantURL string `json:"merchant_url"`
}

func (b *NativeFlowButtonURL) toString() string {
	converted, _ := json.Marshal(b)
	return string(converted)
}

type NativeFlowButtonCopy struct {
	DisplayText string `json:"display_text"`
	ID          string `json:"id"`
	CopyCode    string `json:"copy_code"`
}

func (b *NativeFlowButtonCopy) toString() string {
	converted, _ := json.Marshal(b)
	return string(converted)
}

type NativeFlowButtonReply struct {
	DisplayText string `json:"display_text"`
	ID          string `json:"id"`
	Disabled    string `json:"disabled"`
}

func (b *NativeFlowButtonReply) toString() string {
	converted, _ := json.Marshal(b)
	return string(converted)
}

type NativeFlowListMessageSection struct {
	Title string                            `json:"title"`
	Rows  []NativeFlowListMessageSectionRow `json:"rows"`
}

type NativeFlowListMessageSectionRow struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Header      string `json:"header"`
	ID          string `json:"id"`
}

type ButtonParamsJsonV2 struct {
	Title string `json:"title"`
	Sections []NativeFlowListMessageSection `json:"sections"`
}

func (b *ButtonParamsJsonV2) toString() string {
	converted, _ := json.Marshal(b)
	return string(converted)
}

type Header struct {
	Title string `json:"title"`
	Subtitle string `json:"subtitle"`
	HasMediaAttachment bool `json:"hasMediaAttachment"`
	MediaByte []byte
}