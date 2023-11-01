package client

import (
	"github.com/lucas-kern/tower-of-babel_server/app/model"
)

//ClientUser is the model that governs all account objects retrieved or inserted into the DB
type ClientUser struct {
	First_name    	*string            `json:"first_name,omitempty"`
	Last_name    	 	*string            `json:"last_name,omitempty"`
	Token         	*string            `json:"token,omitempty"`
  Refresh_token 	*string            `json:"refresh_token,omitempty"`
  Base         	 	*model.Base              `json:"Base,omitempty"`
}

// NewClientUser sets up a client appropriate [model.User]
func NewClientUser(user *model.User) *ClientUser {
	return &ClientUser{
		First_name:      user.First_name,
		Last_name:       user.Last_name,
		Token:			 user.Token,
		Refresh_token:   user.Refresh_token,
	}
}