package security

import (
	"github.com/abdivasiyev/project_template/pkg/security/hash"
)

func (p *handler) GenerateHash(plainText string) (string, error) {
	return hash.GenerateHash(hash.Params{
		PlainText:   plainText,
		Memory:      p.memory,
		Iterations:  p.iterations,
		Parallelism: p.parallelism,
		SaltLength:  p.saltLength,
		KeyLength:   p.keyLength,
	})
}
func (p *handler) CompareHash(plainText, encodedHash string) (bool, error) {
	return hash.CompareHash(plainText, encodedHash)
}

func (p *handler) Md5Sum(value any) (string, error) {
	return hash.Md5Sum(value)
}
