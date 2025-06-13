package chat

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatMessage struct {
	Provider string    `json:"provider"`
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Think    bool      `json:"think"`
}
