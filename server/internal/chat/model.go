package chat

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatMessage struct {
	Provider string    `json:"provider" binding:"required,registered_provider"`
	Model    string    `json:"model" binding:"required"`
	Messages []Message `json:"messages" binding:"required"`
	Think    bool      `json:"think"`
}

type Chat struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
