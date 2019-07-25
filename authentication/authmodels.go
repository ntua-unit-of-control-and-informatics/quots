package authentication

import "github.com/dgrijalva/jwt-go"

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthResponse struct {
	ApiKey string `json:"apikey"`
}

type ValidationResponse struct {
	ApiKey string `json:"apikey"`
	Valid  bool   `json:"valid"`
}

type RefreshedToken struct {
	ApiKey string `json:"refreshedtoken"`
}
