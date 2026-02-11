// internal/utils/crypto/password.go
package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// 加密成本，越高越安全但也越慢
	bcryptCost = 12
)

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
