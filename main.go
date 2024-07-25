package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const SERVER_PORT = ":8080"

var idCounter int = 0

type message struct {
	Typ       string
	Msg       string
	TimeStamp string
	UserId    int
}

type user struct {
	Id   int
	Ws   *websocket.Conn
	Cntx *context.Context
}

var users []user

// Create a "HH:MM:SS" string out of time.Time
func makeTimeStampString(t time.Time) string {
	return fmt.Sprintf("%d:%d:%d", t.Hour(), t.Minute(), t.Second())
}

// Send a hello message to the user
func sayHello(user user) error {
	timeStamp := makeTimeStampString(time.Now())
	Hello := message{
		Typ:       "Hello",
		Msg:       fmt.Sprint(user.Id),
		TimeStamp: timeStamp,
	}
	jsonHello, err := json.Marshal(Hello)
	if err != nil {
		return err
	}
	err = user.Ws.Write(*user.Cntx, websocket.MessageText, jsonHello)
	if err != nil {
		return err
	}
	return nil
}

// Broadcast the message byte slice to all users in Users
func broadcast(message []byte) {
	for _, user := range users {
		user.Ws.Write(*user.Cntx, websocket.MessageText, message)
	}
}

// Removes the user with userId form users
func removeUser(userId int) {
	var NewUsers []user
	for _, user := range users {
		if user.Id != userId {
			NewUsers = append(NewUsers, user)
		} else {
			go broadcast([]byte(fmt.Sprint("User: ", user, " Left.")))
		}
	}
	users = NewUsers
}

// Accepts websocekt handshake, sends a hello message to the user,
// listens for and broadcasts incomming messages.
func websocketHandler(w http.ResponseWriter, r *http.Request) {

	//Acceping websocket handshake
	options := websocket.AcceptOptions{InsecureSkipVerify: true}
	c, err := websocket.Accept(w, r, &options)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Creating a new user
	idCounter++
	ctx := r.Context()
	thisUser := user{Ws: c, Cntx: &ctx, Id: idCounter}
	users = append(users, thisUser)

	//Sending hello message to the user
	err = sayHello(thisUser)
	if err != nil {
		fmt.Println(err)
		return
	}

	//While loop listening to incoming messages
	for {
		var msg message
		err := wsjson.Read(ctx, c, &msg)
		if err != nil {
			fmt.Println(err)
			removeUser(thisUser.Id)
			return
		}

		//Do not broadcast hello messages
		if msg.Typ == "Hello" {
			continue
		}

		//Create Message object and broadcast received message
		msg.TimeStamp = makeTimeStampString(time.Now())
		msg.UserId = thisUser.Id
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		go broadcast(jsonMsg)
	}
}

func main() {
	fmt.Println(makeTimeStampString(time.Now()), "Server started on port", SERVER_PORT)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", websocketHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = SERVER_PORT
	}

	http.ListenAndServe("0.0.0.0:"+port, nil)
}
