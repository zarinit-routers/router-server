package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Token = string

func getAuthUrl() string {
	u := url.URL{
		Scheme: "http",
		Host:   getCloudAddress(),
		Path:   CloudAuthPath,
	}
	return u.String()
}

type AuthRequest struct {
	NodeId         string `json:"id"`
	OrganizationId string `json:"groupId"`
	Passphrase     string `json:"passphrase"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func Authenticate() (Token, error) {

	request := AuthRequest{
		NodeId:         GetHostID(),
		OrganizationId: getCloudGroupID(),
		Passphrase:     getCloudPassphrase(),
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed marshal request body to JSON: %s", err)
	}
	bodyBuffer := bytes.NewBuffer(requestBody)

	response, err := http.Post(getAuthUrl(), "application/json", bodyBuffer)
	if err != nil {
		return "", fmt.Errorf("failed send request to %q: %s", getAuthUrl(), err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed read response body: %s", err)
	}

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return "", fmt.Errorf("failed unmarshal response body: %s", err)
	}

	if authResponse.Token == "" {
		return "", fmt.Errorf("returned empty token")
	}

	return authResponse.Token, nil
}
