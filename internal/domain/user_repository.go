package domain

type UserRepository interface {
	Create(user *User) error
	FindByLogin(email string) (*User, error)
	FindByID(id string) (*User, error)
}