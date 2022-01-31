const t = getToken();
let accounts = [];

const get_vps_html = (acc) => {
    let color = "red";
    if (acc.status == "Online") {
        color = "green";
    }

    return `
    <tr>
    <td>${acc.ip}</td> 
    <td>${acc.host}</td> 
    <td>
        <span class="text-${color}-500">
            ${acc.status}
        </span>
    </td>
    <td>${acc.group}</td> 
  </tr>
  `;
};

const add_vps_html = (acc) => {
    let html = get_vps_html(acc);
    document.getElementById("vps_list").innerHTML += html;
};

const send_session = () => {

    const IP = document.getElementById("vps_ip").value.trim();
    const Password = document.getElementById("vps_password").value.trim();
    const Host = document.getElementById("vps_user").value.trim();
    const Type = document.getElementById("vps_group").value.trim();

    const content = {
        sessions: [{
            ip: IP,
            port: "22",
            password: Password,
            group: Type,
            host: Host,
            status: "Offline",
        }]
    };

    socket.send(new Packet("add_session", content).toJson());
}

// make new connection
let socket = null;
try {
    socket = new WebSocket(`ws://${t.ip}:${t.port}/ws`);
} catch (e) {
    console.log(e);
}

// send auth packet on open
socket.onopen = (event) => {
    console.log("Connected to server", event);
    socket.send(
        new Packet("auth", { auth: t.token, response: { message: "web" } }).toJson()
    );
    socket.send(new Packet("get_state", {}).toJson());
};

// handle incoming packets
socket.onmessage = (event) => {
    let packet = JSON.parse(event.data);

    switch (packet.type) {
        case "error":
            popInfo(
                packet.content.response.error
            );
            break;
        case "auth":
            console.log(packet.content.auth);
            break;
        case "state_response":
            accounts = packet.content.state.sessions || [];
            accounts.forEach((acc) => {
                add_vps_html(acc);
            });
            break;
        case "config":
            console.log(packet.content.config);
            break;
        case "response":
            console.log(packet.content.response);
            break;
        case "add_session_response":
            accounts = packet.content.sessions || [];
            accounts.forEach((acc) => {
                add_vps_html(acc);
            });
            break
        default:
            console.log(packet);
    }
};

alrShowedError = false;
socket.onclose = (event) => {
    console.log("Disconnected from server", event);

    if (!alrShowedError) {
        popInfo(
            "There was an error while connecting to the server. Please check if its running and try again."
        );
        alrShowedError = true;
    }
};

socket.onerror = (event) => {
    console.log("Error connecting to server", event);

    if (!alrShowedError) {
        popInfo(
            "There was an error while connecting to the server. Please check if its running and try again."
        );
        alrShowedError = true;
    }
};