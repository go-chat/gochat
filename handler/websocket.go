package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chat/gochat/model"
)

func WebSocketConnect(w http.ResponseWriter, r *http.Request) {

	ws, err := Srv.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		encodeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	Srv.Clients[r] = ws
}

func WebSocketDisConnect(w http.ResponseWriter, r *http.Request) {
	delete(Srv.Clients, r)
}

func WebSocketHeartBeat(w http.ResponseWriter, r *http.Request) {
	ws, ok := Srv.Clients[r]
	if !ok {
		encodeErrorResponse(w, fmt.Errorf("websocket not available"), http.StatusInternalServerError)
		return
	}

	ws.WriteJSON(map[string]interface{}{"connected": true})
}

func WebSocketSendMessage(w http.ResponseWriter, r *http.Request) {
	ws, ok := Srv.Clients[r]
	if !ok {
		encodeErrorResponse(w, fmt.Errorf("websocket not available"), http.StatusInternalServerError)
		return
	}

	var msg *model.Message
	err := ws.ReadJSON(msg)
	if err != nil {
		encodeErrorResponse(w, fmt.Errorf("cannot read message, err = %v", err), http.StatusInternalServerError)
		return
	}

	appErr := Srv.Store.IMessage.Save(msg)
	if appErr != nil {
		encodeAppErrorResponse(w, appErr)
		return
	}
}
