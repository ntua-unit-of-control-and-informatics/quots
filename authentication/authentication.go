package authentication

import (
	"errors"
	"log"
	"os"
	db "quots/database"
	utils "quots/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var user_id string
var user_pass string
var signing_key string
var jwtKey []byte

func NewAuth() {
	user_id = os.Getenv("QUOTS_USER")
	if user_id == "" {
		user_id = "QUOTS"
	}
	user_pass = os.Getenv("QUOTS_PASSWORD")
	if user_pass == "" {
		user_pass = "QUOTS"
	}
	signing_key = utils.RandStringBytesMaskImprSrcUnsafe(18)
	secrets, err := db.GetJWTSignature()
	if err != nil {
		log.Panicf(err.Error())
	}
	if len(secrets) == 0 {
		db.SaveJWTSignature(signing_key)
	} else {
		jwtKey = []byte(secrets[0].Signature)
	}
}

func AuthenticateAdmin(username string, password string) (jwtN string, err error) {
	jwtNew := ""
	errnew := errors.New("Login error")
	if username == user_id && password == user_pass {
		expirationTime := time.Now().Add(50 * time.Minute)
		claims := &Claims{
			Username: username,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		jwtNew, errnew := token.SignedString(jwtKey)
		return jwtNew, errnew
	} else {
		return jwtNew, errnew
	}
}

func ValidateToken(tknStr string) (validated bool, err error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return tkn.Valid, err
	} else {
		return tkn.Valid, err
	}

}

func RenewToken(tknStr string) (tok string, err error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if tkn.Valid {
		expirationTime := time.Now().Add(5 * time.Minute)
		claims.ExpiresAt = expirationTime.Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		return tokenString, err
	} else {
		return tkn.Raw, err
	}

}
