package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
	config "github.com/zarinit-routers/router-server/pkg/cloud/config"
)

type Token = string

type AuthRequest struct {
	NodeId         string `json:"id"`
	OrganizationId string `json:"groupId"`
	Passphrase     string `json:"passphrase"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func authenticate(c config.ConnectionConfig) (Token, error) {
	log.Info("Authenticating", "nodeId", GetHostID())

	request := AuthRequest{
		NodeId:         GetHostID(),
		OrganizationId: c.OrganizationId,
		Passphrase:     c.Passphrase,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed marshal request body to JSON: %s", err)
	}
	bodyBuffer := bytes.NewBuffer(requestBody)

	response, err := http.Post(c.GetAuthURL(), "application/json", bodyBuffer)
	if err != nil {
		return "", fmt.Errorf("failed send request to %q: %s", c.GetAuthURL(), err)
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
