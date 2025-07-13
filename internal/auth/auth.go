package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const secretKey = "supersecretkey"
const AuthCookieName = "Authorization"
const UserIDKey = "userID"

// claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское UserID
type claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

func GenerateAccessToken(userID string) (string, error) {
	return generateToken(userID, 15*time.Minute)
}

func GenerateRefreshToken(userID string) (string, error) {
	return generateToken(userID, 7*24*time.Hour) // 7 дней
}

func generateToken(userID string, duration time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func GetUserIDFromReq(req *http.Request) (userID uuid.UUID, err error) {
	authCookie, err := req.Cookie(AuthCookieName)
	if err != nil {
		zap.L().Sugar().Debugln("No auth cookie found")
		return uuid.Nil, err
	}

	userID, err = getUserID(authCookie.Value)
	if err != nil {
		zap.L().Sugar().Debugln("Error getting user id from cookie. Message: ", err.Error())
		return uuid.Nil, err
	}

	return userID, nil
}

func getUserID(tokenString string) (userID uuid.UUID, err error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}
