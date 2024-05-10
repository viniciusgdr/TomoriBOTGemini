package geminiServices

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GenerationConfig struct {
	Temperature     float64  `json:"temperature"`
	TopK            int      `json:"topK"`
	TopP            int      `json:"topP"`
	MaxOutputTokens int      `json:"maxOutputTokens"`
	StopSequences   []string `json:"stopSequences"`
}

type SafetySettings struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type RequestBody struct {
	Contents         []Content        `json:"contents"`
	GenerationConfig GenerationConfig `json:"generationConfig"`
	SafetySettings   []SafetySettings `json:"safetySettings"`
}

type response struct {
	Message *string `json:"message"`
	Query   *string `json:"query"`
	Command *string `json:"command"`
}

type Response struct {
	Message string `json:"message"`
	Query   string `json:"query"`
	Command string `json:"command"`
}

var prompt = ""

func GeminiChat(input string, history []*genai.Content) (*Response, error) {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_APIKEY")))
	if err != nil {
		return nil, err
	}
	defer client.Close()
	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-pro")
	// Initialize the chat
	cs := model.StartChat()
	cs.History = history

	resp, err := cs.SendMessage(ctx, genai.Text(prompt), genai.Text("input: "+input))
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates found")
	}

	firstCandidate := resp.Candidates[0]
	if len(firstCandidate.Content.Parts) == 0 {
		return nil, fmt.Errorf("no parts found")
	}

	str := resp.Candidates[0].Content.Parts
	var message string

	for _, part := range str {
		message += fmt.Sprintf("%v", part)
	}
	message = strings.ReplaceAll(message, "```json", "")
	message = strings.ReplaceAll(message, "```JSON", "")
	message = strings.ReplaceAll(message, "```", "")
	message = strings.TrimSpace(message)
	var response response
	err = json.Unmarshal([]byte(message), &response)
	if err != nil {
		return nil, err
	}

	fmt.Println(response)
	var res Response
	if response.Message != nil {
		res.Message = *response.Message
	}
	if response.Query != nil {
		res.Query = *response.Query
	}
	if response.Command != nil {
		res.Command = *response.Command
	}

	return &res, nil
}

func gemini(input string) (*Response, error) {
	URL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.0-pro-001:generateContent?key=" + os.Getenv("GEMINI_APIKEY")
	data := RequestBody{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
					{Text: "input: " + input},
				},
			},
		},
		GenerationConfig: GenerationConfig{
			Temperature:     0.9,
			TopK:            1,
			TopP:            1,
			MaxOutputTokens: 1512,
			StopSequences:   []string{},
		},
		SafetySettings: []SafetySettings{
			{Category: "HARM_CATEGORY_HARASSMENT", Threshold: "BLOCK_ONLY_HIGH"},
			{Category: "HARM_CATEGORY_HATE_SPEECH", Threshold: "BLOCK_ONLY_HIGH"},
			{Category: "HARM_CATEGORY_SEXUALLY_EXPLICIT", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
			{Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "BLOCK_ONLY_HIGH"},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// convert to interface
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if result["error"] != nil {
		return nil, fmt.Errorf("error: %v", result["error"])
	}
	candidates := result["candidates"].([]interface{})
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates found")
	}
	firstCandidate := candidates[0].(map[string]interface{})
	if firstCandidate["content"] == nil {
		return nil, fmt.Errorf("no content found")
	}
	content := firstCandidate["content"].(map[string]interface{})
	parts := content["parts"].([]interface{})
	if len(parts) == 0 {
		return nil, fmt.Errorf("no parts found")
	}
	text := ""
	for _, part := range parts {
		text += part.(map[string]interface{})["text"].(string)
	}

	text = strings.ReplaceAll(text, "```json", "")
	text = strings.ReplaceAll(text, "```", "")
	text = strings.TrimSpace(text)

	var response response
	err = json.Unmarshal([]byte(text), &response)
	if err != nil {
		return nil, err
	}

	fmt.Println(response)
	var res Response
	if response.Message != nil {
		res.Message = *response.Message
	}
	if response.Query != nil {
		res.Query = *response.Query
	}
	if response.Command != nil {
		res.Command = *response.Command
	}

	return &res, nil
}

func MakeLoopCallsIfErrorGemini(input string, history []*genai.Content, loopInt int) (*Response, error) {
	response, err := GeminiChat(input, history)
	if err != nil {
		if loopInt > 10 {
			return nil, err
		}
		return MakeLoopCallsIfErrorGemini(input, history, loopInt+1)
	}
	return response, nil
}
