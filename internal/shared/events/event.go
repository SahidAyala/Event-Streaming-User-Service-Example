package events

type Event struct {
	StreamID   string                 `json:"stream_id"`
	Type       string                 `json:"type"`
	Source     string                 `json:"source"`
	Payload    map[string]interface{} `json:"payload"`
	Metadata   map[string]interface{} `json:"metadata"`
}