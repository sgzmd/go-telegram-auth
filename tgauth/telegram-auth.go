package tgauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultExpireTime = time.Hour * 24
	DefaultCookieName = "tg_auth"
)

// Main interface for Telegram authentication.
type TelegramAuth interface {
	// CheckAuth checks for a given set of params (usually from a request) if the user has successfully logged in
	// with Telegram. Returns true/false or error if invalid data.
	CheckAuth(params map[string][]string) (bool, error)

	// GetParamsFromCookie returns the params from the cookie or error if no cookie present
	GetParamsFromCookie(req *http.Request) (map[string][]string, error)

	// SetCookie sets the cookie for the user from the params
	SetCookie(w http.ResponseWriter, params map[string][]string) error
}

type TelegramAuthImpl struct {
	BotToken           string
	AuthUrl            string
	CheckAuthUrl       string
	TelegramCookieName string

	// After how long should the user be logged out? Defaults to 24 hours.
	ExpireTime         time.Duration
}

// NewTelegramAuth creates a new TelegramAuth instance.
//	botToken: The Telegram Bot API token
// 	authUrl: The URL to redirect to when the user is not authenticated (login page)
// 	checkAuthUrl: The URL to redirect for the actual authentication (should match the one in Telegram widget)
// 	telegramCookieName: The name of the cookie to store the Telegram user ID in, e.g. "tg_auth"
func NewTelegramAuth(botToken, authUrl, checkAuthUrl string) TelegramAuth {
	return TelegramAuthImpl{
		BotToken:           botToken,
		AuthUrl:            authUrl,
		CheckAuthUrl:       checkAuthUrl,
		TelegramCookieName: DefaultCookieName,
		ExpireTime:         DefaultExpireTime,
	}
}

// CheckAuth Checks if the user has successfully logged in with Telegram. It will return the
// json string of the user data if the user is logged in, otherwise it will return error.
func (t TelegramAuthImpl) CheckAuth(params map[string][]string) (bool, error) {
	expectedHash := calculateVerificationHash(params, t.BotToken)

	checkHash := params["hash"][0]

	// If the hashes match, then the request was indeed from Telegram
	if expectedHash != checkHash {
		return false, nil
	}

	// Now let's verify auth_date to check that the request is recent
	timestamp, err := strconv.ParseInt(params["auth_date"][0], 10, 64)
	if err != nil {
		return false, err
	}

	// User must login every 24 hours
	if timestamp < (time.Now().Unix() - int64(24*time.Hour.Seconds())) {
		return false, fmt.Errorf("user is not logged in for more than 24 hours")
	}

	return true, nil
}

// GetParamsFromCookie returns the params from the cookie or error if no cookie present
func (t TelegramAuthImpl) GetParamsFromCookie(req *http.Request) (map[string][]string, error) {
	// Get the cookie
	cookie, err := req.Cookie(t.TelegramCookieName)
	if err != nil {
		return nil, err
	}

	// Get the params from the cookie
	return t.getParamsFromCookie(cookie.Value)
}

// SetCookie sets the cookie for the user from the params
func (t TelegramAuthImpl) SetCookie(w http.ResponseWriter, params map[string][]string) error {
	cookie, err2 := t.createCookie(params)
	if err2 != nil {
		return err2
	}

	http.SetCookie(w, cookie)

	return nil
}

func (t TelegramAuthImpl) createCookie(params map[string][]string) (*http.Cookie, error) {
	j, err := json.Marshal(params)
	if err != nil {
		// This should practically never happen.
		return nil, fmt.Errorf("failed to marshal params to JSON: %+v", err)
	}
	// Set the cookie
	cookie := &http.Cookie{
		Name:    t.TelegramCookieName,
		Value:   url.QueryEscape(string(j)),
		Expires: time.Now().Add(t.ExpireTime),
		Path:    "/",
	}
	return cookie, nil
}

func (t TelegramAuthImpl) getParamsFromCookie(value string) (map[string][]string, error) {
	data, err := url.QueryUnescape(value)
	if err != nil {
		return nil, fmt.Errorf("error unescaping cookie value: %s", err)
	}

	params := make(map[string][]string)

	e := json.Unmarshal([]byte(data), &params)
	if e != nil {
		return nil, fmt.Errorf("error unmarshalling cookie value: %s", e)
	}

	return params, nil
}

// calculateVerificationHash: To check telegram login, we need to concat with "\n" all received fields _except_ hash
// sorted in alphabetical order and then calculate hash using sha256, with the bot api key hash
// as the secret.
func calculateVerificationHash(params map[string][]string, token string) string {
	keys := make([]string, 0)
	for k := range params {
		if k != "hash" {
			keys = append(keys, k)
		}
	}

	dataCheckArray := make([]string, len(keys))
	for i, k := range keys {
		// e.g. username=the_user
		dataCheckArray[i] = k + "=" + params[k][0]
	}

	// strings in array should be sorted in alphabetical order
	sort.Strings(dataCheckArray)

	// producing string like id=8889999222&first_name=sgzmd&username=the_user&photo_url=https%3A%2F%2Ft.me%2Fi%2Fu...
	dataCheckStr := strings.Join(dataCheckArray, "\n")

	s256 := sha256.New()
	s256.Write([]byte(token))

	// We will now use this secret key to produce hash-based authentication code
	// from the dataCheckStr produced above
	secretKey := s256.Sum(nil)

	hm := hmac.New(sha256.New, secretKey)
	hm.Write([]byte(dataCheckStr))
	expectedHash := hex.EncodeToString(hm.Sum(nil))

	return expectedHash
}
