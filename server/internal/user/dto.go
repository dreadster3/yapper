package user

type LoginUser struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type RegisterUser struct {
	Name     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
