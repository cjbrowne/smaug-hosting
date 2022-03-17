package users

type UserLogin struct {
	Email    string
	Password string
}

type User struct {
	Id                int64
	Email             string
	PasswordHash      []byte `db:"password_hash"`
	Roles             roleArr
	Balance           int64
	Verified          bool `db:"-"`
	VerificationToken string `db:"verification_token"`
}

type roleArr []Role

func (user User) AsMap() map[string]interface{} {
	if user.Verified {
		return map[string]interface{}{
			"id":            user.Id,
			"email":         user.Email,
			"password_hash": user.PasswordHash,
			"roles":         user.Roles,
			"balance":       user.Balance,
		}
	} else {
		return map[string]interface{}{
			"id":                 user.Id,
			"email":              user.Email,
			"password_hash":      user.PasswordHash,
			"roles":              user.Roles,
			"balance":            user.Balance,
			"verification_token": user.VerificationToken,
		}
	}
}

func (user User) AsUpdateMap() map[string]interface{} {
	if user.Verified {
		return map[string]interface{}{
			"email":         user.Email,
			"password_hash": user.PasswordHash,
			"roles":         user.Roles,
			"balance":       user.Balance,
		}
	} else {
		return map[string]interface{}{
			"email":         user.Email,
			"password_hash": user.PasswordHash,
			"roles":         user.Roles,
			"balance":       user.Balance,
		}
	}
}

var Dummy User
