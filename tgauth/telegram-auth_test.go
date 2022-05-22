package tgauth

import (
	a "github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
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
	a.Equal(t, params, params2)
}

func TestTelegramAuth_SetCookie(t *testing.T) {
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
	c := w.Result().Cookies()[0]

	a.Equal(t, cookie.Value, c.Value)
}
