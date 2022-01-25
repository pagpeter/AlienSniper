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
        this.type = "error";
        this.content = new Content({ error: error });
    }
}

const makePopupHTML = (content, action) => {
    id = new Date().getTime();

    return (html = `
    <div class="modal modal-open" id="${id}">
        <div class="modal-box">
            ${content}

            <div class="modal-action">
                ${action.replace("{id}", id)}
            </div>
        </div>
    </div>`);
};

const closePopup = (id) => {
    const popup = document.getElementById(id);
    popup.classList.remove("modal-open");
    // setTimeout(() => {
    //     popup.remove();
    // }, 500);
};

const showPopup = (html) => {
    document.body.insertAdjacentHTML("beforeend", html);
};

const popInfo = (content) => {
    const action = `<button onclick="closePopup('{id}')" for="my-modal-2" class="btn">Ok</button>`;
    const popup = makePopupHTML(content, action);
    showPopup(popup);
};

function modalOpen(id, event) {
    let modal = document.getElementById(id);
    modal.classList.add(event);
}

function modalClose(id, event) {
    let modal = document.getElementById(id);
    modal.classList.remove(event);
}

const formatTime = (time) => {
    time = parseInt((time + ""));
    if (!time) {
        return "-";
    }

    let date = new Date(time * 1000);

    // this will break in 900 years
    if (date.getUTCFullYear() > 3000) {
        console.log("Date to big", time)
        console.log(date / 100)
        date = new Date(time);
        console.log(date)
    }

    return date.toLocaleString().slice(0, -3) + "." + String(date.getMilliseconds());
    // return `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`
}