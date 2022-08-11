package foo

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
)

// 可以参考的资料：https://blog.boot.dev/cryptography/how-sha-2-works-step-by-step-sha-256/
// https://github.com/golang/go/blob/master/src/crypto/sha256/sha256block.go
func ExampleSHA256() {
	const hw = "hello world"
	sha := sha256.New()
	sha.Write([]byte(hw))
	sum := sha.Sum(nil)
	fmt.Printf("%x", sum)
	// Output:
	// b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9
}

func ExampleCShake() {
	c1 := sha3.NewCShake128([]byte("my_hash"), nil)
	out := make([]byte, 16)
	c1.Write([]byte("hello world"))
	c1.Read(out)
	fmt.Printf("%x", out)
	// Output:
	// 7d5b3c7cb868c51a47b7a5cfb92d6d18
}

func TestJsonSerializer(t *testing.T) {
	m := map[string]string{
		"1": "1",
		"2": "2",
		"3": "3",
		"4": "4",
	}
	b, err := json.Marshal(m)
	require.NoError(t, err)

	b1, err := json.Marshal(m)
	require.NoError(t, err)

	require.Equal(t, b, b1)
}
