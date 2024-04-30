package gemini_helper

import (
	"encoding/json"
	"github.com/google/generative-ai-go/genai"
)

func ParsingTextResponseForChat(response *genai.GenerateContentResponse) (resultString string, err error) {
	marshalResponse, _ := json.MarshalIndent(response, "", "  ")

	var generateResponse ContentResponse
	if err := json.Unmarshal(marshalResponse, &generateResponse); err != nil {
		return "", err
	}
	for _, cad := range *generateResponse.Candidates {
		if cad.Content != nil {
			for _, part := range cad.Content.Parts {
				resultString += part
			}
		}
	}

	return resultString, err
}
