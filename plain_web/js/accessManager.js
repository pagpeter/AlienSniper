getToken = () => {
  return {
    token: localStorage.getItem("token"),
    ip: localStorage.getItem("ip"),
    port: localStorage.getItem("port"),
  };
};

setToken = (token, ip, port) => {
  localStorage.setItem("token", token);
  localStorage.setItem("ip", ip);
  localStorage.setItem("port", port);
};

checkIfCanAccess = () => {
  let token = getToken();
  if (!token.token || !token.ip || !token.port) {
    return false;
  } else {
    return true;
  }
};

makeInputPopupHTML = () => {
  return `
    <input type="checkbox" id="my-modal-2" class="modal-toggle">
    <div for="my-modal-2" class="modal modal-open">
        <div class="modal-box">
            <div class="form-control">
                <h1 class="text-2xl">Enter your information</h1>
                <!-- IP -->
                <label class="label">
                    <span class="label-text">IP Address:</span>
                </label>
                <input type="text" id="ip_info" placeholder="127.0.0.1" class="input input-bordered">

                <!-- Port -->
                <label class="label">
                    <span class="label-text">Port:</span>
                </label>
                <input type="text" id="port_info" placeholder="8080" class="input input-bordered">

                <!-- Token -->
                <label class="label">
                    <span class="label-text">Token:</span>
                </label>
                <input type="text" id="token_info" placeholder="YourReallySecretToken" class="input input-bordered">


                <div class="modal-action">
                    <button onclick="popupInputHandler()" for="my-modal-2" class="btn btn-primary">Connect</button>
                    <label for="my-modal-2" class="btn">Close</label>
                </div>
            </div>
        </div>
    </div>
    `;
};

popupInputHandler = () => {
  let ip = document.getElementById("ip_info").value.trim();
  let port = document.getElementById("port_info").value.trim();
  let token = document.getElementById("token_info").value.trim();

  if (ip == "" || port == "" || token == "") {
    popInfo("Please fill all fields.");
    return;
  }
  setToken(token, ip, port);
  window.location.href = "/";
};

// main

const canAccess = checkIfCanAccess();
console.log("Has access:", canAccess);
if (!canAccess) {
  const popup = makeInputPopupHTML();
  showPopup(popup);
}
