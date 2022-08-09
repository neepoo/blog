package foo

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEsSha512(t *testing.T) {
	bs, err := os.ReadFile("/home/neepoo/下载/elasticsearch-8.3.1-linux-x86_64.tar.gz")
	require.NoError(t, err)
	h := sha512.New()
	_, err = h.Write(bs)
	require.NoError(t, err)
	got := h.Sum(nil)
	want, err := os.ReadFile("/home/neepoo/下载/elasticsearch-8.3.1-linux-x86_64.tar.gz.sha512")
	require.NoError(t, err)
	tokens := bytes.Fields(want)
	require.Equal(t, 2, len(tokens))
	require.Equal(t, hex.EncodeToString(tokens[0]), hex.EncodeToString(got))
}
