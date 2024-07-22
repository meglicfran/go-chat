package main

import (
	"context"
	"fmt"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var options websocket.AcceptOptions
var idCounter int = 0

type Message struct {
	Msg string
}

type User struct {
	Id   int
	Ws   *websocket.Conn
	Cntx *context.Context
}

var Users []User

func Broadcast(message []byte) {
	for _, user := range Users {
		user.Ws.Write(*user.Cntx, websocket.MessageText, message)
	}
}

func wsFunc(w http.ResponseWriter, r *http.Request) {
	options.InsecureSkipVerify = true
	c, err := websocket.Accept(w, r, &options)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := r.Context()
	idCounter++
	thisUser := User{Ws: c, Cntx: &ctx, Id: idCounter}
	Users = append(Users, thisUser)
	c.Write(ctx, websocket.MessageText, []byte(fmt.Sprint("Hello user:", thisUser)))
	fmt.Printf("Users: %+v\n", Users)

	for {
		var msg Message
		err := wsjson.Read(ctx, c, &msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("received: ", msg, "From user:", thisUser)
		go Broadcast([]byte(fmt.Sprint("Received msg: ", string(msg.Msg), " From user: ", thisUser)))
		//c.Write(ctx, websocket.MessageText, []byte(fmt.Sprint("Received msg: ", string(msg.Msg), " From user: ", thisUser)))
	}
}

func main() {
	fmt.Println("Hello from server.")
	//fs := http.FileServer(http.Dir("./"))
	//http.Handle("/", fs)
	http.HandleFunc("/", wsFunc)
	http.ListenAndServe(":8080", nil)
}
