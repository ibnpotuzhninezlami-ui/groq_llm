package main

import (
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "net/http"
    "os"
)

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
}

type Choice struct {
    Message Message `json:"message"`
}

type ChatResponse struct {
    Choices []Choice `json:"choices"`
}

func main() {
    apiKey := flag.String("key", os.Getenv("GROQ_API_KEY"), "Your Groq API key")
    prompt := flag.String("prompt", "Explain quantum computing in simple terms.", "Prompt")
    model := flag.String("model", "llama-3.3-70b-versatile", "Model: e.g., llama-3.3-70b-versatile (best current replacement), llama-3.1-8b-instant (fast/small), openai/gpt-oss-120b (very capable)")
    flag.Parse()

    if *apiKey == "" {
        fmt.Println("Error: Set -key or export GROQ_API_KEY=your_key_here")
        os.Exit(1)
    }

    requestBody := ChatRequest{
        Model: *model,
        Messages: []Message{
            {Role: "system", Content: "You are a helpful and direct assistant."},
            {Role: "user", Content: *prompt},
        },
    }

    jsonData, _ := json.Marshal(requestBody)

    req, _ := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+*apiKey)

    client := &http.Client{}
    resp, _ := client.Do(req)
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    if resp.StatusCode != http.StatusOK {
        fmt.Printf("Error (%d): %s\n", resp.StatusCode, string(body))
        os.Exit(1)
    }

    var chatResp ChatResponse
    json.Unmarshal(body, &chatResp)

    if len(chatResp.Choices) > 0 {
        fmt.Println("LLM response:")
        fmt.Println(chatResp.Choices[0].Message.Content)
    } else {
        fmt.Println("No response.")
    }
}
