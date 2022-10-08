package tgauth

import "net/http"

type UserInfo struct {
	UserName  string
	FirstName string
	PhotoURL  string
	Id        int64
}

type Params map[string]string

// TelegramAuth is the main interface for Telegram authentication.
type TelegramAuth interface {
	// CheckAuth checks for a given set of params (usually from a request) if the user has successfully logged in
	// with Telegram. Returns true/false or error if invalid data.
	CheckAuth(params Params) (bool, error)

	// GetUserInfo returns UserInfo from the map of params.
	GetUserInfo(params Params) (*UserInfo, error)

	// GetParamsFromCookie returns the params from the cookie or error if no cookie present for a given http Request
	GetParamsFromCookie(req *http.Request) (Params, error)

	// GetParamsFromCookieValue returns the params from cookie value obtained by the caller
	GetParamsFromCookieValue(value string) (Params, error)

	// SetCookie sets the cookie for the user from the params
	SetCookie(w http.ResponseWriter, params Params) error

	// CreateCookie creates a cookie which the caller can set directly
	CreateCookie(params Params) (*http.Cookie, error)

	// GetCookieValue returns cookie value to be set by the caller
	GetCookieValue(params Params) (string, error)

	// SetDebug sets the debug flag for verbose logging
	SetDebug(debug bool) error
}
