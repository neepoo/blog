package crypto_mac

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMac(t *testing.T) {
	const (
		size    = 32
		message = "salary:88888"
	)
	key := make([]byte, size)
	ri, err := rand.Read(key)
	require.Equal(t, size, ri)
	require.NoError(t, err)

	// sign
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	authTag := h.Sum(nil)

	// verify
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	got := h.Sum(nil)
	require.True(t, hmac.Equal(authTag, got))
}
