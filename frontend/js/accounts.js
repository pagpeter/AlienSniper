const t = getToken();
let accounts = [];

const get_account_html = (acc) => {
    let status_color = "red";
    let usable = "No";
    if (acc.type != "Pending..." && acc.type) {
        status_color = "green";
        usable = "Yes";
    }

    return `
    <tr>
    <td>${acc.email}</td> 
    <td>${acc.type || "None"}</td> 
    <td>
        <span class="text-${status_color}-500">
            ${usable}
        </span>
    </td>
    <td>
        ${acc.group || "None"}
    </td>
  </tr>
  `;
};

const add_account_html = (acc) => {
    let html = get_account_html(acc);
    document.getElementById("accounts_list").innerHTML += html;
};

const get_account_html = (acc) => {
    let status_color = "red";
    let usable = "Offline";
    if (acc.type != "Pending..." && acc.type) {
        status_color = "green";
        usable = "Online";
    }

    return `
    <tr>
    <td>${acc.email}</td> 
    <td>${acc.type || "None"}</td> 
    <td>
        <span class="text-${status_color}-500">
            ${usable}
        </span>
    </td>
    <td>
        ${acc.group || "None"}
    </td>
  </tr>
  `;
};

const add_account_html = (acc) => {
    let html = get_account_html(acc);
    document.getElementById("accounts_list").innerHTML += html;
};

const add_account_handler = () => {
    const email = document.getElementById("account_email").value.trim();
    const password = document.getElementById("account_password").value.trim();
    const group = document.getElementById("account_group").value.trim();

    const account = {
        email: email,
        password: password,
        group: group || "None",
        type: "Pending...",
        status: "Pending...",
    };

    if (email) {
        if (password) {
            accounts.push(account);
            add_account_html(account);
        }
    }

    socket.send(new Packet("add_account", account).toJson());
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
            type: Type,
            host: Host,
        }]
    };

    socket.send(new Packet("add_session", content).toJson());
}

const mass_add_accounts_handler = () => {
    const lines = document.getElementById("mass_accounts").value;
    const group = document.getElementById("mass_accounts_group").value;
    const tmpAccs = lines.split("\n");

    if (lines.length != 0) {
        tmpAccs.forEach((acc) => {
            if (acc.split(":")[0]) {
                if (acc.split(":")[1]) {
                    const account = {
                        email: acc.split(",")[0],
                        password: acc.split(",")[1],
                        group: group || "None",
                        type: "Pending...",
                        status: "Pending...",
                    };
                    accounts.push(account);
                    add_account_html(account);
                }
            }
        });
    }

    const content = {
        account: {
            group: group || "None",
            type: "Pending...",
            status: "Pending...",
            lines: tmpAccs,
        },
    };

    socket.send(new Packet("add_multiple_accounts", content).toJson());
};

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
            accounts = packet.content.state.accounts || [];
            accounts.forEach((acc) => {
                add_account_html(acc);
            });
            break;
        case "config":
            console.log(packet.content.config);
            break;
        case "response":
            console.log(packet.content.response);
            break;
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