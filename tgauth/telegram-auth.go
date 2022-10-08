package tgauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
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

type TelegramAuthImpl struct {
	BotToken           string
	AuthUrl            string
	CheckAuthUrl       string
	TelegramCookieName string

	// After how long should the user be logged out? Defaults to 24 hours.
	ExpireTime time.Duration

	Debug bool
}

func (t TelegramAuthImpl) SetDebug(debug bool) error {
	t.Debug = debug
	return nil
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
func (t TelegramAuthImpl) CheckAuth(params Params) (bool, error) {
	expectedHash := calculateVerificationHash(params, t.BotToken)

	if t.Debug {
		log.Printf("Calculated hash: %s for params %+v", expectedHash, params)
	}

	if checkHash, ok := params["hash"]; ok {

		// If the hashes match, then the request was indeed from Telegram
		if expectedHash != checkHash {

			if t.Debug {
				log.Printf("Hash mismatch: %s != %s", expectedHash, checkHash)
			}

			return false, nil
		}

		// Now let's verify auth_date to check that the request is recent
		timestamp, err := strconv.ParseInt(params["auth_date"], 10, 64)
		if err != nil {

			if t.Debug {
				log.Printf("Error parsing auth_date: %s", params["auth_date"])
			}

			return false, err
		}

		if t.Debug {
			// prints to log timestamp in string format
			log.Printf("Auth date: %s", time.Unix(timestamp, 0).String())
		}

		// User must login every 24 hours
		if timestamp < (time.Now().Unix() - int64(24*time.Hour.Seconds())) {
			return false, fmt.Errorf("user is not logged in for more than 24 hours")
		}

		if t.Debug {
			log.Printf("User is logged in")
		}
		return true, nil
	} else {
		return false, fmt.Errorf("no 'hash' element in params")
	}
}

// GetUserInfo implements GetUserInfo from the interface
func (t TelegramAuthImpl) GetUserInfo(params Params) (*UserInfo, error) {
	ui := UserInfo{}
	err := paramsToInfo(params, &ui)
	if err != nil {
		return nil, err
	} else {
		return &ui, nil
	}
}

// GetParamsFromCookie returns the params from the cookie or error if no cookie present
func (t TelegramAuthImpl) GetParamsFromCookie(req *http.Request) (Params, error) {
	// Get the cookie
	cookie, err := req.Cookie(t.TelegramCookieName)
	if err != nil {
		return nil, err
	}

	// Get the params from the cookie
	return t.GetParamsFromCookieValue(cookie.Value)
}

// SetCookie sets the cookie for the user from the params
func (t TelegramAuthImpl) SetCookie(w http.ResponseWriter, params Params) error {
	cookie, err2 := t.CreateCookie(params)
	if err2 != nil {
		return err2
	}

	http.SetCookie(w, cookie)

	return nil
}

func (t TelegramAuthImpl) GetCookieValue(params Params) (string, error) {
	c, e := t.CreateCookie(params)
	if e != nil {
		return "", e
	}

	return c.Value, nil
}

func (t TelegramAuthImpl) CreateCookie(params Params) (*http.Cookie, error) {
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

func (t TelegramAuthImpl) GetParamsFromCookieValue(value string) (Params, error) {
	data, err := url.QueryUnescape(value)
	if err != nil {
		return nil, fmt.Errorf("error unescaping cookie value: %s", err)
	}

	params := make(map[string]string)

	e := json.Unmarshal([]byte(data), &params)
	if e != nil {
		return nil, fmt.Errorf("error unmarshalling cookie value: %s", e)
	}

	return params, nil
}

// calculateVerificationHash: To check telegram login, we need to concat with "\n" all received fields _except_ hash
// sorted in alphabetical order and then calculate hash using sha256, with the bot api key hash
// as the secret.
func calculateVerificationHash(params Params, token string) string {
	keys := make([]string, 0)
	for k := range params {
		if k != "hash" {
			keys = append(keys, k)
		}
	}

	dataCheckArray := make([]string, len(keys))
	for i, k := range keys {
		// e.g. username=the_user
		dataCheckArray[i] = k + "=" + params[k]
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

// paramsToInfo converts params map to UserInfo
func paramsToInfo(params Params, ui *UserInfo) error {
	if len(params["id"]) == 0 {
		return fmt.Errorf("no id in params: %+v", params)
	}

	uid, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing id: %s", err)
	}
	ui.Id = uid
	if len(params["username"]) > 0 {
		ui.UserName = params["username"]
	} else {
		return fmt.Errorf("username is empty: %+v", params)
	}

	if len(params["photo_url"]) > 0 {
		ui.PhotoURL = params["photo_url"]
	}

	if len(params["first_name"]) > 0 {
		ui.FirstName = params["first_name"]
	}

	if len(params["photo_url"]) > 0 {
		ui.PhotoURL = params["photo_url"]
	}

	return nil
}
