package testing

import (
	"net/http"

	"github.com/sgzmd/go-telegram-auth/tgauth"
)

type FakeTelegramAuth struct {
	Auth     tgauth.TelegramAuthImpl
	UserName string
	Pass     bool
}

func (f FakeTelegramAuth) GetUserInfo(_ map[string][]string) (*tgauth.UserInfo, error) {
	return &tgauth.UserInfo{
		UserName:  f.UserName,
		FirstName: f.UserName,
		PhotoURL:  "https://www.google.com/s2/favicons?domain=google.com&sz=64",
	}, nil
}

// SetCookie Implements TelegramAuth.SetCookie method for FakeTelegramAuth
func (f FakeTelegramAuth) SetCookie(w http.ResponseWriter, params map[string][]string) error {
	return nil
}

// GetParamsFromCookie Implements TelegramAuth.GetParamsFromCookie method for FakeTelegramAuth
func (f FakeTelegramAuth) GetParamsFromCookie(_ *http.Request) (map[string][]string, error) {
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

// CheckAuth Implements TelegramAuth.CheckAuth method for FakeTelegramAuth
func (f FakeTelegramAuth) CheckAuth(_ map[string][]string) (bool, error) {
	return f.Pass, nil
}

func NewFakeTelegramAuth(pass bool, username string) tgauth.TelegramAuth {
	return FakeTelegramAuth{Pass: pass, UserName: username, Auth: tgauth.TelegramAuthImpl{}}
}
