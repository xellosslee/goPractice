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
		Path:     ServerConfig.CookieConfig.Path,
		MaxAge:   ServerConfig.CookieConfig.MaxAge,
		Secure:   ServerConfig.CookieConfig.Secure, // ssl 만 되는거 같은데 확실하지 않음
		HttpOnly: ServerConfig.CookieConfig.HttpOnly,
	}
}
