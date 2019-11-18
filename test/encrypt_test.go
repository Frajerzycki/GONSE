package test

import (
	"bytes"
	"github.com/ikcilrep/gonse/pkg/nse"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestEncrypt(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 1; i <= 512; i++ {
		data := make([]byte, i)
		keyBytes := make([]byte, 32)
		salt := make([]byte, 16)
		IV, err := nse.GenerateIV(i)
		if err != nil {
			t.Error(err)
		}
		if len(IV) != i {
			t.Errorf("IV length %v is not data length %v", len(IV), i)
		}
		key := new(big.Int)
		rand.Read(data)
		rand.Read(keyBytes)
		rand.Read(salt)
		key.SetBytes(keyBytes)
		ciphertext, err := nse.Encrypt(data, salt, IV, key)
		if err != nil {
			t.Error(err)
		}
		decryptedData, err := nse.Decrypt(ciphertext, salt, IV, key)
		if err != nil {
			t.Error(err)
		}

		if decryptedData == nil {
			t.Errorf("Decrypted data is nil")
		}
		if !bytes.Equal(data, decryptedData) {
			t.Errorf("%v is not %v", data, decryptedData)
		}
	}
}
