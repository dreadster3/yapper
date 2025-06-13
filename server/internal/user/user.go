package user

type UserId string

type User struct {
	Id       UserId
	Name     string
	Email    string
	Password string
}
