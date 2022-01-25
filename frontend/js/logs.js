const t = getToken();

const add_logs = (acc) => {
    var content = "";
    // loop through the logs
    for (const x of(acc.sends || [])) {
        // for each log (for each past snipe)

        // create HTML for the requests
        x.content.forEach((l) => {
            var logHTML = "";
            l.timestamp.forEach((k) => {
                // console.log(k)
                time = String(k).split(":")[0];
                statuscode = String(k).split(":")[1];

                if (statuscode == "200") {
                    bg = "success";
                } else {
                    bg = "error";
                }

                logHTML += `
                  <span class="${
                    statuscode == "200" ? "text-green-500" : "text-red-500"
                  }">[${statuscode}]</span>
                  <span>${formatTime(time)}<br></span>`;
            });

            content += `<div class="bg-${bg} p-2 rounded-md shadow mt-4"><details>
            <summary>
                <h1 class="text-md font-mono">${l.email}</h1>
                <h2 class="text-sm font-mono">${l.ip}</h2>
            </summary>
            <div class="font-mono text-sm mt-2 p-3 bg-neutral ">
                <p>
                  <p>
                    ${logHTML}
                  </p>
                </p>
            </div>
        </details></div>`;
        });

        // create HTML for the log thing

        // logHTML = "";
    }

    statusC = acc.success ? "Yes" : "No";
    bgC = acc.success ? "text-green-500" : "text-red-500";

    return `<div id="${acc.name}" class="modal modal-closed">

    <div class="modal-box">
        <h1 class="text-2xl">Logs for
            <span class="kbd">${acc.name}</span>
        </h1>

        <div class="m-2 p-5 ">

        ${content}  
    
        </div>
        <div class="modal-action">
            <label onclick="modalClose('${acc.name}', 'modal-open')" class="btn">Done</label>
        </div>
    </div>
    </div>
    
    <tr class="hover" onclick="modalOpen('${acc.name}', 'modal-open')">
    <td class="row-data">${acc.name}</td>
    <td class="row-data">${acc.requests}</td>
    <td class="row-data">${acc.delay}</td>
    <td class="row-data">
        <span class="${bgC}">
        ${statusC}
    </span>
    </td>
</tr>`;
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
            console.log(packet.content.error);
            break;
        case "auth":
            console.log(packet.content.auth);
            break;
        case "state_response":
            accs = packet.content.state.logs;
            for (const x of accs) {
                document.getElementById("table1").innerHTML += add_logs(x);
            }
            break;
        case "config":
            console.log(packet.content.config);
            break;
        case "add_task_response":
            // add_task_html(packet.content.task);
            // console.log(packet.content.response);
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