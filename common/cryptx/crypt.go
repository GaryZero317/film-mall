package cryptx

import (
	"encoding/hex"
	"fmt"
	"log"

	"golang.org/x/crypto/scrypt"
)

func PasswordEncrypt(salt, password string) string {
	dk, err := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	if err != nil {
		log.Printf("Error encrypting password: %v", err)
		return ""
	}
	
	result := hex.EncodeToString(dk)
	fmt.Printf("DEBUG Encryption - Password bytes: %x\n", []byte(password))
	fmt.Printf("DEBUG Encryption - Salt bytes: %x\n", []byte(salt))
	fmt.Printf("DEBUG Encryption - Result: %s\n", result)
	
	return result
}
