package configs

import (
	"github.com/go-playground/validator"
	"github.com/gorilla/sessions"
)

type httpConfig struct {
	Service string `json:"service"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func getCookieStore() *sessions.CookieStore {
	// In real-world applications, use env variables to store the session key.
	sessionKey := "test-session-key"
	return sessions.NewCookieStore([]byte(sessionKey))
}
