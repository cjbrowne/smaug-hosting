package message

type MessageType string

const (
	MTSubscribe MessageType = "subscribe"
	MTPublish   MessageType = "publish"
)

type Message struct {
	Type       MessageType       `json:"type"`
	Headers    map[string]string `json:"headers"`
	ContentLen int               `json:"length"`
	Content    interface{}       `json:"content"`
	Game       string            `json:"game"`
}
