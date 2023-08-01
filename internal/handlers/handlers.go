package handlers

import (
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

type WsResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Client connected !")

	response := WsResponse{}

	response.Message = "Response Connected Server"

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err.Error())
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
