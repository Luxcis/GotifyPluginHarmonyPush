package main

type GotifyMessage struct {
	Id       uint32
	Appid    uint32
	Message  string
	Title    string
	Priority uint32
	Date     string
}

type GotifyApplication struct {
	DefaultPriority uint32
	Description     string
	ID              uint32
	Image           string
	Internal        bool
	LastUsed        string
	Name            string
	Token           string
}

type Notification struct {
	Category    string      `json:"category"`
	Title       string      `json:"title"`
	Body        string      `json:"body"`
	Image       string      `json:"image"`
	Badge       Badge       `json:"badge"`
	ClickAction ClickAction `json:"clickAction"`
}

type ClickAction struct {
	ActionType uint32            `json:"actionType"`
	Data       map[string]uint32 `json:"data"`
}

type Badge struct {
	AddNum uint32 `json:"addNum"`
	SetNum uint32 `json:"setNum"`
}

type Payload struct {
	Notification Notification `json:"notification"`
}

type Target struct {
	Token []string `json:"token"`
}

type PushOptions struct {
	TestMessage bool `json:"testMessage"`
}

type Message struct {
	Payload     Payload     `json:"payload"`
	Target      Target      `json:"target"`
	PushOptions PushOptions `json:"pushOptions"`
}
