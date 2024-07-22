console.log("Hello from client");

const ws = new WebSocket("ws://localhost:8080");

ws.addEventListener("open", (event) => {
    message = {
        Msg: "Hello Server!"
    }
    ws.send(JSON.stringify(message));
});
  
ws.addEventListener("message", (event) => {
    console.log(event.data);
});


function SendClicked(){
    text = document.getElementById("inputText").value
    message = {
        msg:text
    }
    console.log("Sending message", message)
    ws.send(JSON.stringify(message))
}