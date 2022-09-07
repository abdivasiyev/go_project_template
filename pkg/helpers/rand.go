package helpers

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func RandomNumber(n int) (string, error) {
	bigInt, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%04d", bigInt.Int64()), nil
}
