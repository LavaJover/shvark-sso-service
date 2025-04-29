package domain

type User struct {
	ID string
	Login string
	Password string // hashed password
}

func NewUser(login, hashedPassword string) (*User, error) {
	return &User{
		Login: login,
		Password: hashedPassword,
	}, nil
}

// Validation functions...