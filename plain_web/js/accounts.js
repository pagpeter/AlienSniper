const IP = "localhost"
const PORT = "8080"
const TOKEN = "GMWJGSAPGATLMODYLUMG"
let accounts = []


const get_account_html = (acc) => {
    // acc = { email: '', password: '', type: '' group: '', status: '' }
    let status_color = "red";
    let usable = "no";
    if ((acc.type != "Pending...") && (acc.type)) {
        status_color = "green";
        usable = "yes";
    }

    console.log(usable, status_color, acc);
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
}

const add_account_html = (acc) => {
    let html = get_account_html(acc);
    document.getElementById("accounts_list").innerHTML += html;
}

const add_account_handler = () => {
    const email = document.getElementById('account_email').value.trim();
    const password = document.getElementById('account_password').value.trim();
    const group = document.getElementById('account_group').value.trim();

    const account = {
        email: email,
        password: password,
        group: group || "None",
        type: "Pending...",
        status: "Pending...",
    }
    console.log(account);
    accounts.push(account);
    socket.send(new Packet('add_account', { account: account }).toJson());
    add_account_html(account);
}

const mass_add_accounts_handler = () => {
    const lines = document.getElementById('mass_accounts').value;
    const group = document.getElementById('mass_accounts_group').value;

    const tmpAccs = lines.split('\n');
    tmpAccs.forEach(acc => {
        const account = {
            email: acc.split(',')[0],
            password: acc.split(',')[1],
            group: group || "None",
            type: "Pending...",
            status: "Pending...",
        }
        console.log(account);
        accounts.push(account);
        // socket.send(new Packet('add_account', { account: account }).toJson());
        add_account_html(account);
    });
    const content = {
        account: {
            group: group || "None",
            type: "Pending...",
            status: "Pending...",
            lines: tmpAccs,
        }
    }

    socket.send(new Packet('add_multiple_accounts', content).toJson());
}




// make new connection
let socket = new WebSocket(`ws://${IP}:${PORT}/ws`)

// send auth packet on open
socket.onopen = event => {
    console.log('Connected to server', event);
    socket.send(new Packet('auth', { auth: TOKEN }).toJson());
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
        case 'state_response':
            // console.log(packet.content.state);
            accounts = packet.content.state.accounts || [];
            accounts.forEach(acc => {
                add_account_html(acc);
            });
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
