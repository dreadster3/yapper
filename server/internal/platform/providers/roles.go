package providers

type Role int

const (
	RoleUser = iota
	RoleAssistant
	RoleSystem
)

func ParseRole(role string) Role {
	switch role {
	case "user":
		return RoleUser
	case "assistant":
		return RoleAssistant
	case "system":
		return RoleSystem
	default:
		return RoleUser
	}
}

func (r Role) String() string {
	switch r {
	case RoleUser:
		return "user"
	case RoleAssistant:
		return "assistant"
	case RoleSystem:
		return "system"
	default:
		return "user"
	}
}
