package dto

type UserRegistrationData struct {
	FullName string
	Email    string
	Password string
}

type UserLoginData struct {
	Email    string
	Password string
}

type UserLogoutData struct {
	Token string
}
