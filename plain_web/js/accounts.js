const IP = "localhost"
const PORT = "8080"
const TOKEN = "GMWJGSAPGATLMODYLUMG"
let accounts = []

// make new connection
let socket = new WebSocket(`ws://${ip}:${port}/ws`)

// send auth packet on open
socket.onopen = event => {
    console.log('Connected to server', event);
    socket.send(new Packet('auth', { auth: token }).toJson());
    socket.send(new Packet('get_state', {}).toJson());
}

// handle incoming packets
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
            accounts = packet.content.state.accounts;
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