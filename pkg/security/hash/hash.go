package hash

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
	"strings"
)

type Params struct {
	PlainText   string
	EncodedHash string
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func Md5Sum(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	md5Sum := md5.Sum(data)
	return fmt.Sprintf("%x", md5Sum), nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateHash(params Params) (string, error) {
	salt, err := generateRandomBytes(params.SaltLength)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(params.PlainText), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, params.Memory, params.Iterations, params.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func CompareHash(plainText, encodedHash string) (bool, error) {
	// Extract the parameters, salt and derived key from the encoded plainText
	// hash.
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other plainText using the same parameters.
	otherHash := argon2.IDKey([]byte(plainText), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (Params, []byte, []byte, error) {
	var (
		p       Params
		version int
	)

	values := strings.Split(encodedHash, "$")
	if len(values) != 6 {
		return p, nil, nil, ErrInvalidHash
	}

	_, err := fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return p, nil, nil, err
	}
	if version != argon2.Version {
		return p, nil, nil, ErrIncompatibleVersion
	}

	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return p, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return p, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return p, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}
