package pages

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 2 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func PodInfoReloaderWS(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			fmt.Printf("Upgrading request to WebSocket failed: %s\n", err)
		}
		return
	}

	var lastMod time.Time
	if n, err := strconv.ParseInt(r.FormValue("lastMod"), 16, 64); err == nil {
		lastMod = time.Unix(0, n)
	}

	go writer(ws, lastMod)
	reader(ws)
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	err := ws.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		fmt.Printf("setting read deadline failed: %s\n", err.Error())
		return
	}

	ws.SetPongHandler(func(string) error {
		err := ws.SetReadDeadline(time.Now().Add(pongWait))
		return err
	})

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, lastMod time.Time) {
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod)
	fileTicker := time.NewTicker(filePeriod)
	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-fileTicker.C:
			var p []byte
			var err error

			p, lastMod, err = podDataIfModified(lastMod)

			if err != nil {
				if s := err.Error(); s != lastError {
					lastError = s
					p = []byte(lastError)
				}
			} else {
				lastError = ""
			}

			if p != nil {
				err := ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err != nil {
					return
				}
				err = ws.WriteMessage(websocket.TextMessage, p)
				if err != nil {
					return
				}
			}
		case <-pingTicker.C:
			err := ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			err = ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
		}
	}
}

// TODO: cache pod data, modify lastMod correctly
func podDataIfModified(lastMod time.Time) ([]byte, time.Time, error) {

	if !time.Now().After(lastMod) {
		return nil, lastMod, nil
	}

	podDataArr, err := getPodData()
	if err != nil {
		return nil, lastMod, err
	}

	data := struct {
		Items []podData
	}{
		Items: *podDataArr,
	}

	var b bytes.Buffer

	err = pod_table_templ.Execute(&b, data)
	if err != nil {
		fmt.Printf("generating pod data table from template failed: %s\n", err.Error())
		return nil, time.Now(), err
	}

	return b.Bytes(), time.Now(), nil
}
