package cloud

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os/exec"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/zarinit-routers/router-server/pkg/commands"
)

const (
	CloudWebSocketPath = "/api/ipc/connect"

	AuthorizationHeader = "Authorization"
	RouterIDHeader      = "X-Router-ID"
)

var (
	conn *websocket.Conn
)

type Request struct {
	Command string         `json:"command"`
	Args    map[string]any `json:"args"`
	ID      string         `json:"requestId"`
}
type Response struct {
	Data  map[string]any `json:"data"`
	ID    string         `json:"requestId"`
	Error string         `json:"error"`
}

func ServeConnection() {
	for {

		if getCloudPassphrase() != "" {
			err := establishConnection()
			if err != nil {
				log.Error("Failed to establish connection", "error", err)
			}
		}

		timeout := getReconnectTimeout()
		log.Debug("Reconnect timeout", "timeout", timeout)
		time.Sleep(timeout)
	}
}

// Currently returns "password"
//
// TODO: implement properly
func getCloudPassphrase() string {
	return "password"
}

// Function return on connection error or connection closing
func establishConnection() error {
	u := getConnectionUrl()

	log.Info("Trying to establish connection", "portalUrl", u)

	headers := http.Header{}
	headers.Add(AuthorizationHeader, getCloudPassphrase())
	headers.Add(RouterIDHeader, GetHostID())

	if connection, _, err := websocket.DefaultDialer.Dial(u, headers); err != nil {
		return err
	} else {
		conn = connection
	}
	log.Info("Successfully connected to portal", "url", u)
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("Failed to read message", "err", err)
			break
		}
		log.Info("Message received", "message", string(message), "messageType", messageType)

		var msg Request
		json.Unmarshal(message, &msg)
		if err := handleRequest(&msg); err != nil {
			log.Error("Failed to handle request, sending internal error response", "err", err)
		}
	}

	return nil
}

// Currently returns 10 seconds
//
// TODO: implement properly
func getReconnectTimeout() time.Duration {
	return 10 * time.Second
}
func getConnectionUrl() string {
	u := url.URL{
		Scheme: "ws",
		Host:   getCloudAddress(),
		Path:   CloudWebSocketPath,
	}
	return u.String()
}
func getCloudAddress() string {
	return "localhost:8080"
}

func handleRequest(r *Request) error {
	cmd, err := commands.CheckCommand(r.Command)
	if err != nil {
		sendError(r.ID, err)
		return err
	}
	data, err := cmd(r.Args)
	if err != nil {
		sendError(r.ID, err)
		return err
	}
	response := Response{
		ID:   r.ID,
		Data: data,
	}
	return sendResponse(response)
}

func GetHostID() string {
	cmd := exec.Command("dmidecode", "--string", "system-uuid")
	output, err := cmd.Output()

	if err != nil {
		log.Error("Failed get host id", "error", err, "cmd", cmd.String())
	}

	return string(output)
}

func sendResponse(response Response) error {
	err := conn.WriteJSON(response)
	if err != nil {
		log.Error("Failed to send response", "error", err, "responseId", response.ID)
	}
	return err
}

func sendError(requestId string, err error) error {
	response := Response{
		ID:    requestId,
		Error: err.Error(),
	}

	return sendResponse(response)
}
