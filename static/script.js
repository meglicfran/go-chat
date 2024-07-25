console.log("Hello from client");

var chat_container = document.getElementById("chat_container")
const ws = new WebSocket("ws://https://go-chat-production-30cb.up.railway.app/ws");
var myId = ""

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
    if (msg.Typ == "Message"){
        console.log(msg.UserId, "==",myId)
        if (msg.UserId==myId) {
            createMessageRightVisual("User " +msg.UserId, msg.Msg, msg.TimeStamp)
        }else{
            createMessageLeftVisual("User " +msg.UserId, msg.Msg, msg.TimeStamp)
        }
    }else{
        myId = msg.Msg
    }
});

function ScrollToBottom(){
    chat_container.scrollTop=chat_container.scrollHeight
}


function SendMessage(){
    text = document.getElementById("messageInput").value.trim();
    if (text != ''){
        message = {
            Typ:"Message",
            Msg:text
        }
        console.log("Sending message ... ", message);
        ws.send(JSON.stringify(message));
        document.getElementById("messageInput").value='';
    }else{
        console.log("Message is empty.");
    }
}

function createMessageLeftVisual(user, text, time){
    html=`  <div class="left">
                <div class="message">
                    <div class="user">
                        ${user}:
                    </div>
                    <div class="border">
                        <div class="text">
                            ${text}
                        </div>
                    </div>
                    <div class="time">
                        ${time}
                    </div>
                </div>
            </div>`
    document.getElementById("chat_container").innerHTML = document.getElementById("chat_container").innerHTML + html
    ScrollToBottom()
}

function createMessageRightVisual(user, text, time){
    html=`  <div class="right">
                <div class="message">
                    <div class="user">
                        ${user}:
                    </div>
                    <div class="border">
                        <div class="text">
                            ${text}
                        </div>
                    </div>
                    <div class="time">
                        ${time}
                    </div>
                </div>
            </div>`
    document.getElementById("chat_container").innerHTML = document.getElementById("chat_container").innerHTML + html
    ScrollToBottom()
}

document.getElementById('messageInput').addEventListener('keypress', function (e) {
    if (e.key === 'Enter') {
        SendMessage();
    }
});
