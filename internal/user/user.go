package user

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/router-server/pkg/storage"
	"golang.org/x/crypto/bcrypt"
)

const DefaultPassword = "root"

func EnsureCreated() error {
	log.Info("Checking default user")
	if storage.GetString("user.password") == "" {
		log.Warn("Default user not found, creating")
		return SetPassword(DefaultPassword)
	}
	return nil
}

func CheckPassword(password string) bool {
	hash := storage.GetString("user.password")
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed hash password: %s", err)
	}
	if err := storage.SetString("user.password", string(hash)); err != nil {
		return fmt.Errorf("failed save hash: %s", err)
	}
	return nil
}
