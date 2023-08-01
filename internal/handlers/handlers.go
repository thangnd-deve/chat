package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var view = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var wsChan = make(chan WsPayLoad)
var clients = make(map[WebSocketConnection]string)

type WebSocketConnection struct {
	*websocket.Conn
}
type WsResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

type WsPayLoad struct {
	Action   string              `json:"action"`
	UserName string              `json:"user_name"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Client connected !")

	response := WsResponse{}

	response.Message = "Response Connected Server"

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err.Error())
	}

	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Err", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayLoad
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			log.Println(err.Error())
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var response WsResponse

	for {
		err := <-wsChan

		response.Action = "Got Here"
		response.Message = fmt.Sprintf("Some message and actione %s", err.Action)
		broadCastToAll(response)
	}
}

func broadCastToAll(response WsResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Websocket Error")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func Hom(w http.ResponseWriter, r *http.Request) {
	err := renderHtml(w, "home.html", nil)

	if err != nil {
		log.Println(err.Error())
		return
	}
}

func renderHtml(w http.ResponseWriter, template string, data jet.VarMap) error {
	view, err := view.GetTemplate(template)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = view.Execute(w, data, nil)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
