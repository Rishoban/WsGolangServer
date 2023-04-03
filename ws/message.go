package ws

type Message struct {
	Type string `json:"type"`
	Id   int    `json:"id"`
	Ts   int64  `json:"ts"`
}
