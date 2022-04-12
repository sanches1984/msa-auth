package jwt

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	jwt := NewService(time.Hour, 6*time.Hour, "secret")

	user := int64(123)
	session := uuid.NewV4()
	token, err := jwt.NewAccessToken(user, session)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	userID, sessionID, err := jwt.ParseToken(token.Value)
	require.NoError(t, err)
	require.Equal(t, user, userID)
	require.Equal(t, session, sessionID)
}

func TestJWT_Expired(t *testing.T) {
	jwt := NewService(time.Second, time.Second, "secret")

	user := int64(123)
	session := uuid.NewV4()
	token, err := jwt.NewAccessToken(user, session)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	time.Sleep(2 * time.Second)

	_, _, err = jwt.ParseToken(token.Value)
	require.EqualError(t, err, "token is expired by 1s")
}

func TestJWTNew(t *testing.T) {
	jwt := NewService(time.Hour, 6*time.Hour, "secret")

	user := int64(123)
	session := uuid.NewV4()
	token, err := jwt.NewAccessToken(user, session)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	time.Sleep(time.Second)

	tokenNew, err := jwt.NewAccessToken(user, session)
	require.NoError(t, err)
	require.NotEmpty(t, tokenNew)
	require.NotEqual(t, token.Value, tokenNew.Value)
}
