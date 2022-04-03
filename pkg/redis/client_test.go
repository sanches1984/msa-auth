package redis

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	client, err := NewClient(Config{
		Host:              "localhost:8112",
		Password:          "",
		Db:                0,
		ConnectionTimeout: 10 * time.Second,
		OperationTimeout:  5 * time.Second,
	})
	require.NoError(t, err)

	defer client.Close()

	err = client.Set("my_key", []byte("hello"))
	require.NoError(t, err)

	data, err := client.Get("my_key")
	require.NoError(t, err)
	require.Equal(t, "hello", string(data))

	err = client.Delete("my_key")
	require.NoError(t, err)
}
