package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	accessExp  = 24 * time.Hour
	refreshExp = 30 * 24 * time.Hour
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (m *Middleware) SetToken(c *gin.Context) {
	userId, existsId := c.Get("userId")
	role, existsRole := c.Get("role")

	if !existsId || !existsRole {
		return
	}

	uid := userId.(uint)
	uRole := role.(string)

	accessToken, err := generateToken(uid, uRole, accessExp, []byte(m.token.AccessToken))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := generateToken(uid, uRole, accessExp, []byte(m.token.RefreshToken))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.SetCookie("access_token", accessToken, int(accessExp.Seconds()), "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, int(refreshExp.Seconds()), "/", "", false, true)
}

func generateToken(userID uint, role string, duration time.Duration, secret []byte) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func parseToken(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
