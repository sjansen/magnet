package server

import (
	"encoding/gob"

	"github.com/crewjam/saml/samlsp"
)

func init() {
	gob.Register(User{})
}

// User describes the currently authenticated user.
type User struct {
	Email     string
	GivenName string
	Surname   string
	Roles     []string
}

// GetAttributes implements SessionWithAttributes.
func (u *User) GetAttributes() samlsp.Attributes {
	attrs := samlsp.Attributes{
		"email":     []string{u.Email},
		"firstName": []string{u.GivenName},
		"lastName":  []string{u.Surname},
		"roles":     u.Roles,
	}
	return attrs
}
