package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeTgAuth(t *testing.T) {
	auth := NewFakeTelegramAuth(true)
	ok, err := auth.CheckAuth(nil)

	assert.Nil(t, err)
	assert.True(t, ok)

	auth2 := NewFakeTelegramAuth(false)
	ok2, err2 := auth2.CheckAuth(nil)

	assert.Nil(t, err2)
	assert.False(t, ok2)
}
