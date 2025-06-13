package providers

type Role int

const (
	RoleUser = iota
	RoleAssistant
)

func ParseRole(role string) Role {
	switch role {
	case "user":
		return RoleUser
	case "assistant":
		return RoleAssistant
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
	default:
		return "user"
	}
}
