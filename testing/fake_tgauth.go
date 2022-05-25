package testing

import (
	"net/http"

	"github.com/sgzmd/go-telegram-auth/tgauth"
)

type FakeTelegramAuth struct {
	Auth tgauth.TelegramAuthImpl
	Pass bool
}

// Implements TelegramAuth.SetCookie method for FakeTelegramAuth
func (f FakeTelegramAuth) SetCookie(w http.ResponseWriter, params map[string][]string) error {
	return nil
}

// Implements TelegramAuth.GetParamsFromCookie method for FakeTelegramAuth
func (f FakeTelegramAuth) GetParamsFromCookie(req *http.Request) (map[string][]string, error) {
	params := map[string][]string{
		"id":         {"123"},
		"first_name": {"John"},
		"username":   {"john"},
		"photo_url":  {"http://example.com/photo.jpg"},
		"auth_date":  {"1234567890"},
		"hash":       {"da26696b03d7e7d67ebe4388fa133425b588b16fc40210e8656fb648eadecd0f"},
	}

	return params, nil
}

// Implements TelegramAuth.CheckAuth method for FakeTelegramAuth
func (f FakeTelegramAuth) CheckAuth(params map[string][]string) (bool, error) {
	return f.Pass, nil
}

func NewFakeTelegramAuth(pass bool) tgauth.TelegramAuth {
	return FakeTelegramAuth{Pass: pass, Auth: tgauth.TelegramAuthImpl{}}
}
