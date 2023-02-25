package jwt

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"log"
)

// DecodeToken 解码token
type DecodeToken struct {
	Iat      int    `json:"iat"`
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Iss      string `json:"iss"`
	IsAdmin  bool   `json:"isAdmin"`
}

// GenerateToken 生成Token
func GenerateToken(claims *jwt.Token, secret string) (token string) {
	hmacSecretString := secret
	hmacSecret := []byte(hmacSecretString)
	token, _ = claims.SignedString(hmacSecret)
	return
}

// VerifyToken 验证Token
func VerifyToken(token string, secret string) *DecodeToken {
	hmacSecretString := secret
	hmacSecret := []byte(hmacSecretString)
	decoded, err := jwt.Parse(
		token, func(token *jwt.Token) (interface{}, error) {
			return hmacSecret, nil
		})
	if err != nil {
		return nil
	}
	if !decoded.Valid {
		return nil
	}
	decodedClaims := decoded.Claims.(jwt.MapClaims)
	var decodedToken DecodeToken
	jsonString, _ := json.Marshal(decodedClaims)
	jsonErr := json.Unmarshal(jsonString, &decodedClaims)
	if jsonErr != nil {
		log.Print(jsonErr)
	}
	return &decodedToken
}
