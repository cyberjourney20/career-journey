package llmrepo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cyberjourney20/career-journey/internal/config"
	"github.com/cyberjourney20/career-journey/internal/repository"
)

type OllamaRepo struct {
	App    *config.AppConfig
	Model  string
	Host   string
	Stream bool
}

func NewOllamaRepo(a *config.AppConfig) repository.LLMRepo {
	return &OllamaRepo{
		App:   a,
		Host:  os.Getenv("LLM_HOST"),
		Model: os.Getenv("LLM_MODEL"),
	}
}

func (o *OllamaRepo) OllamaGenerateResponse(prompt string, stream bool) (string, error) {
	fmt.Println("Running OllamaGenerateResponse with prompt:", prompt)

	request := map[string]interface{}{
		"model":  o.Model,
		"prompt": prompt,
		"stream": stream,
	}
	jsonData, _ := json.Marshal(request)

	resp, err := http.Post(o.Host, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error calling LLM: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	fmt.Println("Raw Ollama Response:", string(body)) // Debugging: print the full response

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error decoding JSON: %w", err)
	}

	fmt.Printf("Parsed JSON Response: %+v\n", response) // Debugging: print parsed JSON

	// Ensure "response" field exists
	respVal, exists := response["response"]
	if !exists {
		fmt.Println("Error: Missing 'response' field in JSON")
		return "", fmt.Errorf("missing 'response' key in Ollama response: %v", response)
	}

	// Ensure "response" is a string before using it
	respText, ok := respVal.(string)
	if !ok {
		fmt.Println("Error: 'response' field is not a string")
		return "", fmt.Errorf("'response' key is not a string: %v", response)
	}

	// Debug: Print response length before slicing
	fmt.Printf("Response Length: %d\n", len(respText))

	// **Check if the response is empty before slicing**
	if len(respText) == 0 {
		fmt.Println("Error: Ollama response is empty")
		return "", fmt.Errorf("ollama response is empty")
	}

	fmt.Println("Returning Response from Ollama:", respText)
	return respText, nil
}
