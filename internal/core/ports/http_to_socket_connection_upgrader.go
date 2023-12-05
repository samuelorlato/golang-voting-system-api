package ports

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type HTTPToSocketConnectionUpgrader interface {
	Upgrade(http.ResponseWriter, *http.Request, http.Header) (*websocket.Conn, error)
}
