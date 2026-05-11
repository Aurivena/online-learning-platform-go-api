package middleware

import "online-learning-platform-go-api/config"

type Middleware struct {
	token *config.TokenConfig
}

func NewMiddleware(token *config.TokenConfig) *Middleware {
	return &Middleware{token: token}
}
