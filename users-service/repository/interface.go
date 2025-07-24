package repository

// RepositoryInterface so Repository can be used in tests with dependency injection
type RepositoryInterface interface {
	GetUsers() ([]User, error)
	GetUserByEmail(email string) (*User, error)
	Disconnect() error
}
