package cloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"github.com/zarinit-routers/router-server/pkg/commands"
)

const (
	CloudWebSocketPath = "/api/ipc/connect"
	CloudAuthPath      = "/api/organizations/authorize-node"

	AuthorizationHeader = "Authorization"
	RouterIDHeader      = "X-Router-ID"
	GroupIDHeader       = "X-Group-ID"
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

		if getCloudPassphrase() != "" && getCloudGroupID() != "" {
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

// TODO: implement properly
func getCloudPassphrase() string {
	return "password"
}

// TODO: implement properly
func getCloudGroupID() string {
	return getDummyId()
}

// Function return on connection error or connection closing
func establishConnection() error {
	u := getConnectionUrl()

	log.Info("Trying to establish connection", "portalUrl", u)

	token, err := Authenticate()
	if err != nil {
		return fmt.Errorf("failed authenticate in cloud: %s", err)
	}

	headers := http.Header{}
	headers.Add(AuthorizationHeader, token)

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

func getReconnectTimeout() time.Duration {
	return viper.GetDuration("cloud.reconnect-timeout")
}
func getConnectionUrl() string {
	u := url.URL{
		Scheme: "ws",
		Host:   getCloudWsAddress(),
		Path:   CloudWebSocketPath,
	}
	return u.String()
}
func getCloudWsAddress() string {
	return fmt.Sprintf("%s:%d", viper.GetString("cloud.address"), viper.GetInt("cloud.ws-port"))
}
func getCloudAddress() string {
	return viper.GetString("cloud.address")
}

func handleRequest(r *Request) error {
	log.Info("Handling request", "id", r.ID, "command", r.Command)
	cmd, err := commands.CheckCommand(r.Command)
	if err != nil {
		sendError(r.ID, err)
		return err
	}
	data, err := cmd(r.Args)
	if err != nil {
		log.Error("Failed to execute command, sending error", "err", err)
		sendError(r.ID, err)
		return err
	}
	response := Response{
		ID:   r.ID,
		Data: data,
	}

	log.Info("Request handled", "id", r.ID, "command", r.Command, "response", response.Data)
	return sendResponse(response)
}

func getDummyId() string {
	log.Warn("Using dummy UUID, do not use this in production. Even in development remove it ASAP")
	return "00000000-0000-0000-0000-000000000000"
}

func GetHostID() string {
	if viper.GetBool("dev-test") {
		log.Warn("I currently send random UUID, remove it ASAP, do not use in production") // TODO: remove this
		return getDummyId()
	} else {
		return viper.GetString("device.id")
	}
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
