package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gotify/plugin-api"
)

// GetGotifyPluginInfo returns gotify plugin info.
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		Name:        "GotifyPluginHarmonyPush",
		Description: "A plugin for HarmonyPush",
		ModulePath:  "github.com/Luxcis/GotifyPluginHarmonyPush",
		Author:      "Luxcis",
		Version:     "1.0.0",
		License:     "MIT",
	}
}

// Plugin is the plugin instance
type Plugin struct {
	ws          *websocket.Conn
	msgHandler  plugin.MessageHandler
	debugLogger *log.Logger
	clientToken string
	gotifyPort  string
	gotifyToken string
	gotifyURL   string
	iconMap     map[uint32]string
	jwt         JwtJson
}

func (p *Plugin) generateJwtToken() string {
	token := GenerateJwtToken(p.jwt)
	p.debugLogger.Println("生成的 JWT: %s\n", token)
	return token
}

func (p *Plugin) sendMsgToHarmony(title string, msg string, icon string, appid uint32) {
	data := Message{
		Payload: Payload{
			Notification: Notification{
				Category: "ACCOUNT",
				Title:    title,
				Body:     msg,
				Image:    icon,
				Badge: Badge{
					AddNum: 1,
				},
				ClickAction: ClickAction{
					ActionType: 0,
					Data: map[string]uint32{
						"appid": appid,
					},
				},
			},
		},
		Target: Target{
			Token: []string{p.clientToken},
		},
		PushOptions: PushOptions{
			TestMessage: false,
		},
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		p.debugLogger.Println("Create json false")
		return
	}
	body := bytes.NewBuffer(payloadBytes)
	// For future debugging
	backupBody := bytes.NewBuffer(body.Bytes())

	req, err := http.NewRequest("POST", "https://push-api.cloud.huawei.com/v3/"+p.jwt.ProjectId+"/messages:send", body)
	if err != nil {
		p.debugLogger.Println("Create request false")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.generateJwtToken())
	req.Header.Set("push-type", "0")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		p.debugLogger.Printf("Send request false: %v\n", err)
		return
	}
	p.debugLogger.Println("HTTP request was sent successfully")

	if resp.StatusCode == http.StatusOK {
		p.debugLogger.Println("The message was forwarded successfully to Harmony")
	} else {
		p.debugLogger.Println("The message was forwarded failed to Harmony")
	}
	// Log infor for debugging
	p.debugLogger.Println("============== Request ==============")
	prettyPrint, err := httputil.DumpRequest(req, true)
	if err != nil {
		p.debugLogger.Printf("%v\n", err)
	}
	p.debugLogger.Printf(string(prettyPrint))
	p.debugLogger.Printf("%v\n", backupBody)

	p.debugLogger.Println("============== Response ==============")
	bodyBytes, err := io.ReadAll(resp.Body)
	p.debugLogger.Printf("%v\n", string(bodyBytes))

	defer resp.Body.Close()
}

func (p *Plugin) connectWebsocket() {
	host := "ws://127.0.0.1:" + p.gotifyPort + "/stream?token=" + p.gotifyToken
	for {
		ws, _, err := websocket.DefaultDialer.Dial(host, nil)
		if err == nil {
			p.ws = ws
			break
		}
		p.debugLogger.Printf("Cannot connect to websocket: %v\n", err)
		time.Sleep(5)
	}
	p.debugLogger.Println("WebSocket connected successfully, ready for forwarding")
}

func (p *Plugin) getWebsocketMsg() {
	go p.connectWebsocket()
	for {
		msg := &GotifyMessage{}
		if p.ws == nil {
			time.Sleep(3)
			continue
		}
		err := p.ws.ReadJSON(msg)
		if err != nil {
			p.debugLogger.Printf("Error while reading websocket: %v\n", err)
			p.connectWebsocket()
			continue
		}
		p.sendMsgToHarmony(msg.Title, msg.Message, p.getIconByAppid(msg.Appid), msg.Appid)
	}
}

func (p *Plugin) getIconByAppid(appid uint32) string {
	_, exists := p.iconMap[appid]
	if !exists {
		api := "http://127.0.0.1:" + p.gotifyPort + "/application?token=" + p.gotifyToken
		resp, err := http.Get(api)
		if err != nil {
			fmt.Println("Error:", err)
			return ""
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			return ""
		}
		var apps []GotifyApplication
		err = json.Unmarshal(body, &apps)
		if err != nil {
			fmt.Println("Error:", err)
			return ""
		}

		for _, app := range apps {
			p.iconMap[app.ID] = app.Image
		}
	}
	return p.gotifyURL + p.iconMap[appid]
}

// SetMessageHandler implements plugin.Messenger
// Invoked during initialization
func (p *Plugin) SetMessageHandler(h plugin.MessageHandler) {
	p.debugLogger = log.New(os.Stdout, "Gotify Plugin Harmony Push: ", log.Lshortfile)
	p.msgHandler = h
}

func (p *Plugin) Enable() error {
	go p.getWebsocketMsg()
	return nil
}

// Disable implements plugin.Plugin
func (p *Plugin) Disable() error {
	if p.ws != nil {
		p.ws.Close()
	}
	return nil
}

// NewGotifyPluginInstance creates a plugin instance for a user context.
func NewGotifyPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	return &Plugin{
		clientToken: os.Getenv("HARMONY_CLIENT_TOKEN"),
		gotifyPort:  os.Getenv("GOTIFY_SERVER_PORT"),
		gotifyToken: os.Getenv("GOTIFY_CLIENT_TOKEN"),
		gotifyURL:   os.Getenv("GOTIFY_SERVER_URL"),
		iconMap:     make(map[uint32]string),
		jwt:         readJwtJsonFile(),
	}
}

func main() {
	panic("this should be built as go plugin")
}
