package configs

import "github.com/gorilla/sessions"

type CookieConfig struct {
	Path     string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

func CookieOption() *sessions.Options {
	return &sessions.Options{
		Path:     GetConfigData().CookieConfig.Path,
		MaxAge:   GetConfigData().CookieConfig.MaxAge,
		Secure:   GetConfigData().CookieConfig.Secure, // ssl 만 되는거 같은데 확실하지 않음
		HttpOnly: GetConfigData().CookieConfig.HttpOnly,
	}
}
