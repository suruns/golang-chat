var url = "ws://192.168.199.125:8000/ws"
var ws = new WebSocket(url)

function Message(user,message,type) {
    this.user = user;
    this.message = message;
    //1:普通信息，2:心跳信息
    this.type = type;
}
var MsgType = {
    TEXT:1,
    HEARTBEAT:2,
    LOGIN:3,
    LOGOUT:4,
    IMAGE:5
}
var send = function (user,data) {
    // utils.createMessage(user.name,data)
    if (ws.readyState === WebSocket.OPEN){
        ws.send(JSON.stringify(new Message(user,data,MsgType.TEXT)))
    } else {
        alert("websocket异常")
    }
}

function sendLogin() {
    ws.send(JSON.stringify(new Message(user,null,MsgType.LOGIN)))

}
ws.onmessage = function (msg) {
    console.log(msg)
    var message = JSON.parse(msg.data)
    switch (message.type){
        case MsgType.TEXT:
            utils.createMessage(message.user.name,message.message)
            break
        case MsgType.LOGIN:
            utils.createInfo(message.user.name)
            var u = document.getElementsByClassName(message.user.name)[0]
            if(message.user.name !== user.name){
                if(u){
                    u.show()
                }else {
                    utils.createUser(message.user.name,message.user.ip)
                }
            }
            break
        case MsgType.LOGOUT:
            utils.createInfo(message.user.name,MsgType.LOGOUT)
            var u = document.getElementsByClassName(message.user.name)[0]
            u.hide();
            break
        case MsgType.IMAGE:
            utils.createImageMessage(message.user.name,message.message)
            break
    }
}

ws.onopen = function () {
    user.ip = currentIp.innerText
    console.log(user)
    showUserList()
    axios.get("/checkIp?ip="+user.ip)
        .then(function (value) {
                    // console.log(value)
            if(value.data.status === 404){
                layer.show()
            } else if(value.data.status === 200){
                user = value.data.entity
                myName.innerText = user.name
                myIp.innerText = user.ip
                sendLogin()
            }
            })
        .catch(function (reason) {
            console.log(reason)
        })
}
