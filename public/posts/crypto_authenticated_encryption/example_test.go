package example_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"reflect"
	"testing"
)

var (
	plainText = []byte("世界和平")
	key       = make([]byte, 16) // aes-128
	nonce     = make([]byte, 12)
	ad        = make([]byte, 32)
)

func init() {
	rand.Read(key)
	rand.Read(nonce)
	rand.Read(ad)
}

func encrypt(plaintext, key, nonce, ad []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(plainText[:0], nonce, plaintext, ad)

	return ciphertext
}

func decrypt(ciphertext, key, nonce, ad []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, ad)
	if err != nil {
		panic(err.Error())
	}

	return plaintext
}

func Test_aes_gcm(t *testing.T) {
	enData := encrypt(plainText, key, nonce, ad)
	deData := decrypt(enData, key, nonce, ad)
	//deData := decrypt(enData, key, nonce, append(ad, 2))
	if !reflect.DeepEqual(deData, plainText) {
		t.Fail()
	}
}
