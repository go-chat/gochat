package helper

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

// NewSalt generate random 128 character for user before register
func NewSalt() string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	symbols := big.NewInt(int64(len(alphanum)))
	states := big.NewInt(0)
	states.Exp(symbols, big.NewInt(int64(128)), nil)
	r, err := rand.Int(rand.Reader, states)

	if err != nil {
		panic(err)
	}

	var bytes = make([]byte, 128)
	r2 := big.NewInt(0)
	symbol := big.NewInt(0)

	for i := range bytes {
		r2.DivMod(r, symbols, symbol)
		r, r2 = r2, r
		bytes[i] = alphanum[symbol.Int64()]
	}

	return string(bytes)
}

func HashPassword(password, salt string) string {
	mac := hmac.New(sha256.New, []byte(salt))
	mac.Write([]byte(password))
	return hex.EncodeToString(mac.Sum(nil))
}
