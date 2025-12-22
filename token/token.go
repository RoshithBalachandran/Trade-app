package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Secret_Key = []byte("My_Secret_Key")

type Claim struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(email string) (string ,error){
	expir := time.Now().Add(30 * time.Minute)
	claim := &Claim{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expir),
			Issuer:    "Trade-app",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(Secret_Key)
}
