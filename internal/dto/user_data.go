package dto

type UserRegistrationData struct {
	FullName string `db:"fullName" json:"fullName"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type UserUpdateData struct {
	ID       int64  `db:"id" json:"id"`
	FullName string `db:"fullName" json:"fullName"`
	Password string `db:"password" json:"password"`
}

type UserRestorePasswordData struct {
	Email string `db:"email" json:"email"`
}

type UserLoginData struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type UserLogoutData struct {
	Token string `db:"token" json:"token"`
}
