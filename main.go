package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var options websocket.AcceptOptions
var idCounter int = 0

type Message struct {
	Typ       string
	Msg       string
	TimeStamp time.Time
	UserId    int
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

func removeUser(userId int) {
	var NewUsers []User
	for _, user := range Users {
		if user.Id != userId {
			NewUsers = append(NewUsers, user)
		} else {
			go Broadcast([]byte(fmt.Sprint("User: ", user, " Left.")))
		}
	}
	Users = NewUsers
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
	fmt.Printf("Users: %+v\n", Users)

	Hello := Message{Typ: "Hello", Msg: fmt.Sprint("Hello user:", thisUser.Id), TimeStamp: time.Now()}
	jsonHello, err := json.Marshal(Hello)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Write(ctx, websocket.MessageText, jsonHello)

	for {
		var msg Message
		err := wsjson.Read(ctx, c, &msg)
		if err != nil {
			fmt.Println(err)
			removeUser(thisUser.Id)
			return
		}
		msg.TimeStamp = time.Now()
		msg.UserId = thisUser.Id
		fmt.Println("received: ", msg, "From user:", thisUser)
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		go Broadcast(jsonMsg)
	}
}

func main() {
	t := time.Now()
	fmt.Println("Hello from server.", t)

	http.HandleFunc("/", wsFunc)
	http.ListenAndServe(":8080", nil)
}
