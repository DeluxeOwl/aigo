package aigo_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	OllamaBase          = "http://localhost:11434"
	ChatCompletionsPath = "/v1/chat/completions"
)

func BuildChatCompletionsURL(base string, path string) string {
	return fmt.Sprintf("%s%s", base, path)
}

func TestOllama(t *testing.T) {
	url := BuildChatCompletionsURL(OllamaBase, ChatCompletionsPath)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, strings.NewReader(`
{
	"model": "qwen3:0.6B",
	"stream": false,
	"messages": [
		{
			"role": "system",
			"content": "You are a helpful assistant."
		},
		{
			"role": "user",
			"content": "Hello!"
		}
	]
}`))
	require.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			fmt.Print(string(buf[:n]))
		}
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
	}
}
