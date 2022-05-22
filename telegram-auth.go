package go_telegram_auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TelegramAuth struct {
	BotToken       string
	AuthUrl        string
	CheckAuthUrl   string
	TelegramCookie string
}

// Checks if the user has successfully logged in with Telegram. It will return the
// json string of the user data if the user is logged in, otherwise it will return error.
func (t TelegramAuth) checkAuth(params map[string][]string) (map[string][]string, error) {
	expectedHash := calculateVerificationHash(params, t.BotToken)

	checkHash := params["hash"][0]

	// If the hashes match, then the request was indeed from Telegram
	if expectedHash != checkHash {
		return nil, fmt.Errorf("Hash mismatch")
	}

	// Now let's verify auth_date to check that the request is recent
	timestamp, err := strconv.ParseInt(params["auth_date"][0], 10, 64)
	if err != nil {
		return nil, err
	}

	// User must login every 24 hours
	if timestamp < (time.Now().Unix() - int64(24*time.Hour.Seconds())) {
		return nil, fmt.Errorf("User is not logged in for more than 24 hours")
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

// CheckAuth checks the Telegram authentication status of the user
func (t TelegramAuth) CheckAuth(params map[string][]string, setCookie bool) (bool, error) {
	_, err := t.checkAuth(params)
	if err != nil {
		return false, err
	}

	return true, nil
}
