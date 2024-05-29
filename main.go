package main

import (
	"bytes"
	"context"
	"net/http"

	"github.com/deastl/hx-socket-chat/views"
	"github.com/deastl/hx-sockets"
	"github.com/deastl/hx-sockets/compat"
)

func main() {
	mux := http.NewServeMux()
	server := compat.NewNetHttp(mux).(compat.NethttpServer)

	server.Listen("some_message", func(ctx *compat.NethttpClient, msg *hx.Message) {

		state := false

		if msg.Includes["state"] != nil {
			state = msg.Includes["state"].(bool)
		}

		var buff bytes.Buffer

		views.Button(state).Render(context.Background(), &buff)

		ctx.SendStr(buff.String())
	})

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	server.Start("/ws") // where the web socket mounts

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.Main().Render(context.Background(), w)
	})
	http.ListenAndServe(":3000", mux)
}
