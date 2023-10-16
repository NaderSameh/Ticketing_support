package token

import (
	"github.com/golang-jwt/jwt"
)

// Payload contains the payload data of the token
type Payload struct {
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}
