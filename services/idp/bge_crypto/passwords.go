package bge_crypto

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(password string) []byte {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		logrus.Errorf("Could not encrypt password: %s", err)
	}

	return hash
}

func Verify(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))

	return err == nil
}