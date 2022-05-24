package tgauth

// TODO: improve testing coverage significantly.

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

	a "github.com/stretchr/testify/assert"
)

func Test_getParamsFromCookie(t *testing.T) {
	auth := NewTelegramAuth("bot_token", "/auth", "/check")
	params := map[string][]string{
		"id":         {"123"},
		"first_name": {"John"},
		"username":   {"john"},
		"photo_url":  {"http://example.com/photo.jpg"},
		"auth_date":  {"1234567890"},
		"hash":       {"1234567890"},
	}

	cookie, err := auth.createCookie(params)
	a.Nil(t, err)
	a.NotNil(t, cookie)

	params2, err := auth.getParamsFromCookie(cookie.Value)
	a.Nil(t, err)
	a.Equal(t, params, params2)
}

func TestTelegramAuth_SetCookie_GetParamsFromCookie(t *testing.T) {
	auth := NewTelegramAuth("bot_token", "/auth", "/check")
	params := map[string][]string{
		"id":         {"123"},
		"first_name": {"John"},
		"username":   {"john"},
		"photo_url":  {"http://example.com/photo.jpg"},
		"auth_date":  {"1234567890"},
		"hash":       {"1234567890"},
	}

	cookie, err := auth.createCookie(params)
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
	auth := NewTelegramAuth("bot_token", "/auth", "/check")
	params := map[string][]string{
		"id":         {"123"},
		"first_name": {"John"},
		"username":   {"john"},
		"photo_url":  {"http://example.com/photo.jpg"},
		"auth_date":  {"1234567890"},
		"hash":       {"1234567890"},
	}

	cookie, err := auth.createCookie(params)
	a.Nil(t, err)
	a.NotNil(t, cookie)

	p2 := map[string][]string{}
	j, e := url.QueryUnescape(cookie.Value)
	a.Nil(t, e)
	json.Unmarshal([]byte(j), &p2)

	a.Equal(t, params, p2)
}

func TestCalculateVerificationHash(t *testing.T) {
	params := map[string][]string{
		"id":         {"123"},
		"first_name": {"John"},
		"username":   {"john"},
		"photo_url":  {"http://example.com/photo.jpg"},
		"auth_date":  {"1234567890"},
		"hash":       {"1234567890"},
	}

	hash := calculateVerificationHash(params, "bot_token")

	// This hash was calculated near manually, with implementation which was
	// manually verified to work with the Telegram API. Any regression here means
	// the implementation is broken. See also: https://xkcd.com/221/
	a.Equal(t, "da26696b03d7e7d67ebe4388fa133425b588b16fc40210e8656fb648eadecd0f", hash)
}
