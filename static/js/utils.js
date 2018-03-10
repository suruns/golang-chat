var utils = {
    createMessage: function (name,message) {
        var msg = document.createElement("div")
        if(this.isMyself(name)){
            msg.className = "message mine-message"
        } else {
            msg.className = "message other-message"
        }
        var msgName = document.createElement("div")
        msgName.className = "name"
        msgName.innerText = name
        msg.appendChild(msgName)
        var msgText = document.createElement("div")
        msgText.className = "text"
        msgText.innerHTML = message
        msg.appendChild(msgText)
        chatMessages.appendChild(msg)
        chatMessages.scrollTop = chatMessages.scrollHeight;
    },
    createImageMessage: function (name,url) {
        var msg = document.createElement("div")
        if(this.isMyself(name)){
            msg.className = "message mine-message"
        } else {
            msg.className = "message other-message"
        }
        var msgName = document.createElement("div")
        msgName.className = "name"
        msgName.innerText = name
        msg.appendChild(msgName)
        var msgText = document.createElement("div")
        msgText.className = "text"
        var msgA = document.createElement("a")
        msgA.href="javascript:;"
        var msgImg = document.createElement("img")
        msgImg.src = url
        msgA.appendChild(msgImg)
        msgText.appendChild(msgA)
        msg.appendChild(msgText)
        chatMessages.appendChild(msg)
        chatMessages.scrollTop = chatMessages.scrollHeight;
    },
    isMyself: function (name) {
        if (name === user.name){
            return true
        }
        return false
    },
    getIp: function (callback) {
        var ip_dups = {}
        var RTCPeerConnection = window.RTCPeerConnection || window.mozRTCPeerConnection || window.webkitRTCPeerConnection
        if (!RTCPeerConnection) {
            var iframe = document.createElement('iframe')
            iframe.sandbox = 'allow-same-origin'
            iframe.style.display = 'none'
            document.body.appendChild(iframe)
            var win = iframe.contentWindow
            window.RTCPeerConnection = win.RTCPeerConnection
            window.mozRTCPeerConnection = win.mozRTCPeerConnection
            window.webkitRTCPeerConnection = win.webkitRTCPeerConnection
            RTCPeerConnection = window.RTCPeerConnection || window.mozRTCPeerConnection || window.webkitRTCPeerConnection
        }
        var mediaConstraints = {
            optional: [{RtpDataChannels: true}]
        }
        var servers = undefined
        if (window.webkitRTCPeerConnection) {
            servers = {iceServers: [{urls: 'stun:stun.services.mozilla.com'}]}
        }
        var pc = new RTCPeerConnection(servers, mediaConstraints)
        pc.onicecandidate = function (ice) {
            if (ice.candidate) {
                var ip_regex = /([0-9]{1,3}(\.[0-9]{1,3}){3})/
                var ip_addr = ip_regex.exec(ice.candidate.candidate)
                if (ip_dups[ip_addr] === undefined) {
                    callback(ip_addr)
                }
                ip_dups[ip_addr] = true
            }
        }
        pc.createDataChannel('')
        pc.createOffer(function (result) {
            pc.setLocalDescription(result, function () {}, function () {})
        }, function () {})
    },
    createUser: function (name,ip) {
        var chatUser = document.createElement("div")
        chatUser.className = "chat-user "+name
        var chatUserName = document.createElement("p")
        chatUserName.innerText = name;
        chatUser.appendChild(chatUserName)
        var chatUserIp = document.createElement("p")
        chatUserIp.innerText = ip;
        chatUser.appendChild(chatUserIp)

        chatUserList.appendChild(chatUser)
    },
    createInfo:function (name,type) {
        var chatInfo = document.createElement("div")
        chatInfo.className = "chat-info"
        var chatSpan = document.createElement("span")
        if(type === MsgType.LOGOUT){
            chatSpan.innerText = name+" 退出了聊天室"
        }else {
            chatSpan.innerText = name+" 登陆了聊天室"
        }
        chatInfo.appendChild(chatSpan)

        chatMessages.appendChild(chatInfo)
    }
}