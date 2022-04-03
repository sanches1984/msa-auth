package jwt

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	jwt := NewService(time.Hour, 6*time.Hour, "secret")

	token, err := jwt.NewAccessToken(123, 456)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	userID, sessionID, err := jwt.ParseToken(token)
	require.NoError(t, err)
	require.Equal(t, int64(123), userID)
	require.Equal(t, int64(456), sessionID)
}

func TestJWT_Expired(t *testing.T) {
	jwt := NewService(time.Second, time.Second, "secret")

	token, err := jwt.NewAccessToken(123, 456)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	time.Sleep(2 * time.Second)

	_, _, err = jwt.ParseToken(token)
	require.EqualError(t, err, "token is expired by 1s")
}
