package core

import (
	"io"
	"net/http"
	"encoding/json"
	"bytes"
	"log"
)

type APINetwork struct {
	Jwt string `json:"by_jwt"`
	Name string `json:"name"`
}

type APIVerificationRequired struct {
	UserAuth string `json:"user_auth,omitempty"`
}

type APILoginError struct {
	Message string `json:"message"`
}

type APILoginWithPasswordResponse struct {
	Network *APINetwork `json:"network,omitempty"`
	VerificationRequired *APIVerificationRequired `json:"verification_required,omitempty"`
	Error *APILoginError `json:"error,omitempty"`
}

const API_BASE_URL = "https://api.bringyour.com"

func makeRoute(path string) string {
	return API_BASE_URL + path
}

func AttemptLoginWithPassword(email string, password string) (*APILoginWithPasswordResponse, error) {
	url := makeRoute("/auth/login-with-password")
	data := map[string]string {
		"user_auth": email,
		"password": password,
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var result APILoginWithPasswordResponse

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	log.Println("Response: ", result)
	return &result, nil
}
