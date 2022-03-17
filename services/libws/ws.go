package libws

import (
	"bitbucket.org/smaug-hosting/services/libhttp"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessageType string

const (
	WSMTSubscribe WSMessageType = "subscribe"
	WSMTMessage                 = "message"
	WSMTHandshake               = "handshake"
)

type WSMessage struct {
	Subject    WSMessageType          `json:"subject"`
	Body       map[string]interface{} `json:"body"`
	Token      string                 `json:"token"`
	Connection *websocket.Conn        `json:"-"`
}

type WebSocketSession struct {
	Id    string
	Token string
}

type WebSocketClient struct {
	Connection *websocket.Conn
	Session    WebSocketSession
}

type WebSocket struct {
	readChan   chan WSMessage
	writeChan  chan WSMessage
	AllClients map[string]*WebSocketClient `json:"-"`
}

func (ws WebSocket) Listen() <-chan WSMessage {
	return ws.readChan
}

func (ws WebSocket) Send() chan<- WSMessage {
	return ws.writeChan
}

func (ws WebSocket) handler(response http.ResponseWriter, request *http.Request) {
	logrus.Tracef("Received websocket request")
	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		logrus.Errorf("Could not upgrade websocket: %s", err)
	} else {
		sessionId := uuid.New().String()
		connected := true
		conn.SetCloseHandler(func(code int, text string) error {
			connected = false
			delete(ws.AllClients, sessionId)
			return nil
		})

		messageType, rawMessage, err := conn.ReadMessage()
		if err != nil {
			logrus.Errorf("Handshake message failed: %s", err)
			return
		}
		if messageType == websocket.TextMessage {
			msg := WSMessage{}
			err := json.Unmarshal(rawMessage, &msg)
			if err != nil {
				logrus.Errorf("Could not unmarshal JSON message: %s", err)
				return
			}

			if msg.Subject != WSMTHandshake {
				logrus.Errorf("Did not receive handshake message first")
				return
			}

			ws.AllClients[sessionId] = &WebSocketClient{
				Connection: conn,
				Session: WebSocketSession{
					Id:    sessionId,
					Token: msg.Body["token"].(string),
				},
			}
		}

		logrus.Tracef("Upgraded connection to websocket")
		go func() {
			logrus.Tracef("Websocket connected, spinning up listener...")
			for connected {
				logrus.Tracef("Waiting for message...")
				messageType, rawMessage, err := conn.ReadMessage()
				if err != nil {
					logrus.Errorf("Could not read message from WS stream: %s", err)
					break
				}
				if messageType == websocket.TextMessage {
					msg := WSMessage{}
					err := json.Unmarshal(rawMessage, &msg)
					if err != nil {
						logrus.Errorf("Could not unmarshal JSON message: %s", err)
						continue
					}
					logrus.Tracef("Got message: %+v", msg)
					msg.Connection = conn
					ws.readChan <- msg
				}
			}
		}()
		go func() {
			logrus.Tracef("Websocket connected, spinning up sender...")
			for msg := range ws.writeChan {
				logrus.Tracef("Sending message: %+v", msg)
				m, err := json.Marshal(msg)
				err = msg.Connection.WriteMessage(websocket.TextMessage, m)
				if err != nil {
					logrus.Warnf("Could not send JSON message: %s", err)
				}
			}
			logrus.Tracef("Cleaned up websocket sender goroutine")
		}()
	}
}

func SetupWebsocket(path string) WebSocket {
	ws := WebSocket{
		readChan:   make(chan WSMessage),
		writeChan:  make(chan WSMessage),
		AllClients: make(map[string]*WebSocketClient),
	}

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     ws.handler,
		Pattern:     path,
		Middleware:  nil,
		Description: "Handle websocket requests",
	})

	http.HandleFunc(path, ws.handler)

	return ws
}
