package actions

import (
	"fmt"
	"tomoribot-geminiai-version/src/commands"
	command_types "tomoribot-geminiai-version/src/commands/types"

	geminiServices "tomoribot-geminiai-version/src/services/gemini"

	"github.com/google/generative-ai-go/genai"
)

var history = make(map[string][]*genai.Content)

func ProcessorGeminiAI(props *command_types.CommandProps) {
	sender := props.Message.Info.Sender
	chat := props.Message.Info.Chat

	keyHistory := sender.ToNonAD().String() + chat.ToNonAD().String()

	if len(history[keyHistory]) > 10 {
		history[keyHistory] = history[keyHistory][1:]
	}

	response, err := geminiServices.MakeLoopCallsIfErrorGemini(props.Arg, history[keyHistory], 0)
	if err != nil {
		fmt.Println("Error in geminiServices.MakeLoopCallsIfErrorGemini", err)
		return
	}

	history[keyHistory] = append(history[keyHistory], &genai.Content{
		Role: "user",
		Parts: []genai.Part{
			genai.Text(props.Arg),
		},
	})
	if response.Message != "" {
		history[keyHistory] = append(history[keyHistory], &genai.Content{
			Role: "model",
			Parts: []genai.Part{
				genai.Text(response.Message),
			},
		})
		props.Reply(response.Message)
	}

	if response.Command != "" {
		command, _ := commands.GetCommand(response.Command)
		command.Execute(props)
	}
}