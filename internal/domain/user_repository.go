package domain

type UserRepository interface {
	Create(user *User) error
	FindByLogin(login string) (*User, error)
	FindByID(id string) (*User, error)
}