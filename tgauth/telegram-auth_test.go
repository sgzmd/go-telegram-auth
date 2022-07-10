package tgauth

// TODO: improve testing coverage significantly.

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

	a "github.com/stretchr/testify/assert"
)

var params = Params{
	"id":         "123",
	"first_name": "John",
	"username":   "john",
	"photo_url":  "http://example.com/photo.jpg",
	"auth_date":  "1234567890",
	"hash":       "1234567890",
}

func NewTelegramAuthImpl() TelegramAuthImpl {
	return TelegramAuthImpl{
		BotToken:           "bot_token",
		AuthUrl:            "/auth",
		CheckAuthUrl:       "/check",
		TelegramCookieName: DefaultCookieName,
		ExpireTime:         DefaultExpireTime,
	}
}

func Test_getParamsFromCookie(t *testing.T) {
	auth := NewTelegramAuthImpl()

	cookie, err := auth.CreateCookie(params)
	a.Nil(t, err)
	a.NotNil(t, cookie)

	params2, err := auth.GetParamsFromCookieValue(cookie.Value)
	a.Nil(t, err)
	a.Equal(t, params, params2)
}

func TestTelegramAuthImpl_GetUserInfo_BadData(t *testing.T) {
	auth := NewTelegramAuthImpl()
	p2 := Params{}
	ui, err := auth.GetUserInfo(p2)
	a.Nil(t, ui)
	a.NotNil(t, err)
}

func TestTelegramAuth_SetCookie_GetParamsFromCookie(t *testing.T) {
	auth := NewTelegramAuthImpl()

	cookie, err := auth.CreateCookie(params)
	a.Nil(t, err)
	a.NotNil(t, cookie)

	w := httptest.NewRecorder()
	auth.SetCookie(w, params)

	w.Flush()
	r := w.Result()
	c := r.Cookies()[0]

	a.Equal(t, cookie.Value, c.Value)
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(c)
	p2, err := auth.GetParamsFromCookie(req)
	a.Nil(t, err)
	a.Equal(t, params, p2)
}

func TestCreateCookie(t *testing.T) {
	auth := NewTelegramAuthImpl()
	cookie, err := auth.CreateCookie(params)
	a.Nil(t, err)
	a.NotNil(t, cookie)

	p2 := Params{}
	j, e := url.QueryUnescape(cookie.Value)
	a.Nil(t, e)
	json.Unmarshal([]byte(j), &p2)

	a.Equal(t, params, p2)
}

func TestCalculateVerificationHash(t *testing.T) {

	hash := calculateVerificationHash(params, "bot_token")

	// This hash was calculated near manually, with implementation which was
	// manually verified to work with the Telegram API. Any regression here means
	// the implementation is broken. See also: https://xkcd.com/221/
	a.Equal(t, "da26696b03d7e7d67ebe4388fa133425b588b16fc40210e8656fb648eadecd0f", hash)
}

func TestParamsToInfo(t *testing.T) {
	auth := NewTelegramAuthImpl()

	info, err := auth.GetUserInfo(params)

	a.Nil(t, err)

	a.Equal(t, "John", info.FirstName)
	a.Equal(t, "john", info.UserName)
	a.Equal(t, "http://example.com/photo.jpg", info.PhotoURL)

	err = paramsToInfo(map[string]string{}, info)
	a.NotNil(t, err)
}

func TestTelegramAuthImpl_CheckAuth(t *testing.T) {
	auth := NewTelegramAuthImpl()
	params := make(map[string]string, 1)
	ok, err := auth.CheckAuth(params)
	a.False(t, ok)
	a.NotNil(t, err)
}
