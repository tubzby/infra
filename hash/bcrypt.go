package hash

import "golang.org/x/crypto/bcrypt"

type bcryptHasher struct {
}

func (b *bcryptHasher) Hash(pwd []byte) (hashed []byte, err error) {
	hashed, err = bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
	}
	return
}

func (b bcryptHasher) Match(hashed, pwd []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hashed, pwd); err != nil {
		return false
	}
	return true
}
