package tokens

import "github.com/sirupsen/logrus"

type Token struct {
	Token        string
}

func (token Token) AsMap() map[string]interface{} {
	claims := TokenClaims{}
	err := ParseToken(token.Token, &claims)
	logrus.Tracef("Converting token to map with user id %d", claims.UserId)
	if err != nil {
		logrus.Errorf("Could not parse token: %s", err)
		return nil
	}

	return map[string]interface{}{
		"token": token.Token,
		"user_id": claims.UserId,
	}
}

type TokenResponse struct {
	Token   string
}