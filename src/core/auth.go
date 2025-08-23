package core

import (
	"os"
	"runtime"
	"log"
	"fmt"
	"errors"
	"path/filepath"
)

type AuthManager struct {
	JwtFilePath string
	Jwt         string
}

func GetDefaultJwtPath() string {
	if runtime.GOOS == "windows" {
		log.Fatalln("Windows is not yet supported")
	}

	home := os.Getenv("HOME")

	if home == "" {
		log.Fatalln("$HOME is not set")
	}

	return fmt.Sprintf("%s/.urnetwork/jwt", home)
}

func NewAuthManager(jwtFilePath string) (*AuthManager, error) {
	if jwtFilePath == "" {
		jwtFilePath = GetDefaultJwtPath()
	}

	jwtDirectory := filepath.Dir(jwtFilePath)
	authManager := & AuthManager {
		JwtFilePath: jwtFilePath,
		Jwt: "",
	}

	err := os.MkdirAll(jwtDirectory, 0755)

	if err != nil {
		return authManager, err
	}

	jwt, err := os.ReadFile(jwtFilePath)

	if err != nil {
		return authManager, err
	}

	authManager.Jwt = string(jwt)
	return authManager, nil
}

func (authManager *AuthManager) PerformAuth(email string, password string) (*APINetwork, error) {
	result, err := AttemptLoginWithPassword(email, password)

	if err != nil {
		return nil, err
	}

	if result.Error != nil {
		return nil, errors.New(result.Error.Message)
	}

	if result.VerificationRequired != nil {
		return nil, errors.New("Your account is not fully verified yet, cannot log in")
	}

	if result.Network == nil {
		return nil, errors.New("The API returned an invalid response")
	}

	authManager.Jwt = result.Network.Jwt

	jwtDirectory := filepath.Dir(authManager.JwtFilePath)

	if err := os.MkdirAll(jwtDirectory, 0755); err != nil {
		return nil, err
	}

	if err := os.WriteFile(authManager.JwtFilePath, []byte(authManager.Jwt), 0644); err != nil {
		return nil, errors.New("Cannot save JWT: " + err.Error())
	}

	log.Println("Logged into network: ", result.Network.Name)
	return result.Network, nil
}
