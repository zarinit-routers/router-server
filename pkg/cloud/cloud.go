package cloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	config "github.com/zarinit-routers/router-server/pkg/cloud/config"
	"github.com/zarinit-routers/router-server/pkg/commands"
)

const (
	AuthorizationHeader = "Authorization"
)

var (
	conn *websocket.Conn
)

type Status struct {
	Connected bool `json:"connected"`
}

func GetStatus() Status {
	st := Status{
		Connected: conn != nil,
	}
	return st
}

type Request struct {
	Command string         `json:"command"`
	Args    map[string]any `json:"args"`
	ID      string         `json:"requestId"`
}
type Response struct {
	Data  any    `json:"data"`
	ID    string `json:"requestId"`
	Error string `json:"error"`
}

func ServeConnection() {
	for {
		tryConnect()

		timeout := getReconnectTimeout()
		log.Debug("Reconnect timeout", "timeout", timeout)
		time.Sleep(timeout)
	}
}

func tryConnect() {
	cloud := config.GetConnectionConfig()
	err := cloud.Validate()
	if err != nil {
		log.Debug("Failed to get cloud config", "err", err) // Error from getting connection config is not a real problem
		return
	}
	if err := establishConnection(cloud); err != nil {
		log.Error("Failed establish connection with cloud", "error", err, "config", cloud)
		return
	}
}

// Function return on connection error or connection closing
func establishConnection(conf *config.ConnectionConfig) error {
	token, err := authenticate(*conf)
	if err != nil {
		return fmt.Errorf("failed authenticate on cloud: %s", err)
	}

	u := conf.GetWebsocketURL()

	log.Info("Trying to establish connection", "portalUrl", u)

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
	return config.GetReconnectTimeout()
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
	return viper.GetString("device.id")
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
