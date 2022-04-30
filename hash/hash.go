package hash

// Hasher .
type Hasher interface {
	Hash(pwd []byte) ([]byte, error)
	Match(hashed, pwd []byte) bool
}

// New create hasher
func New() Hasher {
	return &bcryptHasher{}
}
