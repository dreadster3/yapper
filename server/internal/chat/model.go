package chat

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatMessage struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}
