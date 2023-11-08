package client

import (
	"github.com/lucas-kern/tower-of-babel_server/app/model"
)

//ClientUser is the model that governs all account objects retrieved or inserted into the DB
type ClientUser struct {
	FirstName    	*string            `json:"firstName,omitempty"`
	LastName    	 	*string            `json:"lastName,omitempty"`
	Token         	*string            `json:"token,omitempty"`
  RefreshToken 	*string            `json:"refreshToken,omitempty"`
  Base         	 	*model.Base              `json:"Base,omitempty"`
}

// NewClientUser sets up a client appropriate [model.User]
func NewClientUser(user *model.User) *ClientUser {
	return &ClientUser{
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Token:			 user.Token,
		RefreshToken:   user.RefreshToken,
	}
}