package tokens

import (
	"bitbucket.org/smaug-hosting/services/idp/users"
	µ "bitbucket.org/smaug-hosting/services/micro"
	"bitbucket.org/smaug-hosting/services/registry/services"
	"crypto/rand"
	_ "crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type TokenClaims struct {
	UserId int64
	std    jwt.StandardClaims
}

func (t TokenClaims) Valid() error {
	// todo: actually validate the claims by mapping against the db
	return nil
}

func GenerateServiceToken(service services.Service) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, TokenClaims{
		UserId: -1,
		std: jwt.StandardClaims{
			ExpiresAt: 1500,
			Audience:  "service",
		},
	})

	sigBytes, err := ioutil.ReadFile(µ.MustGetEnv("JWT_PRIVKEY_FILE"))
	if err != nil {
		return "", err
	}

	sigKey, err := jwt.ParseRSAPrivateKeyFromPEM(sigBytes)
	if err != nil {
		return "", err
	}

	tokStr, err := tok.SignedString(sigKey)

	return tokStr, err
}

func GenerateToken(user users.User) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, TokenClaims{
		std: jwt.StandardClaims{
			ExpiresAt: 1500,
			Audience:  "user",
		},
		UserId: user.Id,
	})

	sigBytes, err := ioutil.ReadFile(µ.MustGetEnv("JWT_PRIVKEY_FILE"))
	if err != nil {
		return "", err
	}

	sigKey, err := jwt.ParseRSAPrivateKeyFromPEM(sigBytes)
	if err != nil {
		return "", err
	}

	tokStr, err := tok.SignedString(sigKey)

	return tokStr, err
}

func ParseToken(token string, target *TokenClaims) error {

	_, err := jwt.ParseWithClaims(token, target, func(token *jwt.Token) (interface{}, error) {
		pubKeyBytes, err := ioutil.ReadFile(µ.MustGetEnv("JWT_PUBKEY_FILE"))
		if err != nil {
			return nil, err
		}

		block, _ := pem.Decode(pubKeyBytes)

		crt, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		return crt, nil
	})

	return err
}

func GenerateRefreshToken() []byte {
	tok := make([]byte, 64)

	_, err := rand.Read(tok)
	if err != nil {
		logrus.Errorf("Could not generate token: %s", err)
		return nil
	}

	return tok
}
