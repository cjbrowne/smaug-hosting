package users

import (
	"database/sql/driver"
	"strings"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleMod   Role = "mod"
	RoleAdmin Role = "admin"
)

func (r Role) String() string {
	return string(r)
}

func (r *roleArr) Scan(src interface{}) error {
	s := string(src.([]uint8))
	if s == "" {
		*r = make([]Role, 0)
		return nil
	}

	for _, role := range strings.Split(s, ",") {
		*r = append(*r, MakeRole(role))
	}
	return nil
}

func (r roleArr) Value() (driver.Value, error) {
	roleStrArr := make([]string, 0)
	for _, role := range r {
		roleStrArr = append(roleStrArr, role.String())
	}
	return strings.Join(roleStrArr, ","), nil
}

func MakeRole(s string) Role {
	switch strings.ToLower(s) {
	case "admin":
		return RoleAdmin
	case "mod":
		fallthrough
	case "moderator":
		return RoleMod
	case "user":
		return RoleUser
	default:
		// return safe default
		return RoleUser
	}
}
