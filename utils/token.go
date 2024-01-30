package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id string) (string,error){
     loggedInToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
	 jwt.MapClaims{
		"id":id,
		"exp":time.Now().Add(time.Hour * 24).Unix(),
	 })
	
	 tokenString,err:=loggedInToken.SignedString([]byte(os.Getenv("SECRET")))

	 fmt.Println(tokenString)

	 if err != nil {
		return "", err
		}
	
	 return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
    tokenError := errors.New("Invalid token")
	if err != nil {
		return tokenError
	}
	if !token.Valid {
		return tokenError
	}
	return nil
}