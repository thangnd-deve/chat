package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
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
	Action        string   `json:"action"`
	Message       string   `json:"message"`
	MessageType   string   `json:"message_type"`
	ConnectedUser []string `json:"connected_user"`
}

type WsPayLoad struct {
	Action   string              `json:"action"`
	UserName string              `json:"username"`
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
		channel := <-wsChan

		switch channel.Action {
		case "username":
			clients[channel.Conn] = channel.UserName
			users := getUserList()

			response.Action = "list_user"
			response.ConnectedUser = users
			broadCastToAll(response)
			break
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", channel.UserName, channel.Message)
			broadCastToAll(response)
			break
		}

	}
}

func getUserList() []string {
	var userList []string
	for _, x := range clients {
		if x != "" {
			userList = append(userList, x)
		}
	}
	sort.Strings(userList)

	return userList
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

func Hom(w http.ResponseWriter, _ *http.Request) {
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
