package crypto

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestFoo(t *testing.T) {
	h := sha256.New()
	h.Write([]byte("test"))
	s := h.Sum(nil)
	for i := 0; i < len(s); i++ {
		fmt.Printf("%x", s[i])
	}
}
