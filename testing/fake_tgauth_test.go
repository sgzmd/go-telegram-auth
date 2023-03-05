package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeTgAuth(t *testing.T) {
	auth := NewFakeTelegramAuth(true, "username")
	ok, err := auth.CheckAuth(nil)

	assert.Nil(t, err)
	assert.True(t, ok)

	auth2 := NewFakeTelegramAuth(false, "username")
	ok2, err2 := auth2.CheckAuth(nil)

	assert.Nil(t, err2)
	assert.False(t, ok2)
}

func TestGetFakeParams(t *testing.T) {
	auth := NewFakeTelegramAuth(true, "username")
	p, e := auth.GetParamsFromCookieValue("123")
	assert.Nil(t, e)
	assert.Equal(t, FAKE_PARAMS, p)
}
