const add_logs = (acc) => {
    let content = ""
    for (const x of acc.sends) {
        if (String(x.content).match("200")) {
            status_code = "green"
            bg = "success"
        } else {
            status_code = "red"
            bg = "error"
        }

        content += `<div class="bg-${bg} p-2 rounded-md shadow mt-4"><details>
            <summary>
                <h1 class="text-md font-mono">${x.email}</h1>
                <h2 class="text-sm font-mono">${x.ip}</h2>
            </summary>
            <div class="font-mono text-sm mt-2 p-3 bg-neutral ">
                <p><span class="text-${status_code}-500">${String(x.content)}</span></p>
            </div>
        </details></div>`
    }

    if (acc.success == true) {
        statusC = "Yes"
        bgC = "green"
    } else {
        statusC = "No"
        bgC = "red"
    }

    return `<div id="${acc.name}" class="modal modal-closed">

    <div class="modal-box">
        <h1 class="text-2xl">Logs for
            <span class="kbd">${acc.name}</span>
        </h1>

        <p class="text-2xl">Requests: ${acc.requests}</p>
        <p class="text-2xl">Delay: ${acc.delay}</p>
        <p class="text-2xl">Success: ${acc.success}</p>

        <div class="m-2 p-5">

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
        <span class="text-${bgC}-500">
        ${statusC}
    </span>
    </td>
</tr>`
}