class Content {
    constructor(content) {
        this.error = content.error;
        this.auth = content.auth;
        this.state = content.state;
        this.config = content.config;
        this.response = content.response;
        this.task = content.task;
        this.account = content.account;
    }
    toJson() {
        return JSON.stringify(this);
    }
}

class Packet {
    constructor(type, content) {
        this.type = type;
        this.content = content;
    }
    toJson() {
        return JSON.stringify(this);
    }
    makeError(error) {
        this.type = 'error';
        this.content = new Content({ error: error });
    }
}

client = async (ip, port, token) => {
    let socket = new WebSocket(`ws://${ip}:${port}/ws`)
    socket.onopen = event => {
        console.log('Connected to server', event);
        const m = new Packet('auth', { auth: token }).toJson();
        console.log(m);
        socket.send(m);
    }
    socket.onmessage = (event) => {
        let packet = JSON.parse(event.data);
        switch (packet.type) {
            case 'error':
                console.log(packet.content.error);
                break;
            case 'auth':
                console.log(packet.content.auth);
                break;
            case 'state':
                console.log(packet.content.state);
                break;
            case 'config':
                console.log(packet.content.config);
                break;
            case 'response':
                console.log(packet.content.response);
                break;
            default: 
                console.log(packet);
        }
    }
}

const IP = "localhost"
const PORT = "8080"
const TOKEN = "GMWJGSAPGATLMODYLUMG"

client(IP, PORT, TOKEN);
