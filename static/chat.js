let chatmsg = Vue.createApp({
    data() {
        return {
            msgs: []
        }
    },
    mounted() {
        function connect(obj) {
            let proto = "ws://";
            if (window.location.protocol == "https") {
                proto = "wss://";
            }
            const url = proto + window.location.host + "/ws/chat";
            let ws = new WebSocket(url);
            ws.onopen = ()=>{
                ws.send(JSON.stringify({
                    pkt: "getchat",
                    body: "aa"
                }))
            }
            ws.onmessage = (e)=>{
                const msg = JSON.parse(e.data);
                obj.push(
                    {
                        ts: msg.hts,
                        body: msg.body
                    }
                )
            }
        }
        connect(this.msgs);
    }
});
chatmsg.mount('#chatmsg');
