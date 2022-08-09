package foo

import (
	"bytes"
	"crypto/md5"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMD5Collision(t *testing.T) {
	const (
		ship  = "ship_coll.jpg"
		plane = "plane_coll.jpg"
	)
	f1, err := os.ReadFile(ship)
	require.NoError(t, err)
	h1 := md5.New()
	_, err = io.Copy(h1, bytes.NewReader(f1))
	require.NoError(t, err)

	f2, err := os.ReadFile(plane)
	require.NoError(t, err)
	h2 := md5.New()
	_, err = io.Copy(h2, bytes.NewReader(f2))
	require.NoError(t, err)

	require.NotEqual(t, f1, f2)
	require.Equal(t, h1.Sum(nil), h2.Sum(nil))
}
