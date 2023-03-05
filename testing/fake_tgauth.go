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

var FAKE_PARAMS = tgauth.Params{
	"id":         "123",
	"first_name": "John",
	"username":   "john",
	"photo_url":  "http://example.com/photo.jpg",
	"auth_date":  "1234567890",
	"hash":       "1234567890",
}

func (f FakeTelegramAuth) SetDebug(debug bool) error {
	//TODO implement me
	panic("implement me")
}

func (f FakeTelegramAuth) GetCookieValue(_ tgauth.Params) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeTelegramAuth) CreateCookie(_ tgauth.Params) (*http.Cookie, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeTelegramAuth) GetParamsFromCookieValue(value string) (tgauth.Params, error) {
	return FAKE_PARAMS, nil
}

func (f FakeTelegramAuth) GetUserInfo(_ tgauth.Params) (*tgauth.UserInfo, error) {
	return &tgauth.UserInfo{
		UserName:  f.UserName,
		FirstName: f.UserName,
		PhotoURL:  "https://www.google.com/s2/favicons?domain=google.com&sz=64",
	}, nil
}

// SetCookie Implements TelegramAuth.SetCookie method for FakeTelegramAuth
func (f FakeTelegramAuth) SetCookie(_ http.ResponseWriter, _ tgauth.Params) error {
	return nil
}

// GetParamsFromCookie Implements TelegramAuth.GetParamsFromCookie method for FakeTelegramAuth
func (f FakeTelegramAuth) GetParamsFromCookie(_ *http.Request) (tgauth.Params, error) {
	return FAKE_PARAMS, nil
}

// CheckAuth Implements TelegramAuth.CheckAuth method for FakeTelegramAuth
func (f FakeTelegramAuth) CheckAuth(_ tgauth.Params) (bool, error) {
	return f.Pass, nil
}

func NewFakeTelegramAuth(pass bool, username string) tgauth.TelegramAuth {
	return FakeTelegramAuth{Pass: pass, UserName: username, Auth: tgauth.TelegramAuthImpl{}}
}
