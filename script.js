console.log("Hello from client");

const ws = new WebSocket("ws://localhost:8080");

ws.addEventListener("open", (event) => {
    message = {
        Typ: "Hello",
        Msg: "Hello Server!"
    }
    ws.send(JSON.stringify(message));
});
  
ws.addEventListener("message", (event) => {
    msg = JSON.parse(event.data)
    console.log(msg);
    createMessageVisual("User " +msg.UserId, msg.Msg, msg.TimeStamp)
});


function SendClicked(){
    text = document.getElementById("inputText").value
    message = {
        Typ:"Message",
        Msg:text
    }
    console.log("Sending message", message)
    ws.send(JSON.stringify(message))
}

function createMessageVisual(user, text, time){
    html=`<div class="message">
                    <div class="border">
                        <div class="user">
                            ${user}:
                        </div>
                        <div class="text">
                            ${text}
                        </div>
                    </div>
                    <div class="time">
                        ${time}
                    </div>
                </div>`
    document.getElementById("chat_container").innerHTML = document.getElementById("chat_container").innerHTML + html
}