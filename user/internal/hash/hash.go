package hash

import "golang.org/x/crypto/bcrypt"

func HashPass(pass string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

func VerifyPass(pass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}
